package payment

import (
	"github.com/moxeed/store/common"
	"time"
)

func pollingOpenPaymentFollowUp() {
	var terminals []Terminal

	common.DB.
		Where(&Terminal{LastState: Open}).
		Where("UpdatedAt < ?", time.Now().Add(-15*time.Minute)).
		Limit(10).
		Find(&terminals)

	for _, terminal := range terminals {
		err := Verify(terminal.TerminalCode)
		common.Log(err)
	}

	time.Sleep(time.Minute * common.Configuration.Job.OpenTerminalRetryInMinutes)
}

func init() {
	go pollingOpenPaymentFollowUp()
}
