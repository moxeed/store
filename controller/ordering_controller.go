package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/ordering"
)

func GetOrder(c echo.Context) (err error) {
	model := ordering.GetOrderModel{}
	err = c.Bind(&model)

	result := ordering.GetOrder(model)
	err = c.JSON(200, result)
	return
}

func AddItem(c echo.Context) (err error) {
	model := ordering.AddItemModel{}
	err = c.Bind(&model)

	result := ordering.AddItem(model)
	err = c.JSON(200, result)
	return
}
