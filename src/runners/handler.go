package runners

import (
	"fmt"
	"log"
	"nudam-ctrader-api/api"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
	"nudam-ctrader-api/strategy"
	"sync"
	"time"
)

type IHandler interface {
	HandlerReadMessage(api api.CTraderAPI)
	HandlerStrategy()
	HandlerGetTrendbars(api api.CTraderAPI)
}

type Handler struct {
	positions map[string]bool
	mu        sync.Mutex
}

func NewHandler() IHandler {
	handler := new(Handler)
	handler.positions = make(map[string]bool)
	return handler
}

// Func to start goroutine signal checker.
func (h *Handler) HandlerStrategy() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, symbol := range configs_helper.TraderConfiguration.CurrencyPairs {
			signal, err := strategy.SignalChecker(symbol)
			if err != nil {
				logger.LogError(err, "error getting data from mongodb")
				log.Panic(err)
			}

			if signal == strategy.Short {

			} else if signal == strategy.Long {

			}
		}
	}
}

// Func to start goroutine message reader.
func (h *Handler) HandlerReadMessage(api api.CTraderAPI) {
	for {
		if err := api.ReadMessage(); err != nil {
			logger.LogError(err, "error reading message")
		}
	}
}

// Func to start goroutine trendbar receiver.
func (h *Handler) HandlerGetTrendbars(api api.CTraderAPI) {
	var wg sync.WaitGroup
	for _, symbol := range configs_helper.TraderConfiguration.CurrencyPairs {
		for period := range configs_helper.TraderConfiguration.Periods {
			wg.Add(1)
			go func(symbol, period string) {
				ticker := time.NewTicker(30 * time.Second)
				defer ticker.Stop()

				for range ticker.C {
					if err := api.GetTrendbars(symbol, period); err != nil {
						logger.LogError(err, fmt.Sprintf("error getting trendbars for %s", symbol))
						log.Panic(err)
					}
				}
			}(symbol, period)
		}
	}
	wg.Wait()
}

func (h *Handler) isPositionOpen(symbol string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	open, exists := h.positions[symbol]
	return exists && open
}

func (h *Handler) openPosition(symbol string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.positions[symbol] = true
}

func (h *Handler) closePosition(symbol string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.positions[symbol] = false
}
