package models

import "time"

type Match struct {
	ID             string     `db:"id" json:"id"`
	HostUserID     string     `db:"host_user_id" json:"host_user_id"`
	OpponentUserID *string    `db:"opponent_user_id" json:"opponent_user_id,omitempty"`
	Venue          *string    `db:"venue" json:"venue,omitempty"`
	MatchDate      *time.Time `db:"match_date" json:"match_date,omitempty"`
	OversPerSide   int        `db:"overs_per_side" json:"overs_per_side"`
	TossWinnerSide *string    `db:"toss_winner_side" json:"toss_winner_side,omitempty"`
	TossDecision   *string    `db:"toss_decision" json:"toss_decision,omitempty"`
	WinnerSide     *string    `db:"winner_side" json:"winner_side,omitempty"`
	Status         string     `db:"status" json:"status"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
}
