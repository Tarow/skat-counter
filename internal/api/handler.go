package skat

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/tarow/skat-counter/internal/skat"
	"github.com/tarow/skat-counter/internal/skat/gen/model"
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
	games, err := h.service.List()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	index := template.GameOverview(games)
	return render(c, http.StatusOK, index)
}

func (h Handler) GetGameDetails(c echo.Context) error {
	gameId := c.Param("id")

	parsedId, err := strconv.Atoi(gameId)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	game, err := h.service.Find(int32(parsedId))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	gameDetails := template.GameDetails(*game)

	return render(c, http.StatusOK, gameDetails)
}

func (h Handler) CreateGame(c echo.Context) error {
	form := struct {
		Players []string `form:"player"`
		Online  bool     `form:"online"`
		Stake   float32  `form:"stake"`
	}{}

	err := c.Bind(&form)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	players := []model.Player{}
	for _, p := range form.Players {
		players = append(players, model.Player{Name: p})
	}
	g := skat.Game{}
	g.Players = players

	g.Online = form.Online
	g.Stake = form.Stake
	g.CreatedAt = time.Now()

	g, err = h.service.Create(g)
	if err != nil {
		return err
	}

	details := template.GameDetails(g)

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/games/%v", g.ID))

	return render(c, http.StatusCreated, details)
}

func (h Handler) GetEditGameForm(c echo.Context) error {
	game, err := h.findGame(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	form := components.EditGameForm(*game)
	return render(c, http.StatusOK, form)
}

func (h Handler) EditGame(c echo.Context) error {
	game, err := h.findGame(c.Param("id"))
	if err != nil {
		return err
	}

	form := struct {
		Players []string `form:"player"`
		Online  bool     `form:"online"`
		Stake   float32  `form:"stake"`
	}{}
	c.Bind(&form)

	game.Online = form.Online
	game.Stake = form.Stake

	slices.SortFunc(game.Players, func(i, j model.Player) int {
		return slices.Index(form.Players, i.Name) - slices.Index(form.Players, j.Name)
	})

	updatedGame, err := h.service.Edit(*game)
	if err != nil {
		return err
	}

	details := template.GameDetails(updatedGame)

	return render(c, http.StatusOK, details)
}

func (h Handler) GetCreateGameForm(c echo.Context) error {
	form := components.CreateGameForm()
	return render(c, http.StatusOK, form)
}

func (h Handler) DeleteGame(c echo.Context) error {
	parsedId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.service.Delete(int32(parsedId))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	games, err := h.service.List()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	index := template.GameOverview(games)
	c.Response().Header().Set("HX-Push-Url", "/")

	return render(c, http.StatusOK, index)
}

func (h Handler) AddRound(c echo.Context) error {
	game, err := h.findGame(c.Param("id"))
	if err != nil {
		return err
	}

	var params map[string]string = make(map[string]string)
	c.Bind(&params)

	round := skat.Round{}
	round.GameID = game.ID
	round.CreatedAt = time.Now()
	for _, player := range game.Players {
		role, exists := params[player.Name]
		if !exists {
			continue
		}

		switch role {
		case "declarer":
			round.Declarer = player.ID
		case "opponent":
			round.Opponents = append(round.Opponents, player.ID)
		case "dealer":
			round.Dealer = &player.ID
		}
	}

	wonStr, exists := params["won"]
	if exists {
		won, err := strconv.ParseBool(wonStr)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		round.Won = won
	}

	gameValueStr, exists := params["gamevalue"]
	if exists {
		gameValue, err := strconv.Atoi(gameValueStr)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		round.Value = int32(gameValue)
	}

	round, err = h.service.AddRound(game.ID, round)
	if err != nil {
		return err
	}

	game.Rounds = append(game.Rounds, round)
	gameDetails := template.GameDetails(*game)
	return render(c, http.StatusCreated, gameDetails)
}

func render(ctx echo.Context, status int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(status)

	err := t.Render(context.Background(), ctx.Response().Writer)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, "failed to render response template")
	}

	return nil
}

func (h Handler) findGame(gameId string) (*skat.Game, error) {
	parsedId, err := strconv.Atoi(gameId)
	if err != nil {
		return nil, err
	}

	return h.service.Find(int32(parsedId))
}
