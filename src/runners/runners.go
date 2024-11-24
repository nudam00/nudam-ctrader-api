package runners

import (
	"log"
	"nudam-ctrader-api/api"
)

type IRunner interface {
	StartRoutines()
}

type Runner struct {
	handler IHandler
}

func NewRunner() IRunner {
	runner := new(Runner)
	return runner
}

// Start trading goroutines.
func (r *Runner) StartRoutines() {
	api := api.NewApi()
	err := api.Open()
	if err != nil {
		log.Panic(err)
	}
	defer api.Close()

	r.handler = NewHandler(api)

	go r.handler.HandlerReadMessage()

	go r.handler.HandlerStrategy()

	r.handler.HandlerGetTrendbars()
}
