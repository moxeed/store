package ordering

import (
	"github.com/moxeed/store/catalog"
	"github.com/moxeed/store/common"
)

type ErrorModel struct {
	ReferenceCode uint
	Error         string
}

func getProductCodes(items []*OrderItem) []uint {
	result := make([]uint, 0)
	for _, item := range items {
		result = append(result, item.ProductCode)
	}

	return result
}

func (o *Order) getProductCodes() []uint {
	return getProductCodes(o.Items)
}

func (o *Order) validate() bool {
	productCodes := o.getProductCodes()
	callBacks := catalog.GetProductCallBacks(productCodes)

	calls := make(map[string][]uint)
	items := make(map[uint]*OrderItem)
	for _, item := range o.Items {
		callBack := callBacks[item.ProductCode].ValidationCallBack
		calls[callBack] = append(calls[callBack], item.ReferenceCode)
		items[item.ReferenceCode] = item
	}

	errors := make([]ErrorModel, 0)

	for callBack, referenceCodes := range calls {
		result := make([]ErrorModel, 0)
		state := common.Post(callBack, referenceCodes, &result)
		if !state.IsOk {
			errors = append(errors, result...)
		}

		for _, code := range referenceCodes {
			items[code].addErrors(&result)
		}
	}

	return len(errors) == 0
}

func dispatchCheckout(orderItems []*OrderItem) {
	if len(orderItems) == 0 {
		return
	}

	productCodes := getProductCodes(orderItems)
	callBacks := catalog.GetProductCallBacks(productCodes)

	calls := make(map[string][]uint)
	items := make(map[uint]*OrderItem)

	for _, item := range orderItems {
		callBack := callBacks[item.ProductCode].CheckOutCallBack
		calls[callBack] = append(calls[callBack], item.ReferenceCode)
		items[item.ReferenceCode] = item
	}

	for callBack, referenceCodes := range calls {
		result := make([]ErrorModel, 0)
		state := common.Post(callBack, referenceCodes, &result)

		for _, code := range referenceCodes {
			if state.IsAmbiguous {
				items[code].failedToProcess()
			} else {
				items[code].checkOut(&result)
			}
		}
	}
}

func (o *Order) dispatchCheckout() {
	dispatchCheckout(o.Items)
}
