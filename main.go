package main

import (
	"github.com/labstack/echo/v4"
	api "github.com/tarow/skat-counter/internal/api"
	"github.com/tarow/skat-counter/internal/skat"
)

func main() {
	e := echo.New()

	handler := api.NewHandler(skat.NewService())
	registerRoutes(e, handler)

	e.Logger.Fatal(e.Start(":8080"))
}

type Form struct {
	Playernames []string `form:"playername"`
}

func registerRoutes(e *echo.Echo, handler api.Handler) {
	e.Static("/static", "./static")
	e.GET("/", handler.GetIndex)

	e.POST("/games", handler.CreateGame)
	e.GET("/games/:id", handler.GetGameDetails)
	e.POST("/games/:id/rounds", handler.AddRound)

}
