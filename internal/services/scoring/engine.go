package scoring

import "CricketDuniya-Backend/internal/dto"

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Process(req dto.BallRequest) (bat, bowl, field int) {
	bat = batting(req)
	bowl = bowling(req)
	field = fielding(req)
	return
}
