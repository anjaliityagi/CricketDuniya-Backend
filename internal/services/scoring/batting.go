package scoring

import "CricketDuniya-Backend/internal/dto"

func batting(req dto.BallRequest) int {
	points := 0

	// runs
	switch req.RunsOffBat {
	case 1:
		points += 1
	case 2:
		points += 3
	case 4:
		points += 5
	case 6:
		points += 7
	}

	// boundaries
	if req.IsBoundaryFour {
		points += 5
	}
	if req.IsBoundarySix {
		points += 7
	}

	// dot ball penalty
	if req.IsDotBall {
		points -= 1
	}

	// wicket penalty
	if req.IsWicket {
		points -= 3
	}

	return points
}
