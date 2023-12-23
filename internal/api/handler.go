package skat

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/tarow/skat-counter/internal/skat"
	template "github.com/tarow/skat-counter/templates"
	"github.com/tarow/skat-counter/templates/components"
)

type Handler struct {
	service skat.Service
}

func NewHandler(service skat.Service) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) GetIndex(c echo.Context) error {
	index := template.GameOverview(h.service.List())
	return render(c, http.StatusOK, index)
}

func (h Handler) GetGameDetails(c echo.Context) error {
	gameId := c.Param("id")

	parsedId, err := strconv.Atoi(gameId)
	if err != nil {
		return err
	}

	game := h.service.Find(parsedId)
	if err != nil {
		return err
	}

	gameDetails := template.GameDetails(*game)

	return render(c, http.StatusOK, gameDetails)
}

func (h Handler) CreateGame(c echo.Context) error {
	g := skat.Game{}
	c.Bind(&g)
	fmt.Println("received game", g)
	g = h.service.Create(g)
	fmt.Println("created game", g)
	details := template.GameDetails(g)

	c.Response().Header().Set("Access-Control-Expose-Headers", "*")
	//c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	//c.Response().Header().Set("Access-Control-Allow-Headers", "*")
	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/games/%v", g.Id))

	return render(c, http.StatusCreated, details)
}

func (h Handler) AddRound(c echo.Context) error {
	gameId := c.Param("id")
	parsedId, err := strconv.Atoi(gameId)
	if err != nil {
		return err
	}

	game := h.service.Find(parsedId)
	if err != nil {
		return err
	}
	fmt.Printf("rounds before: %+v", len(game.Rounds))
	var params map[string]string = make(map[string]string)
	c.Bind(&params)

	round := skat.Round{}
	for _, player := range game.Players {
		role, exists := params[player]
		if !exists {
			continue
		}

		switch role {
		case "declarer":
			round.Declarer = player
		case "opponent":
			round.Opponents = append(round.Opponents, player)
		case "dealer":
			round.Dealer = player
		}
	}

	wonStr, exists := params["won"]
	if exists {
		won, err := strconv.ParseBool(wonStr)
		if err != nil {
			return err
		}
		round.Won = won
	}

	gameValueStr, exists := params["gamevalue"]
	if exists {
		gameValue, err := strconv.Atoi(gameValueStr)
		if err != nil {
			return err
		}
		round.Value = gameValue
	}

	game.Rounds = append(game.Rounds, round)
	fmt.Printf("rounds after: %+v", len(game.Rounds))
	gameDetails := template.GameDetails(*game)
	return render(c, http.StatusCreated, gameDetails)
}

func (h Handler) GetCreateGameForm(c echo.Context) error {
	form := components.CreateGameForm()
	return render(c, http.StatusOK, form)
}

func render(ctx echo.Context, status int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(status)

	err := t.Render(context.Background(), ctx.Response().Writer)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to render response template")
	}

	return nil
}
