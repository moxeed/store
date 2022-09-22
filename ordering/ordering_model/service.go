package ordering_model

import "time"

type OrderIdentifier struct {
	ID            *uint
	ReferenceCode *string
}

type StateModel struct {
	State      int
	StateTitle string
	CreatedAt  time.Time
}

type OrderModel struct {
	ID             uint
	UserCode       uint
	CustomerCode   uint
	LastState      int
	LastStateTitle string
	TotalAmount    uint
	Items          []OrderItemModel
	States         []StateModel
	CreatedAt      time.Time
}

type OrderHeaderModel struct {
	ID             uint
	UserCode       uint
	CustomerCode   uint
	LastState      int
	LastStateTitle string
	TotalAmount    uint
	States         []StateModel
	CreatedAt      time.Time
}

type OrderItemModel struct {
	ID             uint
	ProductTitle   string
	Category       string
	ProductCode    uint
	Price          uint
	Quantity       uint
	LastState      int
	LastStateTitle string
	Errors         []string
	States         []StateModel
}

type AddItemModel struct {
	UserCode      uint
	CustomerCode  uint
	ProductCode   uint
	ReferenceCode uint
	Quantity      uint
}

type FlashBuyModel struct {
	UserCode           uint
	CustomerCode       uint
	ProductCode        uint
	ReferenceCode      uint
	Quantity           uint
	OrderReferenceCode string
}
