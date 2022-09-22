package ordering

import (
	"github.com/moxeed/store/catalog"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/payment"
	"gorm.io/gorm"
	"time"
)

const (
	Basket = iota + 1
	PaymentPending
	ProcessFailed
	Processed
	Canceled
	Paid
)

type Order struct {
	gorm.Model
	UserCode      uint
	CustomerCode  uint
	ReferenceCode *string
	LastState     int
	Items         []*OrderItem
	States        []OrderState
}

type OrderItem struct {
	gorm.Model
	ProductCode   uint
	ReferenceCode uint
	ProductTitle  string
	Category      string
	Price         uint
	Quantity      uint
	LastState     int
	OrderID       uint
	Errors        []OrderItemError
	States        []OrderItemState
}

type OrderItemError struct {
	ID          uint
	Error       string
	OrderItemID uint
	CreatedAt   time.Time
}

type OrderItemState struct {
	ID          uint
	State       int
	OrderItemID uint
	CreatedAt   time.Time
}

type OrderState struct {
	ID        uint
	State     int
	OrderID   uint
	CreatedAt time.Time
}

func NewOrder(userCode uint, customerCode uint, referenceCode *string) Order {
	order := Order{
		UserCode:      userCode,
		CustomerCode:  customerCode,
		ReferenceCode: referenceCode,
	}

	order.setState(Basket)

	return order
}

func (o *Order) addItem(productCode uint, referenceCode uint, quantity uint) (err error) {
	product, err := catalog.GetProduct(productCode)

	item := &OrderItem{
		ProductTitle:  product.Title,
		Category:      product.Category.Key,
		ProductCode:   productCode,
		ReferenceCode: referenceCode,
		Price:         product.Price,
		Quantity:      quantity,
		OrderID:       o.ID,
	}

	item.setState(Basket)

	o.Items = append(o.Items, item)

	return
}

func (o *Order) lockForPayment() (bool, bool) {
	isOk := o.validate()

	if o.TotalAmount() == 0 {
		o.setState(Paid)
		o.dispatchCheckout()
		return true, true
	}

	if isOk {

		o.setState(PaymentPending)
		for _, item := range o.Items {
			item.setState(PaymentPending)
		}

		payment.CreatePayment(o.CustomerCode, o.ID, o.TotalAmount())
	}

	return isOk, false
}

func (o *Order) checkOut() {
	inquiry := payment.IsPaid(o.ID)

	if !inquiry.DoesExists {
		payment.CreatePayment(o.CustomerCode, o.ID, o.TotalAmount())
		return
	}

	if !inquiry.IsPaid {
		return
	}

	o.setState(Paid)
	o.dispatchCheckout()
}

func (o *Order) TotalAmount() uint {
	var totalAmount uint = 0
	for _, item := range o.Items {
		totalAmount += item.Price * item.Quantity
	}

	return totalAmount
}

func (o *Order) IsFree() bool {
	return o.TotalAmount() == 0
}

func (o *Order) setState(state int) {
	o.LastState = state
	o.States = append(o.States, OrderState{State: state})
}

func (oi *OrderItem) setState(state int) {
	oi.LastState = state
	oi.States = append(oi.States, OrderItemState{State: state})
}

func (oi *OrderItem) addErrors(errors *[]ErrorModel) {
	for _, e := range *errors {
		if e.ReferenceCode == oi.ReferenceCode {
			oi.Errors = append(oi.Errors, OrderItemError{Error: e.Error})
		}
	}
}

func (oi *OrderItem) checkOut(errors *[]ErrorModel) {
	if len(*errors) == 0 {
		oi.setState(Processed)
		return
	}

	oi.setState(ProcessFailed)
	oi.addErrors(errors)
}

func (oi *OrderItem) failedToProcess() {
	oi.setState(ProcessFailed)
}

func init() {
	common.AutoMigrate(
		&Order{},
		&OrderState{},
		&OrderItem{},
		&OrderItemError{},
		&OrderItemState{})
}
