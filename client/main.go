package main

import (
	"github.com/labstack/echo/v4"
)

type ErrorModel struct {
	ReferenceCode uint
	Error         string
}

func main() {
	router := echo.New()

	router.POST("ok", func(c echo.Context) error {
		return c.JSON(200, "ok")
	})

	router.POST("nok", func(c echo.Context) error {
		result := []ErrorModel{{
			ReferenceCode: 79201,
			Error:         "error",
		}}
		return c.JSON(400, result)
	})

	router.POST("ambiguous", func(c echo.Context) error {
		return c.JSON(500, "error")
	})

	err := router.Start("localhost:8070")
	if err != nil {
		panic(err)
	}
}
