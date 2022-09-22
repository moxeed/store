package proxy

import (
	"fmt"
	"github.com/moxeed/store/catalog/catalog_model"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/ordering/ordering_model"
	"github.com/sirupsen/logrus"
)

func CreateProduct(model catalog_model.CreateProductModel) (*catalog_model.ProductModel, error) {
	config := &common.Configuration.Store
	result := &catalog_model.ProductModel{}
	state := common.Post(config.BaseUrl+config.CreateProduct, model, &result)

	if state.IsOk {
		return result, nil
	}

	logrus.Error("CreateProduct", model, result)
	return result, fmt.Errorf("خطا در ثبت محصول")
}

func AddItem(model ordering_model.AddItemModel) error {
	config := &common.Configuration.Store
	result := make(map[string]interface{})
	state := common.Post(config.BaseUrl+config.AddItem, model, &result)

	if state.IsOk {
		return nil
	}

	logrus.Error("AddItem", model, result)
	return fmt.Errorf("خطا در ثبت سفارش")
}

func FlashBuy(model ordering_model.FlashBuyModel) error {
	config := &common.Configuration.Store
	result := make(map[string]interface{})
	state := common.Post(config.BaseUrl+config.FlashBuy, model, &result)

	if state.IsOk {
		return nil
	}

	logrus.Error("AddItem", model, result)
	return fmt.Errorf("خطا در ثبت سفارش")
}
