package common

type PaymentCompleted struct {
	PaymentCode  uint
	TerminalCode uint
	OrderCode    uint
}

type PaymentCompletedListener interface {
	Handle(completed PaymentCompleted)
}

var subscribers []PaymentCompletedListener
var channel chan PaymentCompleted

func dispatcher(c chan PaymentCompleted) {
	for message := range c {
		for _, s := range subscribers {
			s.Handle(message)
		}
	}
}

func Register(l PaymentCompletedListener) {
	subscribers = append(subscribers, l)
}

func Dispatch(e PaymentCompleted) {
	channel <- e
}

func init() {
	channel = make(chan PaymentCompleted)
	go dispatcher(channel)
}
