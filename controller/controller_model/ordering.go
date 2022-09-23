package controller_model

import "github.com/moxeed/store/ordering/ordering_model"

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
	Order      ordering_model.OrderModel `json:"order"`
	PaymentUrl string                    `json:"paymentUrl,omitempty"`
}
