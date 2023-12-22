package main

import (
	"fmt"

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
	e.GET("/", handler.GetIndex)

	e.Any("/test", func(c echo.Context) error {
		//body, _ := io.ReadAll(c.Request().Body)
		//fmt.Printf("received request on test endpoint: %v\n", string(body))

		f := Form{}
		c.Bind(&f)
		fmt.Printf("Received form: %+v\n", f)
		return nil
	})
}
