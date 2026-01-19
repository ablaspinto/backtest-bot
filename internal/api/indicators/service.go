package indicators

import (
	repo "cmd/internal/adapters/postgresql/sqlc"
	"cmd/internal/dberrors"
	"cmd/internal/pgconverter"
	"context"
)

type Service interface {
	CreateIndicator(ctx context.Context, params createIndicatorParams) (repo.Indicator, error)
	//	CreateIndicatorBatch(ctx context.Context) (int64 , error)
	GetIndicatorByID(ctx context.Context, id int64) (repo.Indicator, error)
	GetIndicatorsByCandle(ctx context.Context, candleId int64) ([]repo.Indicator, error)
	GetIndicatorByType(ctx context.Context, params getIndicatorByTypeParams) (repo.Indicator, error)
	GetIndicatorsByType(ctx context.Context, params getIndicatorsByTypeParams) ([]repo.Indicator, error)
	GetIndicatorsForCandles(ctx context.Context, params getIndicatorsForCandlesParams) ([]repo.GetIndicatorsForCandlesRow, error)
	GetSMAIndicators(ctx context.Context, params getSMAIndicatorsParams) ([]repo.GetSMAIndicatorsRow, error)
	GetEMAIndicators(ctx context.Context, params getEMAIndicatorsParams) ([]repo.GetEMAIndicatorsRow, error)
	GetRSIIndicators(ctx context.Context, params getRSIIndicatorsParams) ([]repo.GetRSIIndicatorsRow, error)
	GetMACDIndicators(ctx context.Context, params getMACDIndicatorsParams) ([]repo.GetMACDIndicatorsRow, error)
	GetBoillingerBands(ctx context.Context, params getBollingerBandsParams) ([]repo.GetBollingerBandsRow, error)
	UpdateIndicator(ctx context.Context, params updateIndicatorParams) (repo.Indicator, error)
	UpdateIndicatorSignal(ctx context.Context, params updateIndicatorSignalParams) error
	MarkCrossover(ctx context.Context, params markCrossoverParams) error
	DeleteIndicator(ctx context.Context, id int64) error
	DeleteIndicatorsByCandle(ctx context.Context, candleId int64) error
	DeleteIndicatorsByType(ctx context.Context, indicatorType string) error
	DeleteOldIndicators(ctx context.Context, calculatedAt string) error
	GetBuySignals(ctx context.Context, params getBuySignalsParams) ([]repo.GetBuySignalsRow, error)
	GetSellSignals(ctx context.Context, params getSellSignalsParams) ([]repo.GetSellSignalsRow, error)
	GetOverboughtIndicators(ctx context.Context, params getOverboughtIndicatorsParams) ([]repo.GetOverboughtIndicatorsRow, error)
	GetOversoldIndicators(ctx context.Context, params getOversoldIndicatorsParams) ([]repo.GetOversoldIndicatorsRow, error)
	GetSignalsByStrength(ctx context.Context, params getSignalsByStrengthParams) ([]repo.GetSignalsByStrengthRow, error)
	GetCrossovers(ctx context.Context, params getCrossoversParams) ([]repo.GetCrossoversRow, error)
	GetGoldenCrosses(ctx context.Context, params getGoldenCrossesParams) ([]repo.GetGoldenCrossesRow, error)
	GetDeathCrosses(ctx context.Context, params getDeathCrossesParams) ([]repo.GetDeathCrossesRow, error)
	GetRecentCrossovers(ctx context.Context, params getRecentCrossoversParams) ([]repo.GetRecentCrossoversRow, error)
	CountIndicators(ctx context.Context, candleId int64) (int64, error)
	CountIndicatorsByType(ctx context.Context, indicatorType string) (int64, error)
	GetIndicatorStats(ctx context.Context, indicatorType string) (repo.GetIndicatorStatsRow, error)
	GetSignalDistribution(ctx context.Context) ([]repo.GetSignalDistributionRow, error)
	GetCrossoverFrequency(ctx context.Context) ([]repo.GetCrossoverFrequencyRow, error)
	GetRSIExtremes(ctx context.Context, params getRsiExtremesParams) ([]repo.GetRSIExtremesRow, error)
	GetBollingerBreakouts(ctx context.Context, params getBollingerBreakoutsParams) ([]repo.GetBollingerBreakoutsRow, error)
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateIndicator(ctx context.Context, params createIndicatorParams) (repo.Indicator, error) {
	ind, err := s.repo.CreateIndicator(ctx, repo.CreateIndicatorParams{
		IndicatorType: params.IndicatorType,
		CandleID:      params.CandleID,
		Period:        pgconverter.PGInt4Converter(params.Period),
		PeriodFast:    pgconverter.PGInt4Converter(params.PeriodFast),
		Value:         pgconverter.PGNumericConverter(params.Value),
		ValueUpper:    pgconverter.PGNumericConverter(params.ValueUpper),
		ValueLower:    pgconverter.PGNumericConverter(params.ValueLower),
		ValueMiddle:   pgconverter.PGNumericConverter(params.ValueMiddle),
	})
	if err != nil {
		return repo.Indicator{}, dberrors.Handle(err)
	}
	return ind, nil

}

func (s *svc) GetIndicatorByID(ctx context.Context, id int64) (repo.Indicator, error) {
	indicator, err := s.repo.GetIndicatorByID(ctx, id)
	if err != nil {
		return repo.Indicator{}, dberrors.Handle(err)
	}
	return indicator, nil
}

func (s *svc) GetIndicatorsByCandle(ctx context.Context, candleId int64) ([]repo.Indicator, error) {
	indicators, err := s.repo.GetIndicatorsByCandle(ctx, candleId)
	if err != nil {
		return []repo.Indicator{}, dberrors.Handle(err)
	}
	return indicators, nil
}

func (s *svc) GetIndicatorByType(ctx context.Context, params getIndicatorByTypeParams) (repo.Indicator, error) {
	indicatorType, err := s.repo.GetIndicatorByType(ctx, repo.GetIndicatorByTypeParams{
		CandleID:      params.CandleID,
		Period:        pgconverter.Int32toInt4(params.Period),
		IndicatorType: params.IndicatorType,
	})
	if err != nil {
		return repo.Indicator{}, dberrors.Handle(err)
	}
	return indicatorType, nil
}

func (s *svc) GetIndicatorsByType(ctx context.Context, params getIndicatorsByTypeParams) ([]repo.Indicator, error) {
	indicators, err := s.repo.GetIndicatorsByType(ctx, repo.GetIndicatorsByTypeParams{
		IndicatorType: params.IndicatorType,
		Limit:         params.Limit,
	})
	if err != nil {
		return []repo.Indicator{}, dberrors.Handle(err)
	}
	return indicators, nil

}

func (s *svc) GetIndicatorsForCandles(ctx context.Context, params getIndicatorsForCandlesParams) ([]repo.GetIndicatorsForCandlesRow, error) {
	indicators, err := s.repo.GetIndicatorsForCandles(ctx, repo.GetIndicatorsForCandlesParams{
		Symbol:      params.Symbol,
		Timestamp_2: params.Timestamp_2,
		Timeframe:   params.Timeframe,
		Timestamp:   params.Timestamp,
	})
	if err != nil {
		return []repo.GetIndicatorsForCandlesRow{}, dberrors.Handle(err)

	}
	return indicators, nil

}

func (s *svc) GetSMAIndicators(ctx context.Context, params getSMAIndicatorsParams) ([]repo.GetSMAIndicatorsRow, error) {
	smaIndicators, err := s.repo.GetSMAIndicators(ctx, repo.GetSMAIndicatorsParams{
		Period:      pgconverter.Int32toInt4(params.Period),
		Symbol:      params.Symbol,
		Timestamp_2: params.Timestamp_2,
		Timestamp:   params.Timestamp,
		Timeframe:   params.Timeframe,
	})
	if err != nil {
		return []repo.GetSMAIndicatorsRow{}, dberrors.Handle(err)
	}
	return smaIndicators, nil

}

func (s *svc) GetEMAIndicators(ctx context.Context, params getEMAIndicatorsParams) ([]repo.GetEMAIndicatorsRow, error) {
	emaIndicators, err := s.repo.GetEMAIndicators(ctx, repo.GetEMAIndicatorsParams{
		Period:      pgconverter.Int32toInt4(params.Period),
		Timeframe:   params.Timeframe,
		Timestamp_2: params.Timestamp_2,
		Timestamp:   params.Timestamp,
	})
	if err != nil {
		return []repo.GetEMAIndicatorsRow{}, dberrors.Handle(err)
	}
	return emaIndicators, nil

}

func (s *svc) GetRSIIndicators(ctx context.Context, params getRSIIndicatorsParams) ([]repo.GetRSIIndicatorsRow, error) {
	rsiIndicators, err := s.repo.GetRSIIndicators(ctx, repo.GetRSIIndicatorsParams{
		Symbol:      params.Symbol,
		Period:      pgconverter.Int32toInt4(params.Period),
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
		Timeframe:   params.Timeframe,
	})
	if err != nil {
		return []repo.GetRSIIndicatorsRow{}, dberrors.Handle(err)
	}
	return rsiIndicators, nil
}

func (s *svc) GetMACDIndicators(ctx context.Context, params getMACDIndicatorsParams) ([]repo.GetMACDIndicatorsRow, error) {
	macdIndicators, err := s.repo.GetMACDIndicators(ctx, repo.GetMACDIndicatorsParams{
		Symbol:      params.Symbol,
		Timeframe:   params.Timeframe,
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
	})
	if err != nil {
		return []repo.GetMACDIndicatorsRow{}, dberrors.Handle(err)
	}
	return macdIndicators, nil
}

func (s *svc) GetBoillingerBands(ctx context.Context, params getBollingerBandsParams) ([]repo.GetBollingerBandsRow, error) {
	bands, err := s.repo.GetBollingerBands(ctx, repo.GetBollingerBandsParams{
		Symbol:      params.Symbol,
		Timeframe:   params.Timeframe,
		Period:      pgconverter.Int32toInt4(params.Period),
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
	})
	if err != nil {
		return []repo.GetBollingerBandsRow{}, dberrors.Handle(err)
	}
	return bands, nil

}

func (s *svc) UpdateIndicator(ctx context.Context, params updateIndicatorParams) (repo.Indicator, error) {
	updatedIndicator, err := s.repo.UpdateIndicator(ctx, repo.UpdateIndicatorParams{
		ID:             params.ID,
		Value:          pgconverter.PGNumericConverter(params.Value),
		ValueUpper:     pgconverter.PGNumericConverter(params.ValueUpper),
		ValueLower:     pgconverter.PGNumericConverter(params.ValueLower),
		Signal:         pgconverter.PGTextConv(params.Signal),
		SignalStrength: pgconverter.PGNumericConverter(params.SignalStrength),
	})
	if err != nil {
		return repo.Indicator{}, dberrors.Handle(err)
	}
	return updatedIndicator, nil
}

func (s *svc) UpdateIndicatorSignal(ctx context.Context, params updateIndicatorSignalParams) error {
	err := s.repo.UpdateIndicatorSignal(ctx, repo.UpdateIndicatorSignalParams{
		ID:             params.ID,
		Signal:         pgconverter.PGTextConv(params.Signal),
		SignalStrength: pgconverter.PGNumericConverter(params.SignalStrength),
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) MarkCrossover(ctx context.Context, params markCrossoverParams) error {
	err := s.repo.MarkCrossover(ctx, repo.MarkCrossoverParams{
		ID:            params.ID,
		CrossoverType: pgconverter.PGTextConv(params.CrossoverType),
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) DeleteIndicator(ctx context.Context, id int64) error {
	err := s.repo.DeleteIndicator(ctx, id)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) DeleteIndicatorsByCandle(ctx context.Context, candleId int64) error {
	err := s.repo.DeleteIndicatorsByCandle(ctx, candleId)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) DeleteIndicatorsByType(ctx context.Context, indicatorType string) error {
	err := s.repo.DeleteIndicatorsByType(ctx, indicatorType)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) DeleteOldIndicators(ctx context.Context, calculatedAt string) error {
	err := s.repo.DeleteOldIndicators(ctx, pgconverter.PGTimeStampConverter(calculatedAt))
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) GetBuySignals(ctx context.Context, params getBuySignalsParams) ([]repo.GetBuySignalsRow, error) {
	buySignals, err := s.repo.GetBuySignals(ctx, repo.GetBuySignalsParams{
		Limit:     params.Limit,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
		Symbol:    params.Symbol,
	})
	if err != nil {
		return []repo.GetBuySignalsRow{}, dberrors.Handle(err)
	}
	return buySignals, nil
}

func (s *svc) GetSellSignals(ctx context.Context, params getSellSignalsParams) ([]repo.GetSellSignalsRow, error) {
	sellSignals, err := s.repo.GetSellSignals(ctx, repo.GetSellSignalsParams{
		Symbol:    params.Symbol,
		Limit:     params.Limit,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return []repo.GetSellSignalsRow{}, dberrors.Handle(err)
	}
	return sellSignals, nil
}

func (s *svc) GetOverboughtIndicators(ctx context.Context, params getOverboughtIndicatorsParams) ([]repo.GetOverboughtIndicatorsRow, error) {
	overboughIndicators, err := s.repo.GetOverboughtIndicators(ctx, repo.GetOverboughtIndicatorsParams{
		Limit:     params.Limit,
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return []repo.GetOverboughtIndicatorsRow{}, dberrors.Handle(err)
	}
	return overboughIndicators, nil
}

func (s *svc) GetOversoldIndicators(ctx context.Context, params getOversoldIndicatorsParams) ([]repo.GetOversoldIndicatorsRow, error) {
	oversoldIndicators, err := s.repo.GetOversoldIndicators(ctx, repo.GetOversoldIndicatorsParams{
		Limit:     params.Limit,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
		Symbol:    params.Symbol,
	})
	if err != nil {
		return []repo.GetOversoldIndicatorsRow{}, dberrors.Handle(err)
	}
	return oversoldIndicators, nil
}
func (s *svc) GetSignalsByStrength(ctx context.Context, params getSignalsByStrengthParams) ([]repo.GetSignalsByStrengthRow, error) {
	signalsStrength, err := s.repo.GetSignalsByStrength(ctx, repo.GetSignalsByStrengthParams{
		Symbol:         params.Symbol,
		SignalStrength: pgconverter.PGNumericConverter(params.SignalStrength),
		Limit:          params.Limit,
		Timeframe:      params.Timeframe,
	})
	if err != nil {
		return []repo.GetSignalsByStrengthRow{}, dberrors.Handle(err)
	}
	return signalsStrength, nil

}

func (s *svc) GetCrossovers(ctx context.Context, params getCrossoversParams) ([]repo.GetCrossoversRow, error) {
	crossovers, err := s.repo.GetCrossovers(ctx, repo.GetCrossoversParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return []repo.GetCrossoversRow{}, dberrors.Handle(err)
	}
	return crossovers, nil
}

func (s *svc) GetGoldenCrosses(ctx context.Context, params getGoldenCrossesParams) ([]repo.GetGoldenCrossesRow, error) {
	goldenCrosses, err := s.repo.GetGoldenCrosses(ctx, repo.GetGoldenCrossesParams{
		Symbol:    params.Symbol,
		Limit:     params.Limit,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return []repo.GetGoldenCrossesRow{}, dberrors.Handle(err)
	}
	return goldenCrosses, nil
}

func (s *svc) GetDeathCrosses(ctx context.Context, params getDeathCrossesParams) ([]repo.GetDeathCrossesRow, error) {
	deathCrosses, err := s.repo.GetDeathCrosses(ctx, repo.GetDeathCrossesParams{
		Limit:     params.Limit,
		Timeframe: params.Timeframe,
		Symbol:    params.Symbol,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return []repo.GetDeathCrossesRow{}, dberrors.Handle(err)
	}
	return deathCrosses, nil
}

func (s *svc) GetRecentCrossovers(ctx context.Context, params getRecentCrossoversParams) ([]repo.GetRecentCrossoversRow, error) {
	crossovers, err := s.repo.GetRecentCrossovers(ctx, repo.GetRecentCrossoversParams{
		Limit:     params.Limit,
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
	})
	if err != nil {
		return []repo.GetRecentCrossoversRow{}, dberrors.Handle(err)
	}
	return crossovers, nil
}

func (s *svc) CountIndicators(ctx context.Context, candleId int64) (int64, error) {
	num, err := s.repo.CountIndicators(ctx, candleId)
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return num, nil

}

func (s *svc) CountIndicatorsByType(ctx context.Context, indicatorType string) (int64, error) {
	num, err := s.repo.CountIndicatorsByType(ctx, indicatorType)
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return num, nil

}

func (s *svc) GetIndicatorStats(ctx context.Context, indicatorType string) (repo.GetIndicatorStatsRow, error) {
	indicatorStats, err := s.repo.GetIndicatorStats(ctx, indicatorType)
	if err != nil {
		return repo.GetIndicatorStatsRow{}, dberrors.Handle(err)
	}
	return indicatorStats, nil
}

func (s *svc) GetSignalDistribution(ctx context.Context) ([]repo.GetSignalDistributionRow, error) {
	signalDis, err := s.repo.GetSignalDistribution(ctx)
	if err != nil {
		return []repo.GetSignalDistributionRow{}, dberrors.Handle(err)
	}
	return signalDis, nil
}

func (s *svc) GetCrossoverFrequency(ctx context.Context) ([]repo.GetCrossoverFrequencyRow, error) {
	crossoverFreq, err := s.repo.GetCrossoverFrequency(ctx)
	if err != nil {
		return []repo.GetCrossoverFrequencyRow{}, dberrors.Handle(err)
	}
	return crossoverFreq, nil
}
func (s *svc) GetRSIExtremes(ctx context.Context, params getRsiExtremesParams) ([]repo.GetRSIExtremesRow, error) {
	rsiExt, err := s.repo.GetRSIExtremes(ctx, repo.GetRSIExtremesParams{
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
		Value:     pgconverter.PGNumericConverter(params.Value),
		Value_2:   pgconverter.PGNumericConverter(params.Value_2),
		Symbol:    params.Symbol,
	})
	if err != nil {
		return []repo.GetRSIExtremesRow{}, dberrors.Handle(err)
	}
	return rsiExt, nil

}

func (s *svc) GetBollingerBreakouts(ctx context.Context, params getBollingerBreakoutsParams) ([]repo.GetBollingerBreakoutsRow, error) {
	bollingerBreakouts, err := s.repo.GetBollingerBreakouts(ctx, repo.GetBollingerBreakoutsParams{
		Symbol:    params.Symbol,
		Limit:     params.Limit,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return []repo.GetBollingerBreakoutsRow{}, dberrors.Handle(err)
	}
	return bollingerBreakouts, nil

}
