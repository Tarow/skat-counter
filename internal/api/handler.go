package skat

import (
	"context"
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

	game, err := h.service.Find(parsedId)
	if err != nil {
		return err
	}

	gameDetails := template.GameDetails(game)

	return render(c, http.StatusOK, gameDetails)
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
