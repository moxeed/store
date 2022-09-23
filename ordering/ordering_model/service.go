package ordering_model

import "time"

type OrderIdentifier struct {
	ID            *uint   `json:"id,omitempty"`
	ReferenceCode *string `json:"referenceCode,omitempty"`
}

type StateModel struct {
	State      int       `json:"state,omitempty"`
	StateTitle string    `json:"stateTitle,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type OrderModel struct {
	ID             uint             `json:"id,omitempty"`
	UserCode       uint             `json:"userCode,omitempty"`
	CustomerCode   uint             `json:"customerCode,omitempty"`
	LastState      int              `json:"lastState,omitempty"`
	LastStateTitle string           `json:"lastStateTitle,omitempty"`
	TotalAmount    uint             `json:"totalAmount,omitempty"`
	Items          []OrderItemModel `json:"items,omitempty"`
	States         []StateModel     `json:"states,omitempty"`
	CreatedAt      time.Time        `json:"createdAt"`
}

type OrderHeaderModel struct {
	ID             uint         `json:"id,omitempty"`
	UserCode       uint         `json:"userCode,omitempty"`
	CustomerCode   uint         `json:"customerCode,omitempty"`
	LastState      int          `json:"lastState,omitempty"`
	LastStateTitle string       `json:"lastStateTitle,omitempty"`
	TotalAmount    uint         `json:"totalAmount,omitempty"`
	States         []StateModel `json:"states,omitempty"`
	CreatedAt      time.Time    `json:"createdAt"`
}

type OrderItemModel struct {
	ID                   uint         `json:"id,omitempty"`
	ProductTitle         string       `json:"productTitle,omitempty"`
	Category             string       `json:"category,omitempty"`
	ProductCode          uint         `json:"productCode,omitempty"`
	ReferenceProductCode uint         `json:"referenceProductCode,omitempty"`
	Price                uint         `json:"price,omitempty"`
	Quantity             uint         `json:"quantity,omitempty"`
	LastState            int          `json:"lastState,omitempty"`
	LastStateTitle       string       `json:"lastStateTitle,omitempty"`
	Errors               []string     `json:"errors,omitempty"`
	States               []StateModel `json:"states,omitempty"`
}

type AddItemModel struct {
	UserCode      uint `json:"userCode,omitempty"`
	CustomerCode  uint `json:"customerCode,omitempty"`
	ProductCode   uint `json:"productCode,omitempty"`
	ReferenceCode uint `json:"referenceCode,omitempty"`
	Quantity      uint `json:"quantity,omitempty"`
}

type FlashBuyModel struct {
	UserCode           uint   `json:"userCode,omitempty"`
	CustomerCode       uint   `json:"customerCode,omitempty"`
	ProductCode        uint   `json:"productCode,omitempty"`
	ReferenceCode      uint   `json:"referenceCode,omitempty"`
	Quantity           uint   `json:"quantity,omitempty"`
	OrderReferenceCode string `json:"orderReferenceCode,omitempty"`
}
