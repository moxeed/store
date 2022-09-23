package proxy

import (
	"fmt"
	"github.com/moxeed/store/catalog/catalog_model"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/controller/controller_model"
	"github.com/moxeed/store/ordering/ordering_model"
	"github.com/sirupsen/logrus"
)

func CreateProduct(model catalog_model.CreateProductModel) (*catalog_model.ProductModel, error) {
	config := &common.Configuration.Store
	result := &catalog_model.ProductModel{}
	state := common.Post(config.BaseUrl+config.CreateProduct, model, result)

	if state.IsOk {
		return result, nil
	}

	logrus.Error("CreateProduct", model, result)
	return result, fmt.Errorf("خطا در ثبت محصول")
}

func AddItem(model ordering_model.AddItemModel) (*ordering_model.OrderModel, error) {
	config := &common.Configuration.Store
	result := &ordering_model.OrderModel{}
	state := common.Post(config.BaseUrl+config.AddItem, model, result)

	if state.IsOk {
		return result, nil
	}

	logrus.Error("AddItem", model, result)
	return result, fmt.Errorf("خطا در ثبت سفارش")
}

func FlashBuy(model ordering_model.FlashBuyModel) (*controller_model.PaymentModel, error) {
	config := &common.Configuration.Store
	result := &controller_model.PaymentModel{}
	state := common.Post(config.BaseUrl+config.FlashBuy, model, result)

	if state.IsOk {
		return result, nil
	}

	logrus.Error("AddItem", model, result)
	return result, fmt.Errorf("خطا در ثبت سفارش")
}

func StartPayment(model controller_model.GetOrderModel) (*controller_model.PaymentModel, error) {
	config := &common.Configuration.Store
	result := &controller_model.PaymentModel{}
	state := common.Post(config.BaseUrl+config.FlashBuy, model, result)

	if state.IsOk {
		return result, nil
	}

	logrus.Error("StartPayment", model, result)
	return result, fmt.Errorf("خطا در شروع پرداخت")
}
