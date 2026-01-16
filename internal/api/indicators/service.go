package indicators

import (
	repo "cmd/internal/adapters/postgresql/sqlc"
	"cmd/internal/dberrors"
	"cmd/internal/pgconverter"
	"context"
)

type Service interface {
	CreateIndicator(ctx context.Context, params createIndicatorParams) (repo.Indicator, error)
	CreateIndicatorBatch(ctx context.Context) error
	GetIndicatorByID(ctx context.Context, id int64) (repo.Indicator, error)
	GetIndicatorsByCandle(ctx context.Context, candleId int64) ([]repo.Indicator, error)
	GetIndicatorByType()
	GetIndicatorsForCandles()
	GetSMAIndicators()
	GetEMAIndicators()
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

func (s *svc) CreateIndicatorBatch(ctx context.Context, params createIndicatorBatchParams) error {
	err := s.repo.CreateIndicatorBatch(ctx, repo.CreateIndicatorBatchParams{
		Signal:        pgconverter.PGTextConv(params.Signal),
		Period:        pgconverter.Int32toInt4(params.Period),
		IndicatorType: params.IndicatorType,
		CandleID:      params.CandleID,
		Value:         pgconverter.PGNumericConverter(params.Value),
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}
