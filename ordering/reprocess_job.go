package ordering

import (
	"github.com/moxeed/store/common"
	"time"
)

func reprocess(items *[]OrderItem) {
	dispatchCheckout(items)
}

func pollingTimeBasedRetry() {
	for {
		var orderItems []OrderItem
		common.DB.
			Where(&OrderItem{LastState: ProcessFailed}).
			Where("UpdatedAt < ?", time.Now().Add(-15*time.Minute)).
			Find(&orderItems)

		reprocess(&orderItems)

		time.Sleep(1 * time.Minute)
	}
}

func init() {
	go pollingTimeBasedRetry()
}
