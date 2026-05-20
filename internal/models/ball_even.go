package models

type BallEvent struct {
	ID        string
	MatchID   string
	InningsID string

	Over int
	Ball int

	BatsmanID string
	BowlerID  string

	Runs       int
	ExtraType  string
	IsWicket   bool
	WicketType string
}
