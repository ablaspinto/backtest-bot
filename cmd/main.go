package main

import (
	"cmd/internal/data"
	"fmt"
)

func main() {
	loader := data.NewLoader()
	bars, err := loader.LoadSingleFile("CME_ESH2000.csv")
	if err != nil {
		fmt.Printf("error loading bars")
	}
	for b := range bars {
		fmt.Printf("BAR: %v", b)
	}

}
