package userconfigs

type createConfigParams struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Strategy           string   `json:"strategy"`
	Volatility         string   `json:"volatility"`
	Drift              string   `json:"drift"`
	StartingPrice      string   `json:"starting_price"`
	CandleInterval     string   `json:"candle_interval"`
	Symbol             string   `json:"symbol"`
	MarketBehavior     string   `json:"market_behavior"`
	TrendStrength      string   `json:"trend_strength"`
	MeanReversionSpeed string   `json:"mean_reversion_speed"`
	AdvancedParams     []byte   `json:"advanced_params"`
	Tags               []string `json:"tags"`
}

type listConfigsPaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type searchConfigsByNameParams struct {
	Column1 string `json:"column_1"`
	Limit   int32  `json:"limit"`
}

type getConfigWithTagsParams struct {
	Column1 []string `json:"column_1"`
	Limit   int32    `json:"limit"`
}

type updateConfigParams struct {
	ID                 int64    `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Strategy           string   `json:"strategy"`
	Volatility         string   `json:"volatility"`
	Drift              string   `json:"drift"`
	StartingPrice      string   `json:"starting_price"`
	CandleInterval     string   `json:"candle_interval"`
	Symbol             string   `json:"symbol"`
	MarketBehavior     string   `json:"market_behavior"`
	TrendStrength      string   `json:"trend_strength"`
	MeanReversionSpeed string   `json:"mean_reversion_speed"`
	AdvancedParams     []byte   `json:"advanced_params"`
	Tags               []string `json:"tags"`
}

type updateConfigPartialParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Volatility  string `json:"volatility"`
	Drift       string `json:"drift"`
	ID          int64  `json:"id"`
}

type addConfigTagParams struct {
	ID          int64       `json:"id"`
	ArrayAppend interface{} `json:"array_append"`
}

type removeConfigTagParams struct {
	ID          int64       `json:"id"`
	ArrayRemove interface{} `json:"array_remove"`
}
type duplicateConfigParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
