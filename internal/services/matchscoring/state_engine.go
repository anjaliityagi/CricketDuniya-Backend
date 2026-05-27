package matchscoring

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
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
	noBatterLeft := false
	needsReplacement := req.IsWicket || strings.EqualFold(req.DismissalType, "retired_hurt")

	if needsReplacement && req.DismissedPlayerID != nil {
		dismissedID := req.DismissedPlayerID.String()
		if req.NextBatterID == nil && dismissedID == nextStriker && dismissedID == nextNonStriker {
			noBatterLeft = true
		}
		if dismissedID == nextStriker {
			if req.NextBatterID == nil {
				nextStriker = nextNonStriker
			} else {
				nextStriker = req.NextBatterID.String()
				needsNextBatter = true
			}
		}
		if dismissedID == nextNonStriker {
			if req.NextBatterID == nil {
				nextNonStriker = nextStriker
			} else {
				nextNonStriker = req.NextBatterID.String()
				needsNextBatter = true
			}
		}
	}

	totalRunsForStrike := internalReq.TotalRuns
	strikeChangesOnBall := internalReq.BallType != "wide" && internalReq.BallType != "no_ball"
	if strikeChangesOnBall && totalRunsForStrike%2 == 1 {
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
	if noBatterLeft {
		inningsCompleted = true
	}

	if inningsMeta.TargetRuns != nil && inningsMeta.InningsNo == 2 && matchUpdate.InningsRuns >= *inningsMeta.TargetRuns {
		inningsCompleted = true
	}

	status := "live"
	if inningsCompleted {
		status = "completed"
	}

	if err := repositories.UpdateInningsStateAfterBallTx(tx, req.InningsID, nextStriker, nextNonStriker, nextBowler, matchUpdate.InningsRuns, matchUpdate.InningsWickets, nextLegalBalls, nextOver, nextBall, status); err != nil {
		return nil, err
	}

	if err := repositories.UpdateInningsTotalsStatusTx(tx, req.InningsID, matchUpdate.InningsRuns, matchUpdate.InningsWickets, status); err != nil {
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

	if inningsCompleted && inningsMeta.InningsNo >= 2 {
		if err := autoCompleteMatchIfNeeded(req.MatchID.String(), inningsMeta, matchUpdate.InningsRuns); err != nil {
			return nil, err
		}
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

type inningsMeta = repositories.InningsMeta

func getInningsMeta(tx *sqlx.Tx, inningsID uuid.UUID) (*inningsMeta, error) {
	return repositories.GetInningsMetaTx(tx, inningsID)
}

func ensureNextInningsTx(tx *sqlx.Tx, meta *inningsMeta, firstInningsRuns int) error {
	if meta.InningsNo != 1 {
		// Super over: after first innings of each super over, create the chase innings.
		if meta.IsSuperOver && meta.InningsNo%2 == 1 {
			nextInningsNo := meta.InningsNo + 1
			count, err := repositories.CountInningsByNoTx(tx, meta.MatchID, nextInningsNo)
			if err != nil {
				return err
			}
			if count > 0 {
				return nil
			}
			target := firstInningsRuns + 1
			if err := repositories.InsertInningsTx(tx, meta.MatchID, meta.BowlingTeamID, meta.BattingTeamID, nextInningsNo, &target); err != nil {
				return err
			}
			inningsID, err := repositories.GetInningsIDByNoTx(tx, meta.MatchID, nextInningsNo)
			if err != nil {
				return err
			}
			superOverNo := 1
			if meta.SuperOverNo != nil {
				superOverNo = *meta.SuperOverNo
			}
			return repositories.LinkSuperOverToInningsTx(tx, meta.MatchID, inningsID, meta.BowlingTeamID, meta.BattingTeamID, superOverNo)
		}
		return nil
	}

	inningsCount, err := repositories.CountSecondInningsTx(tx, meta.MatchID)
	if err != nil {
		return err
	}
	if inningsCount > 0 {
		return nil
	}

	targetRuns := firstInningsRuns + 1
	return repositories.InsertSecondInningsTx(tx, meta.MatchID, meta.BowlingTeamID, meta.BattingTeamID, targetRuns)
}

func getPreviousOverBowlerID(tx *sqlx.Tx, inningsID uuid.UUID) (string, error) {
	return repositories.GetPreviousOverBowlerIDTx(tx, inningsID)
}

func getOrCreateState(tx *sqlx.Tx, req dto.BallInputRequest, meta *inningsMeta) (*InningsState, error) {
	var state InningsState
	err := repositories.GetInningsStateForUpdateTx(tx, req.InningsID, &state)
	if err == nil {
		return &state, nil
	}

	if req.StrikerID == nil || req.NonStrikerID == nil || req.BowlerID == nil {
		return nil, errors.New("striker_id, non_striker_id and bowler_id are required for first ball")
	}

	if err := repositories.InsertInitialInningsStateTx(tx, req); err != nil {
		return nil, err
	}

	if err := repositories.GetInningsStateForUpdateTx(tx, req.InningsID, &state); err != nil {
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
	c, err := repositories.CountStrikersOnTeamTx(tx, strikerID, nonStrikerID, battingTeamID)
	if err != nil {
		return err
	}
	if c < 1 {
		return errors.New("striker/non_striker must belong to batting team")
	}
	c, err = repositories.CountBowlerOnTeamTx(tx, bowlerID, bowlingTeamID)
	if err != nil {
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
	case "dead_ball", "retired_hurt":
		legal = false
		totalRuns = 0
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
	return repositories.SaveBallEventTx(tx, req)
}

func persistFantasyTx(tx *sqlx.Tx, req dto.BallRequest, ballEventID string, battingPoints, bowlingPoints, fieldingPoints int) error {

	if req.StrikerID.String() != "" && battingPoints != 0 {
		
		if err := repositories.UpsertFantasyPointsTx(tx, req.MatchID.String(), req.StrikerID.String(), battingPoints, "batting_points"); err != nil {
			return err
		}
		if err := repositories.InsertPointEventTx(tx, req.MatchID.String(), req.StrikerID.String(), ballEventID, "batting", "ball_batting_points", battingPoints); err != nil {
			return err
		}
	}
	if req.BowlerID.String() != "" && bowlingPoints != 0 {
		if err := repositories.UpsertFantasyPointsTx(tx, req.MatchID.String(), req.BowlerID.String(), bowlingPoints, "bowling_points"); err != nil {
			return err
		}
		if err := repositories.InsertPointEventTx(tx, req.MatchID.String(), req.BowlerID.String(), ballEventID, "bowling", "ball_bowling_points", bowlingPoints); err != nil {
			return err
		}
	}
	if req.FielderID != nil && req.FielderID.String() != "" && fieldingPoints != 0 {
		if err := repositories.UpsertFantasyPointsTx(tx, req.MatchID.String(), req.FielderID.String(), fieldingPoints, "fielding_points"); err != nil {
			return err
		}
		if err := repositories.InsertPointEventTx(tx, req.MatchID.String(), req.FielderID.String(), ballEventID, "fielding", "ball_fielding_points", fieldingPoints); err != nil {
			return err
		}
	}
	return nil
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
	err = repositories.GetInningsStateForUpdateTx(tx, inningsID, &state)
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
	if err := validateActorTeams(tx, meta.BattingTeamID, meta.BowlingTeamID, finalStrikerID, finalNonStrikerID, finalBowlerID); err != nil {
		return nil, err
	}

	if stateExists {
		err = repositories.UpdateInningsStateAfterBallTx(tx, inningsID, finalStrikerID, finalNonStrikerID, finalBowlerID, state.TotalRuns, state.TotalWickets, state.LegalBalls, state.CurrentOver, state.CurrentBall, state.Status)
	} else {
		insertReq := dto.BallInputRequest{InningsID: inningsID, StrikerID: ptrUUID(finalStrikerID), NonStrikerID: ptrUUID(finalNonStrikerID), BowlerID: ptrUUID(finalBowlerID)}
		err = repositories.InsertInitialInningsStateTx(tx, insertReq)
	}
	if err != nil {
		return nil, err
	}

	var s InningsState
	if err := repositories.GetInningsStateTx(tx, inningsID, &s); err != nil {
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
	if err := repositories.GetInningsStateForUpdateTx(tx, inningsID, &state); err != nil {
		return nil, err
	}

	if err := ensureUndoAllowedTx(tx, meta.MatchID, meta.InningsNo); err != nil {
		return nil, err
	}

	deletedBall, err := repositories.GetLastUndoBallTx(tx, inningsID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("nothing to undo")
		}
		return nil, err
	}

	if err := repositories.SoftDeleteBallTx(tx, deletedBall.ID); err != nil {
		return nil, err
	}
	if err := repositories.DeletePointEventsByBallTx(tx, deletedBall.ID); err != nil {
		return nil, err
	}

	if err := rebuildPlayerMatchStatsTx(tx, meta.MatchID); err != nil {
		return nil, err
	}

	rebuilt, err := getRebuiltInningsTotalsTx(tx, inningsID)
	if err != nil {
		return nil, err
	}
	over := rebuilt.Legal / 6
	ball := rebuilt.Legal % 6

	if err := repositories.UpdateInningsStateAfterUndoTx(tx, inningsID, deletedBall.StrikerID, deletedBall.NonStriker, deletedBall.BowlerID, rebuilt.Runs, rebuilt.Wickets, rebuilt.Legal, over, ball); err != nil {
		return nil, err
	}
	if err := repositories.ReopenInningsTx(tx, inningsID, rebuilt.Runs, rebuilt.Wickets); err != nil {
		return nil, err
	}

	if err := cleanupFutureInningsTx(tx, meta.MatchID, meta.InningsNo); err != nil {
		return nil, err
	}
	if err := reopenMatchAfterUndoTx(tx, meta.MatchID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	resp := &dto.InningsStateResponse{
		InningsID:    state.InningsID,
		StrikerID:    deletedBall.StrikerID,
		NonStrikerID: deletedBall.NonStriker,
		BowlerID:     deletedBall.BowlerID,
		TotalRuns:    rebuilt.Runs,
		TotalWickets: rebuilt.Wickets,
		LegalBalls:   rebuilt.Legal,
		CurrentOver:  over,
		CurrentBall:  ball,
	}
	if meta.TargetRuns != nil && meta.InningsNo == 2 {
		reqRuns := *meta.TargetRuns - rebuilt.Runs
		if reqRuns < 0 {
			reqRuns = 0
		}
		ballsLeft := meta.OversPerInnings*6 - rebuilt.Legal
		if ballsLeft < 0 {
			ballsLeft = 0
		}
		resp.RequiredRunsToWin = &reqRuns
		resp.BallsRemaining = &ballsLeft
	}

	return resp, nil
}

type rebuiltInningsTotals struct {
	Runs    int `db:"runs"`
	Wickets int `db:"wickets"`
	Legal   int `db:"legal"`
}

func getRebuiltInningsTotalsTx(tx *sqlx.Tx, inningsID uuid.UUID) (*rebuiltInningsTotals, error) {
	rebuilt, err := repositories.GetRebuiltInningsTotalsTx(tx, inningsID)
	if err != nil {
		return nil, err
	}
	return &rebuiltInningsTotals{Runs: rebuilt.Runs, Wickets: rebuilt.Wickets, Legal: rebuilt.Legal}, nil
}

func ensureUndoAllowedTx(tx *sqlx.Tx, matchID string, inningsNo int) error {
	futureBallCount, err := repositories.CountFutureBallsTx(tx, matchID, inningsNo)
	if err != nil {
		return err
	}
	if futureBallCount > 0 {
		return errors.New("cannot undo this innings after next innings has started")
	}
	return nil
}

func cleanupFutureInningsTx(tx *sqlx.Tx, matchID string, inningsNo int) error {
	if err := repositories.DeleteFutureInningsStateTx(tx, matchID, inningsNo); err != nil {
		return err
	}
	return repositories.DeleteFutureInningsTx(tx, matchID, inningsNo)
}

func reopenMatchAfterUndoTx(tx *sqlx.Tx, matchID string) error {
	return repositories.ReopenMatchAfterUndoTx(tx, matchID)
}

func rebuildPlayerMatchStatsTx(tx *sqlx.Tx, matchID string) error {
	return repositories.RebuildPlayerMatchStatsTx(tx, matchID)
}

func autoCompleteMatchIfNeeded(matchID string, meta *inningsMeta, secondInningsRuns int) error {
	status, err := repositories.GetMatchStatus(matchID)
	if err != nil {
		return err
	}
	if status == "completed" {
		return nil
	}

	winnerMatchTeamID, err := determineWinnerMatchTeamID(matchID, meta, secondInningsRuns)
	if err != nil {
		return err
	}
	if winnerMatchTeamID == "" {
		return nil
	}

	inningsIDs, err := repositories.GetMatchInningsIDs(matchID)
	if err != nil {
		return err
	}
	for _, matchInningsID := range inningsIDs {
		if err := scoring.ApplyNotOutBonus(matchID, matchInningsID); err != nil {
			return err
		}
	}
	if err := scoring.ApplyResultPoints(matchID, winnerMatchTeamID); err != nil {
		return err
	}
	return repositories.FinalizeMatch(matchID, winnerMatchTeamID)
}

func determineWinnerMatchTeamID(matchID string, meta *inningsMeta, secondInningsRuns int) (string, error) {
	// Regular match innings 2 tie triggers super over flow.
	if meta.InningsNo == 2 {
		firstInningsRuns, err := repositories.GetFirstInningsRuns(matchID)
		if err != nil {
			return "", err
		}
		if secondInningsRuns > firstInningsRuns {
			return meta.BattingTeamID, nil
		}
		if secondInningsRuns < firstInningsRuns {
			return meta.BowlingTeamID, nil
		}
		// Tie: create super over 1 (two innings: 3 and 4), do not finalize match yet.
		if err := createNextSuperOver(meta.MatchID, 1); err != nil {
			return "", err
		}
		return "", nil
	}

	// Super over completion check: only decide after the second innings of that super over.
	if meta.IsSuperOver {
		if meta.InningsNo%2 == 1 {
			return "", nil
		}
		superOverNo := 1
		if meta.SuperOverNo != nil {
			superOverNo = *meta.SuperOverNo
		}
		firstSOInningsNo := meta.InningsNo - 1
		firstSORuns, err := repositories.GetInningsRunsByNo(meta.MatchID, firstSOInningsNo)
		if err != nil {
			return "", err
		}
		if secondInningsRuns > firstSORuns {
			return meta.BattingTeamID, nil
		}
		if secondInningsRuns < firstSORuns {
			return meta.BowlingTeamID, nil
		}
		// Tied super over: keep creating next super over until winner is found.
		if err := createNextSuperOver(meta.MatchID, superOverNo+1); err != nil {
			return "", err
		}
		return "", nil
	}

	if meta.TargetRuns != nil {
		if secondInningsRuns >= *meta.TargetRuns {
			return meta.BattingTeamID, nil
		}
		return meta.BowlingTeamID, nil
	}

	firstInningsRuns, err := repositories.GetFirstInningsRuns(matchID)
	if err != nil {
		return "", err
	}
	if secondInningsRuns > firstInningsRuns {
		return meta.BattingTeamID, nil
	}
	if secondInningsRuns < firstInningsRuns {
		return meta.BowlingTeamID, nil
	}
	return "", nil
}

func createNextSuperOver(matchID string, superOverNo int) error {
	tx, err := repositories.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	matchTeams, err := repositories.GetMatchTeamsTx(tx, matchID)
	if err != nil {
		return err
	}

	// super over innings numbers: 3/4 for SO1, 5/6 for SO2
	firstInningsNo := 3 + (superOverNo-1)*2

	count, err := repositories.CountInningsByNoTx(tx, matchID, firstInningsNo)
	if err != nil {
		return err
	}
	if count > 0 {
		return tx.Commit()
	}

	if err := repositories.InsertInningsTx(tx, matchID, matchTeams.TeamAID, matchTeams.TeamBID, firstInningsNo, nil); err != nil {
		return err
	}
	inningsID, err := repositories.GetInningsIDByNoTx(tx, matchID, firstInningsNo)
	if err != nil {
		return err
	}
	if err := repositories.LinkSuperOverToInningsTx(tx, matchID, inningsID, matchTeams.TeamAID, matchTeams.TeamBID, superOverNo); err != nil {
		return err
	}

	return tx.Commit()
}

func strPtr(v string) *string { return &v }

func ptrUUID(id string) *uuid.UUID {
	v := uuid.MustParse(id)
	return &v
}
