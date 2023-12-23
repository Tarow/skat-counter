package skat

import (
	"slices"
	"time"
)

type Game struct {
	Id      int
	Started time.Time
	Ended   *time.Time
	Stake   float32  `form:"stake"`
	Players []string `form:"player"`
	Rounds  []Round  `form:"round"`
	Online  bool     `form:"online"`
}

type Round struct {
	Dealer    string
	Declarer  string
	Opponents []string
	Won       bool
	Value     int
}

func (g Game) IsActive() bool {
	return g.Ended == nil || g.Ended.After(time.Now())
}

func (g Game) GetDate() string {
	return g.Started.Format("Monday, 02.01.2006")
}

func (g Game) GetTotalPlayerScore(player string) int {
	sum := 0
	for _, r := range g.Rounds {
		roundScore := r.GetRoundScore(player)
		if roundScore != nil {
			sum += *roundScore
		}
	}
	return sum
}

func (g Game) GetTotalPayment() float32 {
	sum := 0
	for _, player := range g.Players {
		sum += g.GetTotalPlayerScore(player)
	}

	return float32(sum) * g.Stake
}

func (r Round) GetRoundScore(player string) *int {
	if player == r.Declarer {
		if r.Won {
			return intPtr(0)
		} else {
			return intPtr(2 * r.Value)
		}
	}

	if slices.Contains(r.Opponents, player) {
		if r.Won {
			return &r.Value
		} else {
			return intPtr(0)
		}
	}

	return nil
}

func intPtr(i int) *int {
	return &i
}
