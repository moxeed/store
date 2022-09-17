package ordering

import (
	"github.com/moxeed/store/common"
	"time"
)

func reprocess(items []*OrderItem) {
	dispatchCheckout(items)
}

func pollingTimeBasedFailedOrderRetry() {
	for {
		println("job is running")
		var orderItems []*OrderItem
		common.DB.
			Where(&OrderItem{LastState: ProcessFailed}).
			Where("UpdatedAt < ?", time.Now().Add(-15*time.Minute)).
			Limit(10).
			Find(&orderItems)

		reprocess(orderItems)

		time.Sleep(time.Minute * common.Configuration.Job.FailedOrderRetryInMinutes)
	}
}

func pollingTimeBasedOpenOrderRetry() {
	for {
		println("job is running")
		var orders []Order
		common.DB.
			Where(&Order{LastState: PaymentPending}).
			Where("UpdatedAt < ?", time.Now().Add(-5*time.Minute)).
			Limit(10).
			Find(&orders)

		for _, order := range orders {
			CheckOut(order.ID)
		}

		time.Sleep(time.Minute * common.Configuration.Job.FailedOrderRetryInMinutes)
	}
}

func init() {
	go pollingTimeBasedOpenOrderRetry()
	go pollingTimeBasedFailedOrderRetry()
}
