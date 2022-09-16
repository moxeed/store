package ordering

import (
	"github.com/moxeed/store/common"
	"gorm.io/gorm/clause"
	"time"
)

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
	Items          []OrderItemModel
	States         []StateModel
	CreatedAt      time.Time
}

type OrderItemModel struct {
	ID             uint
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

type GetOrderModel struct {
	ID           *uint `query:"id"`
	UserCode     uint  `query:"userCode"`
	CustomerCode uint  `query:"customerCode"`
}

func GetOrder(model GetOrderModel) OrderModel {
	order := NewOrder(model.UserCode, model.CustomerCode)

	if model.ID != nil {
		common.DB.Preload(clause.Associations).Find(&order, model.ID)
	} else {
		common.DB.Preload(clause.Associations).Where(&order).Find(&order)
	}
	return order.ToModel()
}

func AddItem(model AddItemModel) OrderModel {
	order := NewOrder(model.UserCode, model.CustomerCode)
	common.DB.Preload(clause.Associations).Where(&order).Find(&order)

	order.addItem(model.ProductCode, model.ReferenceCode, model.Quantity)
	common.DB.Save(&order)
	return order.ToModel()
}

func stateText(state int) string {
	switch state {
	case Basket:
		return "سبد خرید"
	case PaymentPending:
		return "منتظر پرداخت"
	case ProcessFailed:
		return "خطای پردازش"
	case Processed:
		return "پردازش شده"
	case Canceled:
		return "لفو"
	case Paid:
		return "پرداخت شده"
	default:
		return ""
	}
}

func (os *OrderState) toModel() StateModel {
	return StateModel{
		State:      os.State,
		StateTitle: stateText(os.State),
		CreatedAt:  os.CreatedAt,
	}
}

func (ois *OrderItemState) toModel() StateModel {
	return StateModel{
		State:      ois.State,
		StateTitle: stateText(ois.State),
		CreatedAt:  ois.CreatedAt,
	}
}

func (oi *OrderItem) toModel() OrderItemModel {
	var errors []string
	var states []StateModel

	for _, e := range oi.Errors {
		errors = append(errors, e.Error)
	}

	for _, s := range oi.States {
		states = append(states, s.toModel())
	}

	return OrderItemModel{
		ID:             oi.ID,
		ProductCode:    oi.ProductCode,
		Price:          oi.Price,
		Quantity:       oi.Quantity,
		LastState:      oi.LastState,
		LastStateTitle: stateText(oi.LastState),
		Errors:         errors,
		States:         states,
	}
}

func (o *Order) ToModel() OrderModel {
	var items []OrderItemModel
	var states []StateModel

	for _, item := range o.Items {
		items = append(items, item.toModel())
	}

	for _, state := range o.States {
		states = append(states, state.toModel())
	}

	return OrderModel{
		ID:             o.ID,
		UserCode:       o.UserCode,
		CustomerCode:   o.CustomerCode,
		LastState:      o.LastState,
		LastStateTitle: stateText(o.LastState),
		Items:          items,
		States:         states,
		CreatedAt:      o.CreatedAt,
	}
}
