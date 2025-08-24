package game

import (
	"errors"

	"github.com/ech21/fast-food-app/internal/types"
)

type Player struct {
	dead      bool         // Dead players are ready for cleanup
	Name      string       `json:"name"`
	FoodEaten []types.Item `json:"foodEaten"`
}

func NewPlayer(name string) (p *Player, err error) {
	if len(name) == 0 {
		return nil, errors.New("Name is empty")
	}
	return &Player{
		dead:      false,
		Name:      name,
		FoodEaten: []types.Item{},
	}, nil
}

func (p *Player) Eat(item types.Item) {
	p.FoodEaten = append(p.FoodEaten, item)
}
