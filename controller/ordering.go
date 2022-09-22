package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/ordering"
	"github.com/moxeed/store/payment"
)

type GetOrderModel struct {
	ID            *uint   `query:"id"`
	ReferenceCode *string `query:"referenceCode"`
}

type GetOrderListModel struct {
	CustomerCode uint `query:"customerCode"`
	Offset       int  `query:"offset"`
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

	result := ordering.OrderModel{}

	result, err = ordering.GetOrder(ordering.OrderIdentifier{
		ID:            model.ID,
		ReferenceCode: model.ReferenceCode,
	})

	err = common.WriteIfNoError(&c, err, result)
	return
}

func GetList(c echo.Context) (err error) {
	model := GetOrderListModel{}
	err = c.Bind(&model)

	result, totalCount := ordering.GetOrderList(model.CustomerCode, model.Offset)

	err = common.WriteIfNoError(&c, nil, struct {
		Rows       []*ordering.OrderHeaderModel
		TotalCount int64
	}{
		Rows:       result,
		TotalCount: totalCount,
	})

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

func getPayment(order *ordering.OrderModel, isFree bool) (paymentModel PaymentModel, err error) {
	paymentModel.Order = *order

	if !isFree {
		var paymentResult string
		paymentResult, err = payment.OpenTerminal(order.ID)

		paymentModel.PaymentUrl = paymentResult
	}

	return
}

func StartPayment(c echo.Context) (err error) {
	model := GetOrderModel{}
	err = c.Bind(&model)

	result, isFree, err := ordering.StartPayment(ordering.OrderIdentifier{
		ID:            model.ID,
		ReferenceCode: model.ReferenceCode,
	})

	factor, err := getPayment(&result, isFree)
	err = common.WriteIfNoError(&c, err, factor)
	return
}

func FlashBuy(c echo.Context) (err error) {
	model := ordering.FlashBuyModel{}
	err = c.Bind(&model)

	result, isFree, err := ordering.FlashBuy(model)

	factor, err := getPayment(&result, isFree)
	err = common.WriteIfNoError(&c, err, factor)
	return
}
