package ordering

import (
	"fmt"
	"github.com/moxeed/store/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

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

func getBasket(customerCode uint) (Order, error) {
	order := Order{CustomerCode: customerCode, LastState: Basket}
	dbResult := common.DB.
		Preload("Items").
		Preload("States").
		Preload("Items.States").
		Preload("Items.Errors").
		Where(&order).First(&order)

	if dbResult.Error == gorm.ErrRecordNotFound {
		return order, fmt.Errorf("سبد پیدا نشد")
	}

	return order, nil
}

func GetBasket(customerCode uint) (OrderModel, error) {
	order, err := getBasket(customerCode)
	return order.toModel(), err
}

func getOrder(ID uint) (Order, error) {
	order := Order{}
	dbResult := common.DB.
		Preload("Items").
		Preload("States").
		Preload("Items.States").
		Preload("Items.Errors").
		First(&order, ID)

	if dbResult.Error == gorm.ErrRecordNotFound {
		return order, fmt.Errorf("سفارش پیدا نشد")
	}

	return order, nil
}

func getOrderByIdentifier(identifier OrderIdentifier) (Order, error) {

	order := Order{ReferenceCode: identifier.ReferenceCode}
	if identifier.ID != nil {
		order.ID = *identifier.ID
	}

	if identifier.ID == nil && identifier.ReferenceCode == nil {
		return order, fmt.Errorf("سفارش پیدا نشد")
	}

	dbResult := common.DB.
		Preload("Items").
		Preload("States").
		Preload("Items.States").
		Preload("Items.Errors").
		Where(&order).
		First(&order)

	if dbResult.Error == gorm.ErrRecordNotFound {
		return order, fmt.Errorf("سفارش پیدا نشد")
	}

	return order, nil
}

func GetOrder(identifier OrderIdentifier) (OrderModel, error) {
	order, err := getOrderByIdentifier(identifier)
	return order.toModel(), err
}

func AddItem(model AddItemModel) (orderModel OrderModel, err error) {
	order := NewOrder(model.UserCode, model.CustomerCode, nil)
	common.DB.Preload(clause.Associations).Where(&order).Find(&order)

	err = order.addItem(model.ProductCode, model.ReferenceCode, model.Quantity)
	common.DB.Save(&order)

	orderModel = order.toModel()
	return
}

func FlashBuy(model FlashBuyModel) (orderModel OrderModel, err error) {
	order := NewOrder(model.UserCode, model.CustomerCode, &model.OrderReferenceCode)

	err = order.addItem(model.ProductCode, model.ReferenceCode, model.Quantity)
	order.lockForPayment()
	common.DB.Save(&order)

	orderModel = order.toModel()
	return
}

func LockForPayment(order *Order) (OrderModel, bool, error) {
	isOk, isFree := order.lockForPayment()
	save(order)

	if isOk {
		return order.toModel(), isFree, nil
	}
	return order.toModel(), isFree, fmt.Errorf("امکان قفل کردن سفارش وجود ندارد")
}

func StartPayment(identifier OrderIdentifier) (OrderModel, bool, error) {
	order, err := getOrderByIdentifier(identifier)

	if err != nil {
		return order.toModel(), false, err
	}

	if order.LastState == Basket {
		return LockForPayment(&order)
	}

	if order.LastState != PaymentPending {
		return order.toModel(), order.IsFree(), fmt.Errorf("سفارش قبلا پرداخت شده است")
	}

	return order.toModel(), order.IsFree(), nil
}

func GetOrderList(customerCode uint, offset int) ([]*OrderHeaderModel, int64) {
	orders := make([]*Order, 0)
	common.DB.Where(&Order{CustomerCode: customerCode}).
		Order("LastState asc, CreatedAt desc").
		Offset(offset).
		Limit(10).
		Find(&orders)

	var totalCount int64 = 0
	common.DB.Where(&Order{CustomerCode: customerCode}).
		Count(&totalCount)

	result := make([]*OrderHeaderModel, 0)

	for _, order := range orders {
		result = append(result, order.toHeaderModel())
	}

	return result, totalCount
}

func CheckOut(orderID uint) {
	order, err := getOrder(orderID)
	if err != nil {
		log.Println(err)
		return
	}

	order.checkOut()
	save(&order)
}

func save(order *Order) {
	common.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&order)
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
		ProductTitle:   oi.ProductTitle,
		Category:       oi.Category,
		ProductCode:    oi.ProductCode,
		Price:          oi.Price,
		Quantity:       oi.Quantity,
		LastState:      oi.LastState,
		LastStateTitle: stateText(oi.LastState),
		Errors:         errors,
		States:         states,
	}
}

func (o *Order) toModel() OrderModel {
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
		TotalAmount:    o.TotalAmount(),
		Items:          items,
		States:         states,
		CreatedAt:      o.CreatedAt,
	}
}

func (o *Order) toHeaderModel() *OrderHeaderModel {
	var states []StateModel

	for _, state := range o.States {
		states = append(states, state.toModel())
	}

	return &OrderHeaderModel{
		ID:             o.ID,
		UserCode:       o.UserCode,
		CustomerCode:   o.CustomerCode,
		LastState:      o.LastState,
		LastStateTitle: stateText(o.LastState),
		TotalAmount:    o.TotalAmount(),
		States:         states,
		CreatedAt:      o.CreatedAt,
	}
}

type PaymentListener struct{}

func (PaymentListener) Handle(completed common.PaymentCompleted) {
	CheckOut(completed.OrderCode)
}

func init() {
	common.Register(PaymentListener{})
}
