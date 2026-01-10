package sessions

type createSessionParams struct {
	SessionName   string   `json:"session_name"`
	Symbol        string   `json:"symbol"`
	Timeframe     string   `json:"timeframe"`
	Strategy      string   `json:"strategy"`
	StartingPrice string   `json:"starting_price"`
	Parameters    []byte   `json:"parameters"`
	Notes         string   `json:"notes"`
	Tags          []string `json:"tags"`
}

type listSessionParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
type getSessionBySymbolsParams struct {
	Symbol string `json:"symbol"`
	Limit  int32  `json:"limit"`
}
type getSessionByStrategyParams struct {
	Strategy string `json:"strategy"`
	Limit    int32  `json:"limit"`
}
type getSessionByStatusParams struct {
	Status string `json:"status"`
	Limit  int32  `json:"limit"`
}

type searchSessionByNameParams struct {
	Column1 string `json:"column_1"`
	Limit   int32  `json:"limit"`
}

type getSessionWithTagsParams struct {
	Column1 []string `json:"column_1"`
	Limit   int32    `json:"limit"`
}
type updateSessionParams struct {
	ID          int64    `json:"id"`
	SessionName string   `json:"session_name"`
	Notes       string   `json:"notes"`
	Tags        []string `json:"tags"`
	IsFavorite  bool     `json:"is_favorite"`
}

type updateSessionEndParams struct {
	ID                 int64  `json:"id"`
	EndedAt            string `json:"ended_at"`
	EndingPrice        string `json:"ending_price"`
	HighestPrice       string `json:"highest_price"`
	LowestPrice        string `json:"lowest_price"`
	TotalCandles       int32  `json:"total_candles"`
	PriceChangePercent string `json:"price_change_percent"`
	Volatility         string `json:"volatility"`
	Status             string `json:"status"`
}
type updateSessionStatus struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}
type updateSessionPricesParams struct {
	ID           int64  `json:"id"`
	HighestPrice string `json:"highest_price"`
	LowestPrice  string `json:"lowest_price"`
}

type addSessionTagParams struct {
	ID          int64       `json:"id"`
	ArrayAppend interface{} `json:"array_append"`
}
type removeSessionTagParams struct {
	ID          int64       `json:"id"`
	ArrayRemove interface{} `json:"array_remove"`
}
