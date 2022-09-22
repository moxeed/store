package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/controller/controller_model"
	"github.com/moxeed/store/payment"
)

func OpenTerminal(c echo.Context) (err error) {
	model := controller_model.OpenTerminalModel{}
	err = c.Bind(&model)

	result, err := payment.OpenTerminal(model.OrderCode)
	return c.JSON(200, result)
}

func Verify(c echo.Context) (err error) {
	model := controller_model.VerifyModel{}
	err = c.Bind(&model)

	result := payment.Verify(model.Authority)
	return c.JSON(200, result)
}
