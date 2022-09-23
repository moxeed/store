package payment

import (
	"fmt"
	"github.com/moxeed/store/common"
	"github.com/moxeed/store/payment/payment_model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreatePayment(customerCode uint, orderCode uint, amount uint) {
	payment := NewOrderPayment(customerCode, orderCode, amount)
	common.DB.Create(payment)
}

func OpenTerminal(orderCode uint) (string, error) {
	orderPayment := OrderPayment{OrderCode: orderCode}
	result := common.DB.Where(&orderPayment).First(&orderPayment)

	if result.Error == gorm.ErrRecordNotFound {
		return "", fmt.Errorf("پرداخت پیدا نشد")
	}

	terminalCode, err := openTerminal(orderPayment.Amount, fmt.Sprintf("پرداخت سفارش به شماره %d", orderCode))

	if err == nil {
		terminal := orderPayment.registerTerminal(terminalCode)
		common.DB.Create(terminal)
	}

	return common.Configuration.ZarinPal.RedirectUrl + terminalCode, err
}

func Verify(terminalCode string) (uint, error) {
	terminal := Terminal{TerminalCode: terminalCode}
	dbResult := common.DB.
		Preload(clause.Associations).
		Preload("OrderPayment.States").
		Where(&terminal).First(&terminal)

	if dbResult.Error == gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("درگاه پیدا نشد")
	}

	if terminal.LastState == Paid {
		return terminal.OrderPayment.OrderCode, nil
	}

	if terminal.LastState == Failed {
		return terminal.OrderPayment.OrderCode, fmt.Errorf("پرداخت معتبر نمی باشد")
	}

	result := terminal.verify()

	if !result {
		return terminal.OrderPayment.OrderCode, fmt.Errorf("پرداخت معتبر نمی باشد")
	}

	common.DB.Save(&terminal)
	if terminal.OrderPayment != nil {
		common.DB.Save(terminal.OrderPayment)
	}

	if terminal.OrderPaymentID != nil {
		common.Dispatch(common.PaymentCompleted{
			PaymentCode:  *terminal.OrderPaymentID,
			TerminalCode: terminal.ID,
			OrderCode:    terminal.OrderPayment.OrderCode,
		})
	}

	return terminal.OrderPayment.OrderCode, nil
}

func IsPaid(orderCode uint) payment_model.InquiryModel {
	orderPayment := OrderPayment{OrderCode: orderCode}
	result := common.DB.Where(&orderPayment).First(&orderPayment)

	if result.Error == gorm.ErrRecordNotFound {
		return payment_model.InquiryModel{
			DoesExists: false,
			IsPaid:     false,
		}
	}

	return payment_model.InquiryModel{
		DoesExists: true,
		IsPaid:     orderPayment.isPaid(),
	}
}
