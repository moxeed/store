package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/controller"
)

func main() {
	router := echo.New()

	product := router.Group("/product")
	product.POST("", controller.AddProduct)

	order := router.Group("/order")
	order.GET("", controller.GetOrder)
	order.POST("", controller.AddItem)

	err := router.Start(fmt.Sprintf("localhost:%d", common.Configuration.ListenPort))

	if err != nil {
		panic(err)
	}
}
