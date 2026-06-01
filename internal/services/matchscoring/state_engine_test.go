package matchscoring

import (
	"CricketDuniya-Backend/internal/dto"
	"testing"

	"github.com/google/uuid"
)

func TestMapBallRequestWideNormalizesToAtLeastOneRun(t *testing.T) {
	req := dto.BallInputRequest{
		InningsID:     uuid.New(),
		MatchID:       uuid.New(),
		BallType:      "wide",
		Extras:        0,
		RunsOffBat:    0,
		TotalRuns:     0,
		IsWicket:      true,
		DismissalType: "run_out",
	}
	state := &InningsState{CurrentOver: 3, CurrentBall: 2}

	ballReq, legal := mapBallRequest(req, uuid.New().String(), uuid.New().String(), uuid.New().String(), state)

	if legal {
		t.Fatalf("wide delivery should not be legal")
	}
	if ballReq.Wides != 1 {
		t.Fatalf("expected wide count to normalize to 1, got %d", ballReq.Wides)
	}
	if ballReq.TotalRuns != 1 {
		t.Fatalf("expected total runs to normalize to 1, got %d", ballReq.TotalRuns)
	}
	if !ballReq.IsWicket {
		t.Fatalf("expected wicket flag to be preserved")
	}
	if ballReq.DismissalType != "run_out" {
		t.Fatalf("expected dismissal type to be preserved, got %q", ballReq.DismissalType)
	}
}

func TestMapBallRequestNoBallNormalizesExtrasAndSetsFreeHitState(t *testing.T) {
	req := dto.BallInputRequest{
		InningsID:  uuid.New(),
		MatchID:    uuid.New(),
		BallType:   "no_ball",
		Extras:     0,
		RunsOffBat: 4,
		TotalRuns:  0,
	}
	state := &InningsState{CurrentOver: 3, CurrentBall: 2, IsFreeHit: true}

	ballReq, legal := mapBallRequest(req, uuid.New().String(), uuid.New().String(), uuid.New().String(), state)

	if legal {
		t.Fatalf("no ball delivery should not be legal")
	}
	if ballReq.NoBalls != 1 {
		t.Fatalf("expected no balls to normalize to 1, got %d", ballReq.NoBalls)
	}
	if ballReq.Extras != 1 {
		t.Fatalf("expected extras to normalize to 1, got %d", ballReq.Extras)
	}
	if ballReq.TotalRuns != 5 {
		t.Fatalf("expected total runs to be 5, got %d", ballReq.TotalRuns)
	}
	if !ballReq.IsFreeHit {
		t.Fatalf("expected request to preserve free-hit state")
	}
}
