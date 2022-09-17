package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/catalog"
)

func AddProduct(c echo.Context) (err error) {

	model := catalog.CreateProductModel{}
	err = c.Bind(&model)
	result := catalog.CreateProduct(model)
	err = c.JSON(200, result)

	return
}
