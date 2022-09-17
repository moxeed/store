package payment

import (
	"github.com/moxeed/store/common"
	"gorm.io/gorm"
	"time"
)

const (
	Open = iota + 1
	Failed
	Paid
)

type SettleAble struct {
	CustomerCode uint
	Amount       int
	SettlementID *uint
	Settlement   *Settlement
}

func (s *SettleAble) settle() {
	s.Settlement = &Settlement{
		Amount:       s.Amount,
		CustomerCode: s.CustomerCode,
	}
}

type Settlement struct {
	ID           uint
	Amount       int
	CustomerCode uint
	CreatedAt    time.Time
}

type OrderPayment struct {
	gorm.Model
	SettleAble
	OrderCode uint
	LastState int
	States    []OrderPaymentState
}

type OrderPaymentState struct {
	ID             uint
	State          int
	CreatedAt      time.Time
	OrderPaymentID uint
}

type Terminal struct {
	gorm.Model
	SettleAble
	TerminalCode   string
	LastState      int
	CardPan        string
	ReferenceCode  string
	Fee            int
	OrderPaymentID *uint
	OrderPayment   *OrderPayment
	States         []TerminalState
}

type TerminalState struct {
	ID         uint
	State      int
	CreatedAt  time.Time
	TerminalID uint
}

func NewOrderPayment(customerCode uint, orderCode uint, amount uint) *OrderPayment {
	payment := OrderPayment{OrderCode: orderCode, SettleAble: SettleAble{
		CustomerCode: customerCode,
		Amount:       -int(amount),
	}}
	payment.setState(Open)
	payment.settle()
	return &payment
}

func (op *OrderPayment) registerTerminal(code string) *Terminal {
	terminal := Terminal{
		SettleAble: SettleAble{
			CustomerCode: op.CustomerCode,
			Amount:       -op.Amount,
		},
		TerminalCode:   code,
		OrderPaymentID: &op.ID,
		OrderPayment:   op,
	}

	terminal.setState(Open)

	return &terminal
}

func (op *OrderPayment) setState(state int) {
	op.LastState = state
	op.States = append(op.States, OrderPaymentState{
		State:          state,
		OrderPaymentID: op.ID,
	})
}

func (op *OrderPayment) isPaid() bool {
	return op.LastState == Paid
}

func (t *Terminal) verify() bool {
	if t.LastState != Open {
		return true
	}

	result := verify(t.TerminalCode, t.Amount)
	//result := VerifyResult{
	//	isRepeated:    false,
	//	isOk:          true,
	//	CardPan:       "8912392*****72823",
	//	ReferenceCode: "1239223",
	//	Fee:           0,
	//}

	if !result.isOk && !result.isRepeated {
		t.setState(Failed)
		return false
	}

	t.CardPan = result.CardPan
	t.ReferenceCode = result.ReferenceCode
	t.Fee = result.Fee
	t.setState(Paid)
	t.settle()

	t.OrderPayment.setState(Paid)
	return true
}

func (t *Terminal) setState(state int) {
	t.LastState = state
	t.States = append(t.States, TerminalState{
		State:      state,
		TerminalID: t.ID,
	})
}

func init() {
	common.AutoMigrate(
		&OrderPayment{},
		&OrderPaymentState{},
		&Terminal{},
		&TerminalState{},
		&Settlement{})
}
