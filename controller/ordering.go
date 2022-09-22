package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/ordering"
	"github.com/moxeed/store/payment"
)

type GetOrderModel struct {
	ID uint `query:"id"`
}

type GetBasketModel struct {
	CustomerCode uint `query:"customerCode"`
}

type LockBasketModel struct {
	CustomerCode uint `query:"customerCode"`
}

type PaymentModel struct {
	Order      ordering.OrderModel
	PaymentUrl string
}

func GetOrder(c echo.Context) (err error) {
	model := GetOrderModel{}
	err = c.Bind(&model)

	result, err := ordering.GetOrder(model.ID)
	err = common.WriteIfNoError(&c, err, result)
	return
}

func GetBasket(c echo.Context) (err error) {
	model := GetBasketModel{}
	err = c.Bind(&model)

	result, err := ordering.GetBasket(model.CustomerCode)
	err = common.WriteIfNoError(&c, err, result)
	return
}

func AddItem(c echo.Context) (err error) {
	model := ordering.AddItemModel{}
	err = c.Bind(&model)

	result, err := ordering.AddItem(model)
	err = common.WriteIfNoError(&c, err, result)
	return
}

func LockForPayment(c echo.Context) (err error) {
	model := LockBasketModel{}
	err = c.Bind(&model)

	result, isFree, err := ordering.LockForPayment(model.CustomerCode)

	var paymentResult string
	if err == nil && !isFree {
		paymentResult, err = payment.OpenTerminal(result.ID)
	}

	err = common.WriteIfNoError(&c, err, PaymentModel{
		Order:      result,
		PaymentUrl: paymentResult,
	})
	return
}
