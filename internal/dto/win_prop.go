package dto

type WinProbabilityResponse struct {
	MatchID                string             `json:"match_id"`
	BattingTeamProbability float64            `json:"batting_team_probability"`
	BowlingTeamProbability float64            `json:"bowling_team_probability"`
	Innings                int                `json:"innings"`
	CalculatedFrom         ProbabilityFactors `json:"calculated_from"`
}

type ProbabilityFactors struct {
	RunsRequired     int `json:"runs_required"`
	BallsRemaining   int `json:"balls_remaining"`
	WicketsRemaining int `json:"wickets_remaining"`
}
