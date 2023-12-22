package skat

import "time"

type Game struct {
	Id      int
	Active  bool
	Started time.Time
	Ended   *time.Time
	Players []string
	Rounds  []Round
	Online  bool
}

type Round struct {
	Dealer    string
	Declarer  string
	Opponents []string
	Won       bool
	Points    uint
}

func (g Game) GetDate() string {
	return g.Started.Format("Monday, 02.01.2006")
}
