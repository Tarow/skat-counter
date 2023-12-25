package skat

import (
	"slices"

	"github.com/tarow/skat-counter/internal/skat/gen/model"
)

type Game struct {
	model.Game
	Players []model.Player
	Rounds  []Round
}

type Round struct {
	model.Round
	Opponents []int32 `alias:"opponent.id"`
}

func (g Game) GetDate() string {
	return g.CreatedAt.Format("Monday, 02.01.2006")
}

func (g Game) GetTotalPlayerScore(player model.Player) int32 {
	sum := int32(0)
	for _, r := range g.Rounds {
		roundScore := r.GetRoundScore(player)
		if roundScore != nil {
			sum += *roundScore
		}
	}
	return sum
}

func (g Game) GetTotalPayment() float32 {
	sum := int32(0)
	for _, player := range g.Players {
		sum += g.GetTotalPlayerScore(player)
	}

	return float32(sum) * g.Stake
}

func (r Round) GetRoundScore(player model.Player) *int32 {
	if player.ID == r.Declarer {
		if r.Won {
			return intPtr(0)
		} else {
			return intPtr(2 * r.Value)
		}
	}

	if slices.Contains(r.Opponents, player.ID) {
		if r.Won {
			return &r.Value
		} else {
			return intPtr(0)
		}
	}

	return nil
}

func intPtr(i int32) *int32 {
	return &i
}
