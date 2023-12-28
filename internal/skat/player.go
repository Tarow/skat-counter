package skat

import (
	"github.com/tarow/skat-counter/internal/skat/gen/model"
	"github.com/tarow/skat-counter/internal/skat/gen/table"
)

func (s Service) ListPlayers() ([]model.Player, error) {
	stmt := table.Player.SELECT(table.Player.AllColumns).FROM(table.Player)

	var result []model.Player
	err := stmt.Query(s.db, &result)

	return result, err
}
