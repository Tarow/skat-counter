package main

import (
	"database/sql"
	"embed"
	"log"

	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	api "github.com/tarow/skat-counter/internal/api"
	"github.com/tarow/skat-counter/internal/skat"
)

//go:embed static/*
var staticAssets embed.FS

func main() {
	db, err := sql.Open("sqlite", "skat.sqlite")
	createTables(db)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()

	handler := api.NewHandler(skat.NewService(db))
	registerRoutes(e, handler)

	e.Logger.Fatal(e.Start(":8080"))
}

type Form struct {
	Playernames []string `form:"playername"`
}

func registerRoutes(e *echo.Echo, handler api.Handler) {
	e.StaticFS("/static", echo.MustSubFS(staticAssets, "static"))
	e.GET("/", handler.GetIndex)

	e.POST("/games", handler.CreateGame)
	e.GET("/games/create", handler.GetCreateGameForm)

	e.GET("/games/:id", handler.GetGameDetails)

	e.GET("/games/:id/edit", handler.GetEditGameForm)
	e.PUT("/games/:id/edit", handler.EditGame)

	e.DELETE("/games/:id", handler.DeleteGame)
	e.POST("/games/:id/rounds", handler.AddRound)

}
func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS player (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name TEXT UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS game (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			stake REAL NOT NULL,
			online BOOLEAN NOT NULL
		);

		CREATE TABLE IF NOT EXISTS game_player (
			game_id INTEGER NOT NULL,
			player_id INTEGER NOT NULL,
			rank INTEGER NOT NULL,
			FOREIGN KEY(game_id) REFERENCES games(id),
			FOREIGN KEY(player_id) REFERENCES players(id)
		);

		CREATE TABLE IF NOT EXISTS round (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			game_id INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			dealer INTEGER,
			declarer INTEGER NOT NULL,
			won BOOLEAN NOT NULL,
			value INTEGER NOT NULL,
			FOREIGN KEY(game_id) REFERENCES games(id),
			FOREIGN KEY(dealer) REFERENCES players(id),
			FOREIGN KEY(declarer) REFERENCES players(id)
		);

		CREATE TABLE IF NOT EXISTS round_opponent (
			round_id INTEGER NOT NULL,
			player_id INTEGER NOT NULL,
			PRIMARY KEY (round_id, player_id),
			FOREIGN KEY(round_id) REFERENCES rounds(id),
			FOREIGN KEY(player_id) REFERENCES players(id)
		);
	`)

	return err
}
