package pgconverter

import (
	"math/big"
	"time"

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

func PGInt4Converter(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}

}

func PGTimeStampConverter(s string) pgtype.Timestamptz {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return pgtype.Timestamptz{Valid: false}
	}
	ts := pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
	return ts
}

func PGToBool(s pgtype.Bool) bool {
	return s.Bool
}

func Int4ToInt32(num pgtype.Int4) int32 {
	return num.Int32

}

func Int32toInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Valid: true, Int32: i}
}
