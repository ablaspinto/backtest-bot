package data

import (
	"errors"
	"strconv"
)

var (
	ErrLoadingFile      = errors.New("Error Loading File")
	ErrWithRecords      = errors.New("Error Reading CSV Records")
	ErrConvertingNumber = errors.New("Error converting Record to Float")
)

type errParser struct {
	err error
}

func (p *errParser) parseFloat(s string) float32 {
	if p.err != nil {
		return 0
	}
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		p.err = err
	}
	return float32(val)

}
