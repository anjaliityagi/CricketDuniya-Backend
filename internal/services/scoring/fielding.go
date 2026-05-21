package scoring

import "CricketDuniya-Backend/internal/dto"

func fielding(req dto.BallRequest) int {
	points := 0

	if req.FielderID != nil && req.IsWicket && req.DismissalType == "caught" {
		points += 2
	}

	if req.FielderID != nil && req.IsWicket && (req.DismissalType == "run_out" || req.DismissalType == "stumped") {
		points += 5
	}

	return points
}
