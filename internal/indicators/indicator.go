package indicators

import (
	"cmd/internal/data"
	"fmt"
	"math"
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

func RSI(bars []data.Bar, period float32) []float32 {
	var rsiValues []float32
	var gains []float32
	var losses []float32

	for i := range len(bars) {
		currPrice := bars[i].Close
		if i == 0 {
			continue
		}
		prevPrice := bars[i-1].Close
		difference := currPrice - prevPrice
		if difference > 0 {
			gains = append(gains, difference)
			losses = append(losses, 0)
		} else if difference < 0 {
			absVal := math.Abs(float64(difference))
			losses = append(losses, float32(absVal))
			gains = append(gains, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, 0)
		}
	}
	simpleGainAvg := float32(0)
	simpleLossAvg := float32(0)
	for i := range int32(period) {
		simpleGainAvg += gains[i]
		simpleLossAvg += losses[i]
	}
	simpleGainAvg = simpleGainAvg / float32(period)
	simpleLossAvg = simpleLossAvg / float32(period)

	RS := simpleGainAvg / simpleLossAvg
	firstRSI := 100 - (100 / (1 + RS))
	rsiValues = append(rsiValues, firstRSI)

	for i := int(period); i < len(gains); i++ {
		currGain := gains[i]
		currLoss := losses[i]
		simpleGainAvg = ((simpleGainAvg)*(period-1) + currGain) / period
		simpleLossAvg = ((simpleLossAvg)*(period-1) + currLoss) / period
		if simpleLossAvg == 0 {
			rsiValues = append(rsiValues, 100)
			continue
		}
		secondRs := simpleGainAvg / simpleLossAvg
		secRSI := 100 - (100 / (1 + secondRs))
		rsiValues = append(rsiValues, secRSI)
	}
	return rsiValues
}
