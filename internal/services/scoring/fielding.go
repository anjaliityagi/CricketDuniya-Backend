package scoring

import "CricketDuniya-Backend/internal/dto"

func fielding(req dto.BallRequest) int {
	points := 0

	if req.FielderID != nil && req.IsWicket {
		points += 2
	}

	if req.DismissedPlayerID != nil && req.IsWicket {
		points += 5
	}

	return points
}
