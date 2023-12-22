package skat

import (
	"fmt"
	"time"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

var db = []Game{
	{
		Id:      0,
		Active:  false,
		Started: time.Now().AddDate(0, 0, -1),
		Ended:   nil,
		Players: []string{"Jannik", "Moritz", "Manuel", "Niklas"},
		Online:  false,
		Rounds: []Round{
			{
				Dealer:    "Niklas",
				Declarer:  "Moritz",
				Opponents: []string{"Jannik", "Manuel"},
				Won:       true,
			},
			{
				Dealer:    "Niklas",
				Declarer:  "Moritz",
				Opponents: []string{"Jannik", "Manuel"},
				Won:       false,
			},
		},
	},
	{
		Id:      1,
		Active:  true,
		Started: time.Now(),
		Ended:   nil,
		Players: []string{"Jannik", "Moritz", "Manuel", "Niklas"},
		Online:  true,
		Rounds: []Round{
			{
				Dealer:    "Niklas",
				Declarer:  "Moritz",
				Opponents: []string{"Jannik", "Manuel"},
				Won:       true,
			},
			{
				Dealer:    "Niklas",
				Declarer:  "Moritz",
				Opponents: []string{"Jannik", "Manuel"},
				Won:       false,
			},
		},
	},
}

func (s Service) List() []Game {
	return db
}

func (s Service) Find(gameId int) (Game, error) {
	for _, e := range db {
		if gameId == e.Id {
			return e, nil
		}
	}
	return Game{}, fmt.Errorf("no game found with id %v", gameId)
}
