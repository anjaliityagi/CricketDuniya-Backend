package dto

import "github.com/google/uuid"

type BallRequest struct {
	InningsID uuid.UUID `json:"innings_id"`
	MatchID   uuid.UUID `json:"match_id"`

	StrikerID    uuid.UUID `json:"striker_id"`
	NonStrikerID uuid.UUID `json:"non_striker_id"`
	BowlerID     uuid.UUID `json:"bowler_id"`

	BallNo     int    `json:"ball_no"`
	DeliveryNo int    `json:"delivery_no"`
	BallType   string `json:"ball_type"`

	RunsScored int  `json:"runs_scored"`
	RunsOffBat int  `json:"runs_off_bat"`
	Extras     int  `json:"extras"`
	TotalRuns  int  `json:"total_runs"`
	IsFreeHit  bool `json:"is_free_hit"`

	IsDotBall bool `json:"is_dot_ball"`

	IsBoundaryFour bool `json:"is_boundary_four"`
	IsBoundarySix  bool `json:"is_boundary_six"`

	IsWicket      bool   `json:"is_wicket"`
	DismissalType string `json:"dismissal_type"`

	DismissedPlayerID *uuid.UUID `json:"dismissed_player_id"`
	FielderID         *uuid.UUID `json:"fielder_id"`

	Wides   int `json:"wides"`
	NoBalls int `json:"no_balls"`
	Byes    int `json:"byes"`
	LegByes int `json:"leg_byes"`
}

type BallInputRequest struct {
	InningsID uuid.UUID `json:"innings_id" binding:"required"`
	MatchID   uuid.UUID `json:"match_id" binding:"required"`

	StrikerID    *uuid.UUID `json:"striker_id"`
	NonStrikerID *uuid.UUID `json:"non_striker_id"`
	BowlerID     *uuid.UUID `json:"bowler_id"`

	BallType   string `json:"ball_type" binding:"required"`
	RunsOffBat int    `json:"runs_off_bat"`
	Extras     int    `json:"extras"`
	TotalRuns  int    `json:"total_runs"`

	IsWicket          bool       `json:"is_wicket"`
	DismissalType     string     `json:"dismissal_type"`
	DismissedPlayerID *uuid.UUID `json:"dismissed_player_id"`
	FielderID         *uuid.UUID `json:"fielder_id"`
	NextBatterID      *uuid.UUID `json:"next_batter_id"`
}

type InningsStateResponse struct {
	InningsID         string  `json:"innings_id"`
	StrikerID         *string `json:"striker_id"`
	NonStrikerID      *string `json:"non_striker_id"`
	BowlerID          *string `json:"bowler_id"`
	IsFreeHit         bool    `json:"is_free_hit"`
	TotalRuns         int     `json:"total_runs"`
	TotalWickets      int     `json:"total_wickets"`
	LegalBalls        int     `json:"legal_balls"`
	CurrentOver       int     `json:"current_over"`
	CurrentBall       int     `json:"current_ball"`
	OverCompleted     bool    `json:"over_completed"`
	InningsCompleted  bool    `json:"innings_completed"`
	NeedsNextBowler   bool    `json:"needs_next_bowler"`
	NeedsNextBatter   bool    `json:"needs_next_batter"`
	RequiredRunsToWin *int    `json:"required_runs_to_win,omitempty"`
	BallsRemaining    *int    `json:"balls_remaining,omitempty"`
}

type UpdateInningsStateRequest struct {
	StrikerID    *uuid.UUID `json:"striker_id"`
	NonStrikerID *uuid.UUID `json:"non_striker_id"`
	BowlerID     *uuid.UUID `json:"bowler_id"`
}
