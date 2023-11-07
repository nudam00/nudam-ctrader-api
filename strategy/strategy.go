package strategy

import "log"

func GetSignal(closingPrices []float64) {
	ema26 := calculateEMA(closingPrices, 26)
	ema50 := calculateEMA(closingPrices, 50)

	for i := 50; i < len(closingPrices); i++ {
		log.Printf("%f", ema26[i])
		if ema50[i] > ema26[i] {
			log.Printf("Dzień %d: Trend spadkowy (EMA50 > EMA26)\n", i+1)
		} else {
			log.Printf("Dzień %d: Trend wzrostowy (EMA50 <= EMA26)\n", i+1)
		}
	}
}

func calculateEMA(closingPrices []float64, period int) []float64 {
	alpha := 2.0 / (float64(period) + 1.0)
	emaValues := make([]float64, len(closingPrices))

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += closingPrices[i]
	}
	emaValues[period-1] = sum / float64(period)

	for i := period; i < len(closingPrices); i++ {
		emaValues[i] = (closingPrices[i]-emaValues[i-1])*alpha + emaValues[i-1]
	}

	return emaValues
}
