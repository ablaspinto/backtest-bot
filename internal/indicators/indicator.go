package indicators

import (
	"cmd/internal/data"
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
