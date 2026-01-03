package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

type Loader interface {
	LoadSingleFile(filename string) ([]Bar, error)
}

type ldr struct{}

func NewLoader() Loader {
	return &ldr{}

}

func (l *ldr) LoadSingleFile(filename string) ([]Bar, error) {
	p := &errParser{}
	var bars []Bar
	path := filepath.Join("internal", "historial_data", "archive", filename)
	file, err := os.Open(path)
	if err != nil {
		fmt.Print("Error opening File\n")
		return []Bar{}, ErrLoadingFile
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Print("Error with records\n")
		return []Bar{}, ErrWithRecords
	}
	for i := len(records) - 1; i > 0; i-- {
		record := records[i]
		if i == 0 {
			continue
		}
		date := record[0]
		open := p.parseFloat(record[1])
		high := p.parseFloat(record[2])
		low := p.parseFloat(record[3])
		closeVal := p.parseFloat(record[4])
		volume := p.parseFloat(record[5])
		openInterest := p.parseFloat(record[6])
		newBar := Bar{
			Date:         date,
			Open:         open,
			High:         high,
			Low:          low,
			Close:        closeVal,
			Volume:       volume,
			OpenInterest: openInterest,
		}
		bars = append(bars, newBar)
	}
	return bars, nil

}
