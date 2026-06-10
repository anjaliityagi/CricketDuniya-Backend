package models

import "time"

type Innings struct {
	ID            string     `json:"id"`
	MatchID       string     `json:"match_id" db:"match_id"`
	BattingTeamID string     `json:"batting_team_id" db:"batting_team_id"`
	BowlingTeamID string     `json:"bowling_team_id" db:"bowling_team_id"`
	InningsNumber int        `json:"innings_number" db:"innings_no"`
	TargetRuns    *int       `json:"target_runs" db:"target_runs"`
	TotalRuns     int        `json:"total_runs" db:"total_runs"`
	Wickets       int        `json:"wickets" db:"total_wickets"`
	TotalOvers    float64    `json:"total_overs" db:"total_overs"`
	StartedAt     *time.Time `json:"started_at" db:"started_at"`
	EndedAt       *time.Time `json:"ended_at" db:"completed_at"`
	Status        string     `json:"status" db:"status"`
}
