package skat

import (
	"math/rand"
	"time"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

var DB = []Game{
	{
		Id:      0,
		Started: time.Now().AddDate(0, 0, -1),
		Stake:   2,
		Ended:   nil,
		Players: []string{"Jannik", "Moritz", "Manuel", "Niklas"},
		Online:  false,
		Rounds: []Round{
			{
				Dealer:    "Jannik",
				Declarer:  "Moritz",
				Opponents: []string{"Niklas", "Manuel"},
				Won:       true,
				Value:     24,
			},
			{
				Dealer:    "Moritz",
				Declarer:  "Manuel",
				Opponents: []string{"Jannik", "Niklas"},
				Won:       false,
				Value:     24,
			},
		},
	},
}

func (s Service) List() []Game {
	return DB
}

func (s Service) Find(gameId int) *Game {
	for i, e := range DB {
		if gameId == e.Id {
			return &DB[i]
		}
	}
	return nil
}

func (s Service) Create(g Game) Game {
	g.Id = rand.Intn(100)
	g.Started = time.Now()
	g.Ended = nil

	DB = append(DB, g)
	return g
}
