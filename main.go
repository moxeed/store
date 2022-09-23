package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/controller"
	"github.com/nullseed/logruseq"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.AddHook(logruseq.NewSeqHook(common.Configuration.Seq.Url,
		logruseq.OptionAPIKey(common.Configuration.Seq.ApiKey)))
	router := echo.New()

	product := router.Group("/product")
	product.POST("", controller.AddProduct)

	order := router.Group("/order")
	order.GET("", controller.GetOrder)
	order.GET("/list", controller.GetList)
	order.GET("/basket", controller.GetBasket)
	order.POST("", controller.AddItem)
	order.POST("/flash", controller.FlashBuy)
	order.POST("/pay", controller.StartPayment)

	payment := router.Group("/payment")
	payment.GET("/terminal", controller.OpenTerminal)
	payment.GET("/verify*", controller.Verify)

	err := router.Start(fmt.Sprintf("localhost:%d", common.Configuration.ListenPort))

	if err != nil {
		panic(err)
	}
}
