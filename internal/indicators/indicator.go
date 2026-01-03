package indicators

import (
	"cmd/internal/data"
	"fmt"
)

func SMA(bars []data.Bar, period int) []float32 {
	var smaAverages []float32
	var currAverage float32

	for i := range len(bars) {
		start := i - period + 1

		if start >= 0 {
			for j := range period {
				closeVal := bars[start+j].Close
				currAverage += closeVal
			}
			currSmaAvg := currAverage / float32(period)
			smaAverages = append(smaAverages, currSmaAvg)
			currAverage = 0
		}
	}
	return smaAverages
}

func EMA(bars []data.Bar, period float32, firstSmaValue float32) []float32 {
	var emaAverages []float32
	// EMA = (priceToday * k) + (emaYesterday * (1 - k))
	// k = 2 / (period + 1)
	k := float32(2) / (period + float32(1))
	j := 0
	fmt.Printf("K: %v\n", k)
	for i := 11; i < len(bars)-1; i++ { // start at 10, since we averaged out for the first sMA value
		closePrice := bars[i].Close
		if i == 11 {
			ema := (closePrice * k) + (firstSmaValue * (1 - k))
			emaAverages = append(emaAverages, ema)
			continue
		}
		ema := (closePrice * k) + (emaAverages[j] * (1 - k))
		emaAverages = append(emaAverages, ema)
		j++
	}

	return emaAverages
}
