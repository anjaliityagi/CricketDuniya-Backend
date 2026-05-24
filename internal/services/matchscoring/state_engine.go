package matchscoring

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services/scoring"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InningsState struct {
	InningsID    string  `db:"innings_id"`
	StrikerID    *string `db:"striker_id"`
	NonStrikerID *string `db:"non_striker_id"`
	BowlerID     *string `db:"bowler_id"`
	TotalRuns    int     `db:"total_runs"`
	TotalWickets int     `db:"total_wickets"`
	LegalBalls   int     `db:"legal_balls"`
	CurrentOver  int     `db:"current_over"`
	CurrentBall  int     `db:"current_ball"`
	Status       string  `db:"status"`
	UpdatedAt    string  `db:"updated_at"`
}

type ProcessBallResult struct {
	State dto.InningsStateResponse
}

func (e *Engine) ProcessBall(req dto.BallInputRequest) (*ProcessBallResult, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	inningsMeta, err := getInningsMeta(tx, req.InningsID)
	if err != nil {
		return nil, err
	}

	state, err := getOrCreateState(tx, req, inningsMeta)
	if err != nil {
		return nil, err
	}
	if state.Status != "live" {
		return nil, errors.New("innings is not live")
	}

	overLimitBalls := inningsMeta.OversPerInnings * 6
	if state.LegalBalls >= overLimitBalls {
		return nil, errors.New("innings already completed by overs")
	}

	activeStriker, activeNonStriker, activeBowler, err := resolveActors(req, state)
	if err != nil {
		return nil, err
	}

	if err := validateActorTeams(tx, inningsMeta.BattingTeamID, inningsMeta.BowlingTeamID, activeStriker, activeNonStriker, activeBowler); err != nil {
		return nil, err
	}

	if state.LegalBalls > 0 && state.CurrentBall == 0 {
		previousOverBowlerID, err := getPreviousOverBowlerID(tx, req.InningsID)
		if err != nil {
			return nil, err
		}
		if previousOverBowlerID != "" && previousOverBowlerID == activeBowler {
			return nil, errors.New("same bowler cannot bowl consecutive overs")
		}
	}
	internalReq, legalBall := mapBallRequest(req, activeStriker, activeNonStriker, activeBowler, state)

	matchUpdate, err := e.Process(internalReq)
	if err != nil {
		return nil, err
	}

	bat, bowl, field, err := scoring.NewEngine().Process(internalReq)
	if err != nil {
		return nil, err
	}

	ballEventID, err := saveBallTx(tx, internalReq)
	if err != nil {
		return nil, err
	}

	if err := persistFantasyTx(tx, internalReq, ballEventID, bat, bowl, field); err != nil {
		return nil, err
	}

	nextStriker := activeStriker
	nextNonStriker := activeNonStriker
	nextBowler := activeBowler
	needsNextBatter := false

	if req.IsWicket && req.DismissedPlayerID != nil {
		dismissedID := req.DismissedPlayerID.String()
		if dismissedID == nextStriker {
			if req.NextBatterID == nil {
				return nil, errors.New("next_batter_id is required for striker dismissal")
			}
			nextStriker = req.NextBatterID.String()
			needsNextBatter = true
		}
		if dismissedID == nextNonStriker {
			if req.NextBatterID == nil {
				return nil, errors.New("next_batter_id is required for non-striker dismissal")
			}
			nextNonStriker = req.NextBatterID.String()
			needsNextBatter = true
		}
	}

	totalRunsForStrike := internalReq.TotalRuns
	if totalRunsForStrike%2 == 1 {
		nextStriker, nextNonStriker = nextNonStriker, nextStriker
	}

	overCompleted := false
	nextOver := state.CurrentOver
	nextBall := state.CurrentBall
	nextLegalBalls := state.LegalBalls
	if legalBall {
		nextLegalBalls++
		nextBall++
		if nextBall >= 6 {
			overCompleted = true
			nextOver++
			nextBall = 0
			nextStriker, nextNonStriker = nextNonStriker, nextStriker
		}
	}

	inningsCompleted := false
	if nextLegalBalls >= overLimitBalls {
		inningsCompleted = true
	}

	if matchUpdate.InningsWickets >= 10 {
		inningsCompleted = true
	}
	if inningsMeta.TargetRuns != nil && inningsMeta.InningsNo == 2 && matchUpdate.InningsRuns >= *inningsMeta.TargetRuns {
		inningsCompleted = true
	}

	status := "live"
	if inningsCompleted {
		status = "completed"
	}

	if _, err := tx.Exec(`
		UPDATE innings_state
		SET striker_id = $2,
			non_striker_id = $3,
			bowler_id = $4,
			total_runs = $5,
			total_wickets = $6,
			legal_balls = $7,
			current_over = $8,
			current_ball = $9,
			status = $10,
			updated_at = NOW()
		WHERE innings_id = $1
	`, req.InningsID, nextStriker, nextNonStriker, nextBowler, matchUpdate.InningsRuns, matchUpdate.InningsWickets, nextLegalBalls, nextOver, nextBall, status); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(`
		UPDATE innings
		SET total_runs = $2,
			total_wickets = $3,
			status = $4
		WHERE id = $1
	`, req.InningsID, matchUpdate.InningsRuns, matchUpdate.InningsWickets, status); err != nil {
		return nil, err
	}

	if inningsCompleted {
		if err := ensureNextInningsTx(tx, inningsMeta, matchUpdate.InningsRuns); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp := dto.InningsStateResponse{
		InningsID:        req.InningsID.String(),
		StrikerID:        strPtr(nextStriker),
		NonStrikerID:     strPtr(nextNonStriker),
		BowlerID:         strPtr(nextBowler),
		TotalRuns:        matchUpdate.InningsRuns,
		TotalWickets:     matchUpdate.InningsWickets,
		LegalBalls:       nextLegalBalls,
		CurrentOver:      nextOver,
		CurrentBall:      nextBall,
		OverCompleted:    overCompleted,
		InningsCompleted: inningsCompleted,
		NeedsNextBowler:  overCompleted && !inningsCompleted,
		NeedsNextBatter:  needsNextBatter,
	}
	if inningsMeta.TargetRuns != nil && inningsMeta.InningsNo == 2 {
		reqRuns := *inningsMeta.TargetRuns - matchUpdate.InningsRuns
		if reqRuns < 0 {
			reqRuns = 0
		}
		ballsLeft := overLimitBalls - nextLegalBalls
		if ballsLeft < 0 {
			ballsLeft = 0
		}
		resp.RequiredRunsToWin = &reqRuns
		resp.BallsRemaining = &ballsLeft
	}

	return &ProcessBallResult{State: resp}, nil
}

type inningsMeta struct {
	MatchID         string `db:"match_id"`
	InningsNo       int    `db:"innings_no"`
	BattingTeamID   string `db:"batting_team_id"`
	BowlingTeamID   string `db:"bowling_team_id"`
	OversPerInnings int    `db:"overs_per_innings"`
	TargetRuns      *int   `db:"target_runs"`
}

func getInningsMeta(tx *sqlx.Tx, inningsID uuid.UUID) (*inningsMeta, error) {
	var m inningsMeta
	err := tx.Get(&m, `
		SELECT i.match_id, i.innings_no, i.batting_team_id, i.bowling_team_id,
			   m.overs_per_innings, i.target_runs
		FROM innings i
		JOIN matches m ON m.id = i.match_id
		WHERE i.id = $1
	`, inningsID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func ensureNextInningsTx(tx *sqlx.Tx, meta *inningsMeta, firstInningsRuns int) error {
	if meta.InningsNo != 1 {
		return nil
	}

	var inningsCount int
	if err := tx.Get(&inningsCount, `SELECT COUNT(1) FROM innings WHERE match_id = $1 AND innings_no = 2`, meta.MatchID); err != nil {
		return err
	}
	if inningsCount > 0 {
		return nil
	}

	targetRuns := firstInningsRuns + 1
	_, err := tx.Exec(`
		INSERT INTO innings (
			match_id,
			batting_team_id,
			bowling_team_id,
			innings_no,
			target_runs,
			status,
			started_at
		) VALUES ($1, $2, $3, 2, $4, 'live', NOW())
	`, meta.MatchID, meta.BowlingTeamID, meta.BattingTeamID, targetRuns)
	return err
}

func getPreviousOverBowlerID(tx *sqlx.Tx, inningsID uuid.UUID) (string, error) {
	var bowlerID sql.NullString
	err := tx.Get(&bowlerID, `
		SELECT bowler_id
		FROM ball_events
		WHERE innings_id = $1
		  AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`, inningsID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	if !bowlerID.Valid {
		return "", nil
	}
	return bowlerID.String, nil
}

func getOrCreateState(tx *sqlx.Tx, req dto.BallInputRequest, meta *inningsMeta) (*InningsState, error) {
	var state InningsState
	err := tx.Get(&state, `SELECT * FROM innings_state WHERE innings_id = $1 FOR UPDATE`, req.InningsID)
	if err == nil {
		return &state, nil
	}

	if req.StrikerID == nil || req.NonStrikerID == nil || req.BowlerID == nil {
		return nil, errors.New("striker_id, non_striker_id and bowler_id are required for first ball")
	}

	_, err = tx.Exec(`
		INSERT INTO innings_state (innings_id, striker_id, non_striker_id, bowler_id, total_runs, total_wickets, legal_balls, current_over, current_ball, status)
		VALUES ($1,$2,$3,$4,0,0,0,0,0,'live')
	`, req.InningsID, req.StrikerID, req.NonStrikerID, req.BowlerID)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&state, `SELECT * FROM innings_state WHERE innings_id = $1 FOR UPDATE`, req.InningsID); err != nil {
		return nil, err
	}
	_ = meta
	return &state, nil
}

func resolveActors(req dto.BallInputRequest, state *InningsState) (striker, nonStriker, bowler string, err error) {
	if req.StrikerID != nil {
		striker = req.StrikerID.String()
	} else if state.StrikerID != nil {
		striker = *state.StrikerID
	}
	if req.NonStrikerID != nil {
		nonStriker = req.NonStrikerID.String()
	} else if state.NonStrikerID != nil {
		nonStriker = *state.NonStrikerID
	}
	if req.BowlerID != nil {
		bowler = req.BowlerID.String()
	} else if state.BowlerID != nil {
		bowler = *state.BowlerID
	}
	if striker == "" || nonStriker == "" || bowler == "" {
		return "", "", "", errors.New("active striker, non_striker and bowler are required")
	}
	return striker, nonStriker, bowler, nil
}

func validateActorTeams(tx *sqlx.Tx, battingTeamID, bowlingTeamID, strikerID, nonStrikerID, bowlerID string) error {
	var c int
	if err := tx.Get(&c, `SELECT COUNT(1) FROM team_players WHERE id IN ($1,$2) AND team_id = $3 AND deleted_at IS NULL`, strikerID, nonStrikerID, battingTeamID); err != nil {
		return err
	}
	if c != 2 {
		return errors.New("strikers must belong to batting team")
	}
	if err := tx.Get(&c, `SELECT COUNT(1) FROM team_players WHERE id = $1 AND team_id = $2 AND deleted_at IS NULL`, bowlerID, bowlingTeamID); err != nil {
		return err
	}
	if c != 1 {
		return errors.New("bowler must belong to bowling team")
	}
	return nil
}

func mapBallRequest(req dto.BallInputRequest, strikerID, nonStrikerID, bowlerID string, state *InningsState) (dto.BallRequest, bool) {
	bt := strings.ToLower(strings.TrimSpace(req.BallType))
	if bt == "" {
		bt = "normal"
	}
	totalRuns := req.TotalRuns
	if totalRuns == 0 {
		totalRuns = req.RunsOffBat + req.Extras
	}
	wides := 0
	noBalls := 0
	byes := 0
	legByes := 0
	legal := true
	switch bt {
	case "wide":
		wides = req.Extras
		legal = false
	case "no_ball":
		noBalls = req.Extras
		legal = false
	case "bye":
		byes = req.Extras
	case "leg_bye":
		legByes = req.Extras
	}
	deliveryNo := state.CurrentBall + 1
	if !legal {
		deliveryNo = state.CurrentBall + 1
	}
	return dto.BallRequest{
		InningsID:         req.InningsID,
		MatchID:           req.MatchID,
		StrikerID:         uuid.MustParse(strikerID),
		NonStrikerID:      uuid.MustParse(nonStrikerID),
		BowlerID:          uuid.MustParse(bowlerID),
		BallNo:            state.CurrentOver + 1,
		DeliveryNo:        deliveryNo,
		BallType:          bt,
		RunsScored:        totalRuns,
		RunsOffBat:        req.RunsOffBat,
		Extras:            req.Extras,
		TotalRuns:         totalRuns,
		IsDotBall:         totalRuns == 0,
		IsBoundaryFour:    req.RunsOffBat == 4,
		IsBoundarySix:     req.RunsOffBat == 6,
		IsWicket:          req.IsWicket,
		DismissalType:     req.DismissalType,
		DismissedPlayerID: req.DismissedPlayerID,
		FielderID:         req.FielderID,
		Wides:             wides,
		NoBalls:           noBalls,
		Byes:              byes,
		LegByes:           legByes,
	}, legal
}

func saveBallTx(tx *sqlx.Tx, req dto.BallRequest) (string, error) {
	query := `
	INSERT INTO ball_events (
		id, innings_id, match_id, striker_id, non_striker_id, bowler_id,
		ball_no, delivery_no, ball_type, runs_scored, extras, total_runs,
		is_dot_ball, is_wicket, dismissal_type,
		dismissed_player_id, fielder_id, wides, no_balls, byes, leg_byes, created_at
	) VALUES (
		gen_random_uuid(), $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,
		$12,$13,$14,$15,$16,$17,$18,$19,$20, NOW()
	) RETURNING id
	`
	var id string
	if err := tx.QueryRowx(query,
		req.InningsID, req.MatchID, req.StrikerID, req.NonStrikerID, req.BowlerID,
		req.BallNo, req.DeliveryNo, req.BallType, req.RunsScored, req.Extras, req.TotalRuns,
		req.IsDotBall, req.IsWicket, req.DismissalType,
		req.DismissedPlayerID, req.FielderID, req.Wides, req.NoBalls, req.Byes, req.LegByes,
	).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func persistFantasyTx(tx *sqlx.Tx, req dto.BallRequest, ballEventID string, battingPoints, bowlingPoints, fieldingPoints int) error {

	if req.StrikerID.String() != "" && battingPoints != 0 {
		if err := upsertFantasyPointsTx(tx, req.MatchID.String(), req.StrikerID.String(), battingPoints, "batting_points"); err != nil {
			return err
		}
		if err := insertPointEventTx(tx, req.MatchID.String(), req.StrikerID.String(), ballEventID, "batting", "ball_batting_points", battingPoints); err != nil {
			return err
		}
	}
	if req.BowlerID.String() != "" && bowlingPoints != 0 {
		if err := upsertFantasyPointsTx(tx, req.MatchID.String(), req.BowlerID.String(), bowlingPoints, "bowling_points"); err != nil {
			return err
		}
		if err := insertPointEventTx(tx, req.MatchID.String(), req.BowlerID.String(), ballEventID, "bowling", "ball_bowling_points", bowlingPoints); err != nil {
			return err
		}
	}
	if req.FielderID != nil && req.FielderID.String() != "" && fieldingPoints != 0 {
		if err := upsertFantasyPointsTx(tx, req.MatchID.String(), req.FielderID.String(), fieldingPoints, "fielding_points"); err != nil {
			return err
		}
		if err := insertPointEventTx(tx, req.MatchID.String(), req.FielderID.String(), ballEventID, "fielding", "ball_fielding_points", fieldingPoints); err != nil {
			return err
		}
	}
	return nil
}

func upsertFantasyPointsTx(tx *sqlx.Tx, matchID, matchPlayerID string, points int, bucket string) error {
	query := `UPDATE player_match_stats
		SET ` + bucket + ` = COALESCE(` + bucket + `, 0) + $3,
			fantasy_points = COALESCE(fantasy_points, 0) + $3,
			updated_at = NOW()
		WHERE match_id = $1
		  AND team_player_id = $2`
	res, err := tx.Exec(query, matchID, matchPlayerID, points)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}
	_, err = tx.Exec(`INSERT INTO player_match_stats (match_id, player_id, team_player_id,`+bucket+`, fantasy_points, updated_at) SELECT $1, tp.player_id, tp.id, $3, $3, NOW() FROM team_players tp WHERE tp.id = $2 AND tp.player_id IS NOT NULL`, matchID, matchPlayerID, points)
	return err
}

func insertPointEventTx(tx *sqlx.Tx, matchID, matchPlayerID, ballEventID, category, ruleName string, points int) error {
	_, err := tx.Exec(`INSERT INTO point_events (match_id, user_id, ball_event_id, category, rule_name, points)
		SELECT $1, tp.player_id, $2, $3::point_category, $4, $5
		FROM team_players tp
		WHERE tp.id = $6
		  AND tp.player_id IS NOT NULL`,
		matchID, ballEventID, category, ruleName, points, matchPlayerID)
	return err
}

func (e *Engine) OverrideState(inningsID uuid.UUID, req dto.UpdateInningsStateRequest) (*dto.InningsStateResponse, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	meta, err := getInningsMeta(tx, inningsID)
	if err != nil {
		return nil, err
	}

	var state InningsState
	err = tx.Get(&state, `SELECT * FROM innings_state WHERE innings_id = $1 FOR UPDATE`, inningsID)
	stateExists := err == nil
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	finalStrikerID := ""
	finalNonStrikerID := ""
	finalBowlerID := ""
	if stateExists {
		if state.StrikerID != nil {
			finalStrikerID = *state.StrikerID
		}
		if state.NonStrikerID != nil {
			finalNonStrikerID = *state.NonStrikerID
		}
		if state.BowlerID != nil {
			finalBowlerID = *state.BowlerID
		}
	}
	if req.StrikerID != nil {
		finalStrikerID = req.StrikerID.String()
	}
	if req.NonStrikerID != nil {
		finalNonStrikerID = req.NonStrikerID.String()
	}
	if req.BowlerID != nil {
		finalBowlerID = req.BowlerID.String()
	}

	if finalStrikerID == "" || finalNonStrikerID == "" || finalBowlerID == "" {
		return nil, errors.New("striker_id, non_striker_id and bowler_id are required to set innings state")
	}
	if finalStrikerID == finalNonStrikerID {
		return nil, errors.New("striker and non_striker must be different players")
	}
	if err := validateActorTeams(tx, meta.BattingTeamID, meta.BowlingTeamID, finalStrikerID, finalNonStrikerID, finalBowlerID); err != nil {
		return nil, err
	}

	if stateExists {
		_, err = tx.Exec(`
			UPDATE innings_state
			SET striker_id = $2,
				non_striker_id = $3,
				bowler_id = $4,
				updated_at = NOW()
			WHERE innings_id = $1
		`, inningsID, finalStrikerID, finalNonStrikerID, finalBowlerID)
	} else {
		_, err = tx.Exec(`
			INSERT INTO innings_state (
				innings_id, striker_id, non_striker_id, bowler_id,
				total_runs, total_wickets, legal_balls, current_over, current_ball, status
			) VALUES ($1, $2, $3, $4, 0, 0, 0, 0, 0, 'live')
		`, inningsID, finalStrikerID, finalNonStrikerID, finalBowlerID)
	}
	if err != nil {
		return nil, err
	}

	var s InningsState
	if err := tx.Get(&s, `SELECT * FROM innings_state WHERE innings_id = $1`, inningsID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.InningsStateResponse{
		InningsID:    s.InningsID,
		StrikerID:    s.StrikerID,
		NonStrikerID: s.NonStrikerID,
		BowlerID:     s.BowlerID,
		TotalRuns:    s.TotalRuns,
		TotalWickets: s.TotalWickets,
		LegalBalls:   s.LegalBalls,
		CurrentOver:  s.CurrentOver,
		CurrentBall:  s.CurrentBall,
	}, nil
}

func (e *Engine) UndoLastBall(inningsID uuid.UUID) (*dto.InningsStateResponse, error) {
	var state InningsState
	if err := database.DB.Get(&state, `SELECT * FROM innings_state WHERE innings_id = $1`, inningsID); err != nil {
		return nil, err
	}
	if state.LegalBalls == 0 && state.TotalRuns == 0 && state.TotalWickets == 0 {
		return nil, errors.New("nothing to undo")
	}

	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var lastBallID string
	if err := tx.Get(&lastBallID, `
		SELECT id FROM ball_events
		WHERE innings_id = $1 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`, inningsID); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(`UPDATE ball_events SET is_deleted = TRUE, deleted_at = NOW() WHERE id = $1`, lastBallID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`DELETE FROM point_events WHERE ball_event_id = $1`, lastBallID); err != nil {
		return nil, err
	}

	var rebuilt struct {
		Runs    int `db:"runs"`
		Wickets int `db:"wickets"`
		Legal   int `db:"legal"`
	}
	if err := tx.Get(&rebuilt, `
		SELECT COALESCE(SUM(total_runs),0) AS runs,
			   COALESCE(SUM(CASE WHEN is_wicket THEN 1 ELSE 0 END),0) AS wickets,
			   COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball') THEN 1 ELSE 0 END),0) AS legal
		FROM ball_events
		WHERE innings_id = $1 AND is_deleted = FALSE
	`, inningsID); err != nil {
		return nil, err
	}
	over := rebuilt.Legal / 6
	ball := rebuilt.Legal % 6

	var last struct {
		Striker    *string `db:"striker_id"`
		NonStriker *string `db:"non_striker_id"`
		Bowler     *string `db:"bowler_id"`
	}
	_ = tx.Get(&last, `
		SELECT striker_id, non_striker_id, bowler_id
		FROM ball_events
		WHERE innings_id = $1 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`, inningsID)

	if _, err := tx.Exec(`
		UPDATE innings_state
		SET total_runs = $2,
			total_wickets = $3,
			legal_balls = $4,
			current_over = $5,
			current_ball = $6,
			striker_id = COALESCE($7, striker_id),
			non_striker_id = COALESCE($8, non_striker_id),
			bowler_id = COALESCE($9, bowler_id),
			status = 'live',
			updated_at = NOW()
		WHERE innings_id = $1
	`, inningsID, rebuilt.Runs, rebuilt.Wickets, rebuilt.Legal, over, ball, last.Striker, last.NonStriker, last.Bowler); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`UPDATE innings SET total_runs = $2, total_wickets = $3, status = 'live' WHERE id = $1`, inningsID, rebuilt.Runs, rebuilt.Wickets); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var out InningsState
	if err := database.DB.Get(&out, `SELECT * FROM innings_state WHERE innings_id = $1`, inningsID); err != nil {
		return nil, err
	}
	return &dto.InningsStateResponse{
		InningsID:    out.InningsID,
		StrikerID:    out.StrikerID,
		NonStrikerID: out.NonStrikerID,
		BowlerID:     out.BowlerID,
		TotalRuns:    out.TotalRuns,
		TotalWickets: out.TotalWickets,
		LegalBalls:   out.LegalBalls,
		CurrentOver:  out.CurrentOver,
		CurrentBall:  out.CurrentBall,
	}, nil
}

func strPtr(v string) *string { return &v }
