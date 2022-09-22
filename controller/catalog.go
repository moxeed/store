package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/catalog"
	"github.com/moxeed/store/catalog/catalog_model"
)

func AddProduct(c echo.Context) (err error) {

	model := catalog_model.CreateProductModel{}
	err = c.Bind(&model)
	result := catalog.CreateProduct(model)
	err = c.JSON(200, result)

	return
}
