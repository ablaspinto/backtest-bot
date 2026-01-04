package indicators

import (
	"cmd/internal/data"
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

func EMA(bars []data.Bar, period int) []float32 {
	result := make([]float32, len(bars))
	var sum float32
	for i := range period {
		sum += bars[i].Close
		if i < period-1 {
			result[i] = float32(math.NaN())
		}
	}
	result[period-1] = sum / float32(period)

	// Calculate EMA for remaining bars
	k := float32(2) / float32(period+1)
	for i := period; i < len(bars); i++ {
		result[i] = (bars[i].Close * k) + (result[i-1] * (1 - k))
	}

	return result
}

func emaMACD(bars []float32, period int) []float32 {
	result := make([]float32, len(bars))
	var sum float32
	var startIndex int
	// Fill all invalid indices with NaN
	for i := range bars {
		if !math.IsNaN(float64(bars[i])) {
			startIndex = i
			break
		}
	}
	for i := 0; i < startIndex+period-1; i++ {
		result[i] = float32(math.NaN())
	}

	// Calculate SMA seed
	for i := range period {
		sum += bars[i+startIndex]
	}
	result[startIndex+period-1] = sum / float32(period)

	// Calculate EMA for remaining bars
	k := float32(2) / float32(period+1)
	for i := startIndex + period; i < len(bars); i++ {
		result[i] = (bars[i] * k) + (result[i-1] * (1 - k))
	}

	return result
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

func MACD(bars []data.Bar) []float32 {
	var macd []float32
	var histogram []float32
	first := EMA(bars, 12)
	second := EMA(bars, 26)
	for i := range len(first) {
		macd = append(macd, first[i]-second[i])
	}
	signal := emaMACD(macd, 9)
	for i := range len(signal) {
		histogram = append(histogram, macd[i]-signal[i])
	}
	return histogram

}
