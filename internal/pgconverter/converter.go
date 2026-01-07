package pgconverter

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func PGTextConv(text string) pgtype.Text {
	return pgtype.Text{String: text, Valid: true}

}

func PGBoolConv(val bool) pgtype.Bool {
	return pgtype.Bool{Bool: val, Valid: true}

}

func PGNumericConverter(i string) pgtype.Numeric {
	bigIntVal := new(big.Int)
	bigIntVal.SetString(i, 10)
	return pgtype.Numeric{Valid: true, Int: bigIntVal}

}
