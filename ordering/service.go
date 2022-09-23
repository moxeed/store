package ordering

import (
	"fmt"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/ordering/ordering_model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

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

func GetBasket(customerCode uint) (ordering_model.OrderModel, error) {
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

func getOrderByIdentifier(identifier ordering_model.OrderIdentifier) (Order, error) {

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

func GetOrder(identifier ordering_model.OrderIdentifier) (ordering_model.OrderModel, error) {
	order, err := getOrderByIdentifier(identifier)
	return order.toModel(), err
}

func AddItem(model ordering_model.AddItemModel) (orderModel ordering_model.OrderModel, err error) {
	order := NewOrder(model.UserCode, model.CustomerCode, nil)
	common.DB.Preload(clause.Associations).Where(&order).Find(&order)

	err = order.addItem(model.ProductCode, model.ReferenceCode, model.Quantity)
	common.DB.Save(&order)

	orderModel = order.toModel()
	return
}

func StartPayment(identifier ordering_model.OrderIdentifier) (orderModel ordering_model.OrderModel, isFree bool, err error) {
	order, err := getOrderByIdentifier(identifier)
	isFree = false

	if err != nil {
		return
	}

	orderModel = order.toModel()
	isFree = order.IsFree()

	order.lockForPayment()
	save(&order)
	return
}

func FlashBuy(model ordering_model.FlashBuyModel) (orderModel ordering_model.OrderModel, isFree bool, err error) {
	order := NewOrder(model.UserCode, model.CustomerCode, &model.OrderReferenceCode)

	err = order.addItem(model.ProductCode, model.ReferenceCode, model.Quantity)
	if err != nil {
		return
	}

	isOk := order.lockForPayment()
	save(&order)

	if !isOk {
		err = fmt.Errorf("وضعیت سبد معتبر نمی باشد")
	}

	isFree = order.IsFree()
	orderModel = order.toModel()
	return
}

func GetOrderList(customerCode uint, offset int) (result []*ordering_model.OrderHeaderModel, totalCount int64) {
	orders := make([]*Order, 0)
	common.DB.Where(&Order{CustomerCode: customerCode}).
		Order("LastState asc, CreatedAt desc").
		Offset(offset).
		Limit(10).
		Find(&orders)

	common.DB.Where(&Order{CustomerCode: customerCode}).
		Count(&totalCount)

	for _, order := range orders {
		result = append(result, order.toHeaderModel())
	}

	return
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

func (os *OrderState) toModel() ordering_model.StateModel {
	return ordering_model.StateModel{
		State:      os.State,
		StateTitle: stateText(os.State),
		CreatedAt:  os.CreatedAt,
	}
}

func (ois *OrderItemState) toModel() ordering_model.StateModel {
	return ordering_model.StateModel{
		State:      ois.State,
		StateTitle: stateText(ois.State),
		CreatedAt:  ois.CreatedAt,
	}
}

func (oi *OrderItem) toModel() ordering_model.OrderItemModel {
	var errors []string
	var states []ordering_model.StateModel

	for _, e := range oi.Errors {
		errors = append(errors, e.Error)
	}

	for _, s := range oi.States {
		states = append(states, s.toModel())
	}

	return ordering_model.OrderItemModel{
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

func (o *Order) toModel() ordering_model.OrderModel {
	var items []ordering_model.OrderItemModel
	var states []ordering_model.StateModel

	for _, item := range o.Items {
		items = append(items, item.toModel())
	}

	for _, state := range o.States {
		states = append(states, state.toModel())
	}

	return ordering_model.OrderModel{
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

func (o *Order) toHeaderModel() *ordering_model.OrderHeaderModel {
	var states []ordering_model.StateModel

	for _, state := range o.States {
		states = append(states, state.toModel())
	}

	return &ordering_model.OrderHeaderModel{
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
