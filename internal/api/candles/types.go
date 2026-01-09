package candles

type createCandleParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Open      string `json:"open"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Close     string `json:"close"`
	Volume    string `json:"volume"`
}

type getCandleParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
}

type getLatestCandleParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
}

type getRecentCandleParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Limit     int32  `json:"limit"`
}

type getCandlesInRangeParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type getCandlesAfterTimestampParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}
type getCandlesBeforeTimestampParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}

type listCandlesPaginatedParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type getCandlesBySymbolParams struct {
	Column1     []string `json:"column_1"`
	Timeframe   string   `json:"timeframe"`
	Timestamp   int64    `json:"timestamp"`
	Timestamp_2 int64    `json:"timestamp_2"`
}

type countCandlesParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
}

type countCandlesInRangeParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type getCandleStatsParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
}
type getCandleStatsInRangeParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type getVolumeLeadersParams struct {
	Timestamp   int64 `json:"timestamp"`
	Timestamp_2 int64 `json:"timestamp_2"`
	Limit       int32 `json:"limit"`
}

type updateCandleParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Open      string `json:"open"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Close     string `json:"close"`
	Volume    string `json:"volume"`
}

type deleteCandlesByTimeFrameParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
}
type deleteCandlesInRangeParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}
type checkCandleExistsParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
}

type getOldsCandleParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
}
type getCandleGapsParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}
