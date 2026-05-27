package dto

import "time"

type CreateMatchRequest struct {
	TeamAID         string `json:"team_a_id"`
	TeamBID         string `json:"team_b_id"`
	Location        string `json:"location"`
	MatchDate       string `json:"match_date" binding:"required"`
	OversPerInnings int    `json:"overs_per_innings" binding:"required,min=1"`
}

type GetMatchesQuery struct {
	Status     string `form:"status"`
	TeamID     string `form:"team_id"`
	HostUserID string `form:"host_user_id"`
	Search     string `form:"search"`
}

type MatchResponse struct {
	ID string `db:"id" json:"id"`

	HostUserID string `db:"host_user_id" json:"host_user_id"`
	CreatedBy  string `db:"created_by" json:"created_by"`

	TeamAID   *string `db:"team_a_id" json:"team_a_id"`
	TeamAName *string `db:"team_a_name" json:"team_a_name"`

	TeamBID   *string `db:"team_b_id" json:"team_b_id"`
	TeamBName *string `db:"team_b_name" json:"team_b_name"`

	TeamOneScore *string `db:"team_one_score" json:"team_one_score"`
	TeamTwoScore *string `db:"team_two_score" json:"team_two_score"`

	TeamAMatchTeamID *string `db:"team_a_match_team_id" json:"team_a_match_team_id"`
	TeamBMatchTeamID *string `db:"team_b_match_team_id" json:"team_b_match_team_id"`

	Location *string `db:"location" json:"location"`

	Status string `db:"status" json:"status"`

	MatchDate *time.Time `db:"match_date" json:"match_date"`

	OversPerInnings int `db:"overs_per_innings" json:"overs_per_innings"`

	TossDecision *string `db:"toss_decision" json:"toss_decision"`

	TossWinnerTeamID *string `db:"toss_winner_team_id" json:"toss_winner_team_id"`

	FirstPickTeamID *string `db:"first_pick_team_id" json:"first_pick_team_id"`

	WinnerTeamID *string `db:"winner_match_team_id" json:"winner_match_team_id"`
}

type SetFirstPickRequest struct {
	FirstPickTeamID string `json:"first_pick_team_id" binding:"required"`
}

type CompleteMatchRequest struct {
	WinnerMatchTeamID string `json:"winner_match_team_id"`
	WinnerTeamID      string `json:"winner_team_id"`
}

type InningsResponse struct {
	ID                 string  `db:"id" json:"id"`
	MatchID            string  `db:"match_id" json:"match_id"`
	InningsNo          int     `db:"innings_no" json:"innings_no"`
	BattingMatchTeamID *string `db:"batting_team_id" json:"batting_match_team_id"`
	BowlingMatchTeamID *string `db:"bowling_team_id" json:"bowling_match_team_id"`
	TotalRuns          int     `db:"total_runs" json:"total_runs"`
	TotalWickets       int     `db:"total_wickets" json:"total_wickets"`
	LegalBalls         int     `db:"legal_balls" json:"legal_balls"`
	CurrentOver        int     `db:"current_over" json:"current_over"`
	CurrentBall        int     `db:"current_ball" json:"current_ball"`
}

type ScorecardPlayerStats struct {
	MatchTeamPlayerID string  `db:"team_player_id" json:"match_team_player_id"`
	MatchTeamID       *string `db:"team_id" json:"match_team_id"`
	UserID            *string `db:"user_id" json:"user_id"`
	PlayerName        string  `db:"player_name" json:"player_name"`
	RunsScored        int     `db:"runs_scored" json:"runs_scored"`
	BallsFaced        int     `db:"balls_faced" json:"balls_faced"`
	Fours             int     `db:"fours" json:"fours"`
	Sixes             int     `db:"sixes" json:"sixes"`
	IsOut             bool    `db:"is_out" json:"is_out"`
	RunsConceded      int     `db:"runs_conceded" json:"runs_conceded"`
	WicketsTaken      int     `db:"wickets_taken" json:"wickets_taken"`
	OversBowled       float64 `db:"overs_bowled" json:"overs_bowled"`
	FantasyPoints     int     `db:"fantasy_points" json:"fantasy_points"`
}

type RecentBall struct {
	ID            string  `db:"id" json:"id"`
	InningsID     string  `db:"innings_id" json:"innings_id"`
	BallNo        int     `db:"ball_no" json:"ball_no"`
	DeliveryNo    int     `db:"delivery_no" json:"delivery_no"`
	BallType      string  `db:"ball_type" json:"ball_type"`
	TotalRuns     int     `db:"total_runs" json:"total_runs"`
	IsWicket      bool    `db:"is_wicket" json:"is_wicket"`
	StrikerID     *string `db:"striker_id" json:"striker_id"`
	NonStrikerID  *string `db:"non_striker_id" json:"non_striker_id"`
	BowlerID      *string `db:"bowler_id" json:"bowler_id"`
	DismissalType *string `db:"dismissal_type" json:"dismissal_type"`
}

type MatchScorecardResponse struct {
	Innings             []InningsResponse      `json:"innings"`
	Batting             []ScorecardPlayerStats `json:"batting"`
	Bowling             []ScorecardPlayerStats `json:"bowling"`
	RecentBalls         []RecentBall           `json:"recent_balls"`
	CurrentStrikerID    *string                `json:"current_striker_id"`
	CurrentNonStrikerID *string                `json:"current_non_striker_id"`
	CurrentBowlerID     *string                `json:"current_bowler_id"`
}

type MatchSquadPlayer struct {
	MatchTeamPlayerID string  `db:"team_player_id" json:"match_team_player_id"`
	MatchTeamID       string  `db:"team_id" json:"match_team_id"`
	UserID            *string `db:"user_id" json:"user_id"`
	PlayerName        string  `db:"player_name" json:"player_name"`
	PhoneNumber       *string `db:"phone_number" json:"phone_number"`
	IsPlayingXI       bool    `db:"is_playing_xi" json:"is_playing_xi"`
	IsCaptain         bool    `db:"is_captain" json:"is_captain"`
	IsUmpire          bool    `db:"is_umpire" json:"is_umpire"`
	BattingOrder      *int    `db:"batting_order" json:"batting_order"`
}

type UpdateLineupPlayer struct {
	MatchTeamPlayerID string `json:"match_team_player_id" binding:"required"`
	IsPlayingXI       bool   `json:"is_playing_xi"`
	IsCaptain         bool   `json:"is_captain"`
	IsUmpire          bool   `json:"is_umpire"`
	BattingOrder      *int   `json:"batting_order"`
}

type UpdateMatchLineupRequest struct {
	Players []UpdateLineupPlayer `json:"players" binding:"required,min=1"`
}
