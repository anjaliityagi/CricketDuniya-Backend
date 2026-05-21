package dto

import "time"

type UpdateProfileRequest struct {
	Name         *string `json:"name" binding:"omitempty,min=1,max=100"`
	BattingStyle *string `json:"batting_style"`
	BowlingStyle *string `json:"bowling_style"`
}

type UserProfileResponse struct {
	User     UserProfileUser    `json:"user"`
	Summary  UserProfileSummary `json:"summary"`
	Batting  UserBattingStats   `json:"batting"`
	Bowling  UserBowlingStats   `json:"bowling"`
	Fielding UserFieldingStats  `json:"fielding"`
}

type UserProfileUser struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	PhoneNumber  string    `db:"phone_number" json:"phone_number"`
	BattingStyle *string   `db:"batting_style" json:"batting_style"`
	BowlingStyle *string   `db:"bowling_style" json:"bowling_style"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type UserProfileSummary struct {
	MatchesPlayed int     `db:"matches_played" json:"matches_played"`
	Won           int     `db:"won" json:"won"`
	Lost          int     `db:"lost" json:"lost"`
	MVPs          int     `db:"mvps" json:"mvps"`
	WinPercentage float64 `db:"win_percentage" json:"win_percentage"`
	Points        int     `db:"points" json:"points"`
}

type UserBattingStats struct {
	Average    float64 `db:"average" json:"average"`
	StrikeRate float64 `db:"strike_rate" json:"strike_rate"`
	HighScore  int     `db:"high_score" json:"high_score"`
	Runs       int     `db:"runs" json:"runs"`
	Innings    int     `db:"innings" json:"innings"`
	Fours      int     `db:"fours" json:"fours"`
	Sixes      int     `db:"sixes" json:"sixes"`
}

type UserBowlingStats struct {
	OversBowled  float64 `db:"overs_bowled" json:"overs_bowled"`
	Wickets      int     `db:"wickets" json:"wickets"`
	RunsConceded int     `db:"runs_conceded" json:"runs_conceded"`
	Maidens      int     `db:"maidens" json:"maidens"`
	Economy      float64 `db:"economy" json:"economy"`
}

type UserFieldingStats struct {
	Catches  int `db:"catches" json:"catches"`
	Stumping int `db:"stumping" json:"stumping"`
	RunOuts  int `db:"run_outs" json:"run_outs"`
}
