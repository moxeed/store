package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/controller/controller_model"
	"github.com/moxeed/store/ordering"
	"github.com/moxeed/store/ordering/ordering_model"
	"github.com/moxeed/store/payment"
)

func GetOrder(c echo.Context) (err error) {
	model := controller_model.GetOrderModel{}
	err = c.Bind(&model)

	result := ordering_model.OrderModel{}

	result, err = ordering.GetOrder(ordering_model.OrderIdentifier{
		ID:            model.ID,
		ReferenceCode: model.ReferenceCode,
	})

	err = common.WriteIfNoError(&c, err, result)
	return
}

func GetList(c echo.Context) (err error) {
	model := controller_model.GetOrderListModel{}
	err = c.Bind(&model)

	result, totalCount := ordering.GetOrderList(model.CustomerCode, model.Offset)

	err = common.WriteIfNoError(&c, nil, struct {
		Rows       []*ordering_model.OrderHeaderModel
		TotalCount int64
	}{
		Rows:       result,
		TotalCount: totalCount,
	})

	return
}

func GetBasket(c echo.Context) (err error) {
	model := controller_model.GetBasketModel{}
	err = c.Bind(&model)

	result, err := ordering.GetBasket(model.CustomerCode)
	err = common.WriteIfNoError(&c, err, result)
	return
}

func AddItem(c echo.Context) (err error) {
	model := ordering_model.AddItemModel{}
	err = c.Bind(&model)

	result, err := ordering.AddItem(model)
	err = common.WriteIfNoError(&c, err, result)
	return
}

func getPayment(order *ordering_model.OrderModel, isFree bool, prevError error) (paymentModel controller_model.PaymentModel, err error) {
	paymentModel.Order = *order

	if !isFree && order.LastState == ordering.PaymentPending && prevError == nil {
		var paymentResult string
		paymentResult, err = payment.OpenTerminal(order.ID)

		paymentModel.PaymentUrl = paymentResult
	}

	return
}

func StartPayment(c echo.Context) (err error) {
	model := controller_model.GetOrderModel{}
	err = c.Bind(&model)

	result, isFree, err := ordering.StartPayment(ordering_model.OrderIdentifier{
		ID:            model.ID,
		ReferenceCode: model.ReferenceCode,
	})

	factor, err := getPayment(&result, isFree, err)
	err = common.WriteIfNoError(&c, err, factor)
	return
}

func FlashBuy(c echo.Context) (err error) {
	model := ordering_model.FlashBuyModel{}
	err = c.Bind(&model)

	result, isFree, err := ordering.FlashBuy(model)

	factor, err := getPayment(&result, isFree, err)
	err = common.WriteIfNoError(&c, err, factor)
	return
}
