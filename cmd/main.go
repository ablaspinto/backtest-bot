package main

import (
	"cmd/internal/data"
	"cmd/internal/indicators"
	"fmt"
)

func main() {
	loader := data.NewLoader()
	bars, err := loader.LoadSingleFile("CME_ESH2001.csv")
	if err != nil {
		fmt.Printf("error loading bars")
	}
	smaObj := indicators.SMA(bars, 10)
	fmt.Printf("SMA BARS: %v\n", smaObj)
	emaObj := indicators.EMA(bars, 10)
	fmt.Printf("EMA BARS: %v\n", emaObj)
	rsiObj := indicators.RSI(bars, 14)
	fmt.Printf("RSI : %v\n", rsiObj)

	macd := indicators.MACD(bars)
	fmt.Printf("MACD HISTOGRAM: %v\n", macd)

}
