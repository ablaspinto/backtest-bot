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
