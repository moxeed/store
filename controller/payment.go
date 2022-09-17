package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/payment"
)

type OpenTerminalModel struct {
	OrderCode uint `query:"orderCode"`
}

type VerifyModel struct {
	Authority string `query:"Authority"`
}

func OpenTerminal(c echo.Context) (err error) {
	model := OpenTerminalModel{}
	err = c.Bind(&model)

	result, err := payment.OpenTerminal(model.OrderCode)
	return c.JSON(200, result)
}

func Verify(c echo.Context) (err error) {
	model := VerifyModel{}
	err = c.Bind(&model)

	result := payment.Verify(model.Authority)
	return c.JSON(200, result)
}
