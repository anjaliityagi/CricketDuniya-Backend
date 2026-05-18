package models

import "time"

type Innings struct {
	ID            string     `json:"id"`
	MatchID       string     `json:"match_id"`
	BattingTeamID string     `json:"batting_team_id"`
	BowlingTeamID string     `json:"bowling_team_id"`
	InningsNumber int        `json:"innings_number"`
	TargetRuns    *int       `json:"target_runs"`
	TotalRuns     int        `json:"total_runs"`
	Wickets       int        `json:"wickets"`
	TotalOvers    float64    `json:"total_overs"`
	StartedAt     *time.Time `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
}
