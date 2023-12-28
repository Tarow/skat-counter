package skat

import (
	"github.com/go-jet/jet/v2/sqlite"
	"github.com/tarow/skat-counter/internal/skat/gen/model"
	"github.com/tarow/skat-counter/internal/skat/gen/table"
)

func (s Service) AddRound(round Round) (Round, error) {
	opponents := round.Opponents

	tx, err := s.db.Begin()
	if err != nil {
		return Round{}, err
	}
	defer tx.Rollback()

	// Insert round
	stmt := table.Round.INSERT(
		table.Round.AllColumns.Except(table.Round.ID),
	).MODEL(round).RETURNING(table.Round.AllColumns)

	err = stmt.Query(tx, &round)
	if err != nil {
		return Round{}, err
	}
	round.Opponents = opponents

	roundOpponents := make([]model.RoundOpponent, 0)
	for _, opponentId := range opponents {
		roundOpponents = append(roundOpponents, model.RoundOpponent{
			RoundID:  round.ID,
			PlayerID: opponentId,
		})
	}
	stmt = table.RoundOpponent.INSERT(
		table.RoundOpponent.AllColumns,
	).MODELS(roundOpponents)

	_, err = stmt.Exec(tx)
	if err != nil {
		return Round{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Round{}, err
	}

	return round, nil
}

func (s Service) UpdateRound(round Round) (Round, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Round{}, err
	}
	defer tx.Rollback()

	// Delete old opponents
	deleteOpponents := table.RoundOpponent.DELETE().
		WHERE(table.RoundOpponent.RoundID.EQ(sqlite.Int32(round.ID)))
	_, err = deleteOpponents.Exec(tx)
	if err != nil {
		return Round{}, err
	}

	// Update round
	updateRound := table.Round.UPDATE(
		table.Round.AllColumns.Except(table.Round.ID, table.Round.CreatedAt),
	).WHERE(table.Round.ID.EQ(sqlite.Int32(round.ID))).
		MODEL(round)
	_, err = updateRound.Exec(tx)
	if err != nil {
		return Round{}, err
	}

	// Insert new round opponents
	roundOpponents := make([]model.RoundOpponent, 0)
	for _, opponentId := range round.Opponents {
		roundOpponents = append(roundOpponents, model.RoundOpponent{
			RoundID:  round.ID,
			PlayerID: opponentId,
		})
	}
	insertOpponents := table.RoundOpponent.INSERT(
		table.RoundOpponent.AllColumns,
	).MODELS(roundOpponents)

	_, err = insertOpponents.Exec(tx)
	if err != nil {
		return Round{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Round{}, err
	}

	return round, nil
}

func (s Service) DeleteRound(roundId int32) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete opponents
	deleteOpponents := table.RoundOpponent.DELETE().
		WHERE(table.RoundOpponent.RoundID.EQ(sqlite.Int32(roundId)))
	_, err = deleteOpponents.Exec(tx)
	if err != nil {
		return err
	}

	// Delete round
	deleteRound := table.Round.DELETE().
		WHERE(table.Round.ID.EQ(sqlite.Int32(roundId)))
	_, err = deleteRound.Exec(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
