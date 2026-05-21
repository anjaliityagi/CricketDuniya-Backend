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
	BallType   string `json:"ball_type"` // normal, wide, no_ball, etc

	RunsScored int `json:"runs_scored"`
	RunsOffBat int `json:"runs_off_bat"`
	Extras     int `json:"extras"`
	TotalRuns  int `json:"total_runs"`

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
