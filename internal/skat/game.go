package skat

import (
	"slices"

	"github.com/go-jet/jet/v2/sqlite"
	"github.com/tarow/skat-counter/internal/skat/gen/model"
	"github.com/tarow/skat-counter/internal/skat/gen/table"
)

func (s Service) List() ([]Game, error) {
	stmt := sqlite.SELECT(
		table.Game.AllColumns,
		table.Player.AllColumns,
		table.Round.AllColumns,
		table.RoundOpponent.PlayerID.AS("opponent.id"),
	).FROM(table.Game.
		INNER_JOIN(table.GamePlayer, table.Game.ID.EQ(table.GamePlayer.GameID)).
		INNER_JOIN(table.Player, table.Player.ID.EQ(table.GamePlayer.PlayerID)).
		LEFT_JOIN(table.Round, table.Round.GameID.EQ(table.Game.ID)).
		LEFT_JOIN(table.RoundOpponent, table.RoundOpponent.RoundID.EQ(table.Round.ID).AND(
			table.RoundOpponent.PlayerID.EQ(table.Player.ID),
		)),
	).ORDER_BY(table.Game.CreatedAt.DESC(), table.GamePlayer.Rank.ASC(), table.Round.CreatedAt.ASC())

	games := make([]Game, 0)
	err := stmt.Query(s.db, &games)
	if err != nil {
		return []Game{}, err
	}

	return games, nil
}

func (s Service) Find(gameId int32) (*Game, error) {
	stmt := sqlite.SELECT(
		table.Game.AllColumns,
		table.Player.AllColumns,
		table.Round.AllColumns,
		table.RoundOpponent.PlayerID.AS("opponent.id"),
	).FROM(table.Game.
		INNER_JOIN(table.GamePlayer, table.Game.ID.EQ(table.GamePlayer.GameID)).
		INNER_JOIN(table.Player, table.Player.ID.EQ(table.GamePlayer.PlayerID)).
		LEFT_JOIN(table.Round, table.Round.GameID.EQ(table.Game.ID)).
		LEFT_JOIN(table.RoundOpponent, table.RoundOpponent.RoundID.EQ(table.Round.ID).AND(
			table.RoundOpponent.PlayerID.EQ(table.Player.ID),
		)),
	).WHERE(table.Game.ID.EQ(sqlite.Int32(gameId))).
		ORDER_BY(table.GamePlayer.Rank.ASC(), table.Round.CreatedAt.ASC())

	game := Game{}
	err := stmt.Query(s.db, &game)
	if err != nil {
		return &Game{}, err
	}

	return &game, nil
}

func (s Service) Create(g Game) (Game, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Game{}, err
	}

	defer tx.Rollback()

	// Create players if they dont exist
	stmt := table.Player.
		INSERT(table.Player.Name).
		MODELS(g.Players).
		ON_CONFLICT().DO_NOTHING().RETURNING(table.Player.AllColumns)

	_, err = stmt.Exec(tx)
	if err != nil {

		return Game{}, err
	}

	// Insert game
	stmt = table.Game.INSERT(
		table.Game.Online,
		table.Game.Stake,
		table.Game.CreatedAt,
	).MODEL(g)

	res, err := stmt.Exec(tx)
	if err != nil {
		return Game{}, err
	}

	gameId, err := res.LastInsertId()
	if err != nil {
		return Game{}, err
	}

	playerNames := []sqlite.Expression{}
	playerNamesStr := []string{}
	for _, p := range g.Players {
		playerNames = append(playerNames, sqlite.String(p.Name))
		playerNamesStr = append(playerNamesStr, p.Name)
	}

	// Link players to game
	selectPlayers := table.Player.SELECT(table.Player.AllColumns).
		FROM(table.Player).
		WHERE(table.Player.Name.IN(playerNames...))

	players := make([]model.Player, 0)
	err = selectPlayers.Query(tx, &players)
	if err != nil {
		return Game{}, err
	}

	gamePlayers := []model.GamePlayer{}
	for _, p := range players {
		gamePlayers = append(gamePlayers, model.GamePlayer{
			GameID:   int32(gameId),
			PlayerID: p.ID,
			Rank:     int32(slices.Index(playerNamesStr, p.Name)),
		})
	}
	stmt = table.GamePlayer.INSERT(table.GamePlayer.AllColumns).
		MODELS(gamePlayers)
	_, err = stmt.Exec(tx)
	if err != nil {
		return Game{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Game{}, err
	}

	g.ID = int32(gameId)
	return g, nil
}

func (s Service) UpdateGame(g Game) (Game, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Game{}, err
	}
	defer tx.Rollback()

	stmt := table.Game.UPDATE(table.Game.Stake, table.Game.Online).
		MODEL(g).
		WHERE(table.Game.ID.EQ(sqlite.Int32(g.ID)))

	_, err = stmt.Exec(tx)
	if err != nil {
		return Game{}, err
	}

	for rank, player := range g.Players {
		stmt := table.GamePlayer.UPDATE(table.GamePlayer.Rank).
			MODEL(model.GamePlayer{
				GameID:   g.ID,
				PlayerID: player.ID,
				Rank:     int32(rank),
			}).
			WHERE(table.GamePlayer.GameID.EQ(sqlite.Int32(g.ID)).
				AND(table.GamePlayer.PlayerID.EQ(sqlite.Int32(player.ID))))
		_, err = stmt.Exec(tx)
		if err != nil {
			return Game{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return Game{}, err
	}

	return g, nil
}

func (s Service) DeleteGame(gameId int32) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete round opponents
	stmt := table.RoundOpponent.DELETE().WHERE(
		table.RoundOpponent.RoundID.IN(
			table.Round.SELECT(table.Round.ID).FROM(table.Round).WHERE(
				table.Round.GameID.EQ(sqlite.Int32(gameId)),
			),
		))

	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}

	//Delete rounds
	stmt = table.Round.DELETE().WHERE(table.Round.GameID.EQ(sqlite.Int32(gameId)))
	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}

	//Delete game players references
	stmt = table.GamePlayer.DELETE().WHERE(table.GamePlayer.GameID.EQ(sqlite.Int32(gameId)))
	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}

	//Delete game
	stmt = table.Game.DELETE().WHERE(table.Game.ID.EQ(sqlite.Int32(gameId)))
	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}

	// Delete orphan players that arent refenced in any game
	stmt = table.Player.DELETE().WHERE(table.Player.ID.NOT_IN(
		table.GamePlayer.SELECT(table.GamePlayer.PlayerID),
	))
	_, err = stmt.Exec(tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
