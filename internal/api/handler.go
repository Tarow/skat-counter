package skat

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/tarow/skat-counter/internal/skat"
	template "github.com/tarow/skat-counter/templates"
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

func render(ctx echo.Context, status int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(status)

	err := t.Render(context.Background(), ctx.Response().Writer)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to render response template")
	}

	return nil
}
