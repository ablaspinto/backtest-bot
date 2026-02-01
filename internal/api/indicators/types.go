package indicators

type createIndicatorParams struct {
	CandleID           int64  `json:"candle_id"`
	IndicatorType      string `json:"indicator_type"`
	Value              string `json:"value"`
	ValueUpper         string `json:"value_upper"`
	ValueLower         string `json:"value_lower"`
	ValueMiddle        string `json:"value_middle"`
	SignalLine         string `json:"signal_line"`
	Histogram          string `json:"histogram"`
	Period             int32  `json:"period"`
	PeriodFast         int32  `json:"period_fast"`
	PeriodSlow         int32  `json:"period_slow"`
	PeriodSignal       int32  `json:"period_signal"`
	StandardDeviations string `json:"standard_deviations"`
	Signal             string `json:"signal"`
	SignalStrength     string `json:"signal_strength"`
	IsCrossover        bool   `json:"is_crossover"`
	CrossoverType      string `json:"crossover_type"`
	Metadata           []byte `json:"metadata"`
}

type createIndicatorBatchParams struct {
	CandleID      int64  `json:"candle_id"`
	IndicatorType string `json:"indicator_type"`
	Value         string `json:"value"`
	Period        int32  `json:"period"`
	Signal        string `json:"signal"`
}

type getIndicatorByTypeParams struct {
	CandleID      int64  `json:"candle_id"`
	IndicatorType string `json:"indicator_type"`
	Period        int32  `json:"period"`
}

type getIndicatorsByTypeParams struct {
	IndicatorType string `json:"indicator_type"`
	Limit         int32  `json:"limit"`
}

type getIndicatorsForCandlesParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type getSMAIndicatorsParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Period      int32  `json:"period"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type getEMAIndicatorsParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Period      int32  `json:"period"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type getRSIIndicatorsParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Period      int32  `json:"period"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type getMACDIndicatorsParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type getBollingerBandsParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Period      int32  `json:"period"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type updateIndicatorParams struct {
	ID             int64  `json:"id"`
	Value          string `json:"value"`
	ValueUpper     string `json:"value_upper"`
	ValueLower     string `json:"value_lower"`
	Signal         string `json:"signal"`
	SignalStrength string `json:"signal_strength"`
}

type updateIndicatorSignalParams struct {
	ID             int64  `json:"id"`
	Signal         string `json:"signal"`
	SignalStrength string `json:"signal_strength"`
}

type markCrossoverParams struct {
	ID            int64  `json:"id"`
	CrossoverType string `json:"crossover_type"`
}

type getBuySignalsParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}

type getSellSignalsParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}

type getOverboughtIndicatorsParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}

type getOversoldIndicatorsParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}

type getSignalsByStrengthParams struct {
	Symbol         string `json:"symbol"`
	Timeframe      string `json:"timeframe"`
	SignalStrength string `json:"signal_strength"`
	Limit          int32  `json:"limit"`
}

type getCrossoversParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
}

type getGoldenCrossesParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}

type getDeathCrossesParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}

type getRecentCrossoversParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Limit     int32  `json:"limit"`
}

type getRsiExtremesParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Value     string `json:"value"`
	Value_2   string `json:"value_2"`
	Timestamp int64  `json:"timestamp"`
}
type getBollingerBreakoutsParams struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Timestamp int64  `json:"timestamp"`
	Limit     int32  `json:"limit"`
}

type getMAConvergenceParams struct {
	Symbol      string `json:"symbol"`
	Timeframe   string `json:"timeframe"`
	Period      int32  `json:"period"`
	Period_2    int32  `json:"period_2"`
	Timestamp   int64  `json:"timestamp"`
	Timestamp_2 int64  `json:"timestamp_2"`
}

type checkIndicatorExistsParams struct {
	CandleID      int64  `json:"candle_id"`
	IndicatorType string `json:"indicator_type"`
	Period        int32  `json:"period"`
}

type getLatestIndicatorParams struct {
	Symbol        string `json:"symbol"`
	Timeframe     string `json:"timeframe"`
	IndicatorType string `json:"indicator_type"`
	Period        int32  `json:"period"`
}

type getIndicatorHistoryParams struct {
	Symbol        string `json:"symbol"`
	Timeframe     string `json:"timeframe"`
	IndicatorType string `json:"indicator_type"`
	Limit         int32  `json:"limit"`
}
