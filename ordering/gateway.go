package ordering

import (
	"github.com/moxeed/store/catalog"
	"github.com/moxeed/store/common"
)

type ErrorModel struct {
	ReferenceCode uint
	Error         string
}

func getProductCodes(items *[]OrderItem) []uint {
	var result []uint
	for _, item := range *items {
		result = append(result, item.ProductCode)
	}

	return result
}

func (o *Order) getProductCodes() []uint {
	return getProductCodes(&o.Items)
}

func (o *Order) validate() bool {
	productCodes := o.getProductCodes()
	callBacks := catalog.GetProductCallBacks(productCodes)

	calls := make(map[string][]uint)
	items := make(map[uint]*OrderItem)
	for _, item := range o.Items {
		callBack := callBacks[item.ProductCode].ValidationCallBack
		calls[callBack] = append(calls[callBack], item.ReferenceCode)
		items[item.ReferenceCode] = &item
	}

	var errors []ErrorModel

	for callBack, referenceCodes := range calls {
		var result []ErrorModel
		isSuccess := common.Post(callBack, referenceCodes, &result)
		if isSuccess {
			errors = append(errors, result...)
		}

		for _, code := range referenceCodes {
			items[code].addErrors(&result)
		}
	}

	return len(errors) == 0
}

func dispatchCheckout(orderItems *[]OrderItem) bool {
	if len(*orderItems) == 0 {
		return true
	}

	productCodes := getProductCodes(orderItems)
	callBacks := catalog.GetProductCallBacks(productCodes)

	calls := make(map[string][]uint)
	items := make(map[uint]*OrderItem)

	for _, item := range *orderItems {
		callBack := callBacks[item.ProductCode].CheckOutCallBack
		calls[callBack] = append(calls[callBack], item.ReferenceCode)
		items[item.ReferenceCode] = &item
	}

	var errors []ErrorModel

	for callBack, referenceCodes := range calls {
		var result []ErrorModel
		isOk := common.Post(callBack, referenceCodes, &result)

		if isOk {
			errors = append(errors, result...)
		}

		for _, code := range referenceCodes {
			if isOk {
				items[code].checkOut(&result)
			} else {
				items[code].failedToProcess()
			}
		}
	}

	return len(errors) == 0
}

func (o *Order) dispatchCheckout() bool {
	return dispatchCheckout(&o.Items)
}
