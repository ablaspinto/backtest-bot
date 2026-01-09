package candles

import (
	repo "cmd/internal/adapters/postgresql/sqlc"
	"cmd/internal/dberrors"
	"cmd/internal/pgconverter"
	"context"
)

type Service interface {
	CreateCandle(ctx context.Context, params createCandleParams) (repo.Candle, error)
	CreateCandleBatch(ctx context.Context, params []createCandleParams) (int64, error)
	GetCandleByID(ctx context.Context, id int64) (repo.Candle, error)
	GetCandle(ctx context.Context, params getCandleParams) (repo.Candle, error)
	GetLatestCandle(ctx context.Context, params getLatestCandleParams) (repo.Candle, error)
	GetRecentCandles(ctx context.Context, params getRecentCandleParams) ([]repo.Candle, error)
	GetCandlesInRange(ctx context.Context, params getCandlesInRangeParams) ([]repo.Candle, error)
	GetCandlesAfterTimestamp(ctx context.Context, params getCandlesAfterTimestampParams) ([]repo.Candle, error)
	GetCandlesBeforeTimestamp(ctx context.Context, params getCandlesBeforeTimestampParams) ([]repo.Candle, error)
	ListCandlesPaginated(ctx context.Context, params listCandlesPaginatedParams) ([]repo.Candle, error)
	GetCandlesBySymbol(ctx context.Context, params getCandlesBySymbolParams) ([]repo.Candle, error)
	CountCandles(ctx context.Context, params countCandlesParams) (int64, error)
	CountCandlesInRange(ctx context.Context, params countCandlesInRangeParams) (int64, error)
	GetCandleStats(ctx context.Context, params getCandleStatsParams) (repo.GetCandleStatsRow, error)
	GetCandleStatsInRange(ctx context.Context, params getCandleStatsInRangeParams) (repo.GetCandleStatsInRangeRow, error)
	GetVolumeLeaders(ctx context.Context, params getVolumeLeadersParams) ([]repo.GetVolumeLeadersRow, error)
	UpdateCandle(ctx context.Context, params updateCandleParams) (repo.Candle, error)
	DeleteCandle(ctx context.Context, id int64) error
	DeleteCandlesBySymbol(ctx context.Context, symbol string) error
	DeleteCandlesByTimeFrame(ctx context.Context, params deleteCandlesByTimeFrameParams) error
	DeleteOldCandles(ctx context.Context, timestamps int64) error
	DeleteCandlesInRange(ctx context.Context, params deleteCandlesInRangeParams) error
	GetDistinctSymbols(ctx context.Context, symbol string) ([]string, error)
	GetDistinctTimeFrames(ctx context.Context, symbol string) ([]string, error)
	GetSymbolTimeframePairs(ctx context.Context) ([]repo.GetSymbolTimeframePairsRow, error)
	CheckCandleExists(ctx context.Context, params checkCandleExistsParams) (bool, error)
	GetOldestCandle(ctx context.Context, params getOldsCandleParams) (repo.Candle, error)
	GetCandleGaps(ctx context.Context, params getCandleGapsParams) ([]repo.GetCandleGapsRow, error)
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateCandle(ctx context.Context, params createCandleParams) (repo.Candle, error) {
	candle, err := s.repo.CreateCandle(ctx, repo.CreateCandleParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
		Low:       pgconverter.PGNumericConverter(params.Low),
		High:      pgconverter.PGNumericConverter(params.High),
		Open:      pgconverter.PGNumericConverter(params.Open),
		Close:     pgconverter.PGNumericConverter(params.Close),
		Timestamp: params.Timestamp,
		Volume:    pgconverter.PGNumericConverter(params.Volume),
	})
	if err != nil {
		return repo.Candle{}, dberrors.Handle(err)
	}
	return candle, nil
}
func (s *svc) CreateCandleBatch(ctx context.Context, params []createCandleParams) (int64, error) {
	candleBatch := make([]repo.CreateCandleBatchParams, len(params))
	for i, c := range params {
		candleBatch[i] = repo.CreateCandleBatchParams{
			Symbol:    c.Symbol,
			Timeframe: c.Timeframe,
			Timestamp: c.Timestamp,
			Open:      pgconverter.PGNumericConverter(c.Open),
			High:      pgconverter.PGNumericConverter(c.High),
			Low:       pgconverter.PGNumericConverter(c.Low),
			Close:     pgconverter.PGNumericConverter(c.Close),
			Volume:    pgconverter.PGNumericConverter(c.Volume),
		}
	}
	num, err := s.repo.CreateCandleBatch(ctx, candleBatch)
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return num, nil
}

func (s *svc) GetCandleByID(ctx context.Context, id int64) (repo.Candle, error) {
	c, err := s.repo.GetCandleByID(ctx, id)
	if err != nil {
		return repo.Candle{}, dberrors.Handle(err)
	}
	return c, nil
}
func (s *svc) GetCandle(ctx context.Context, params getCandleParams) (repo.Candle, error) {
	c, err := s.repo.GetCandle(ctx, repo.GetCandleParams{
		Timeframe: params.Timeframe,
		Symbol:    params.Symbol,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return repo.Candle{}, dberrors.Handle(err)
	}
	return c, nil
}
func (s *svc) GetLatestCandle(ctx context.Context, params getLatestCandleParams) (repo.Candle, error) {
	c, err := s.repo.GetLatestCandle(ctx, repo.GetLatestCandleParams{
		Timeframe: params.Timeframe,
		Symbol:    params.Symbol,
	})
	if err != nil {
		return repo.Candle{}, dberrors.Handle(err)
	}
	return c, nil
}
func (s *svc) GetRecentCandles(ctx context.Context, params getRecentCandleParams) ([]repo.Candle, error) {
	candles, err := s.repo.GetRecentCandles(ctx, repo.GetRecentCandlesParams{
		Timeframe: params.Timeframe,
		Symbol:    params.Symbol,
		Limit:     params.Limit,
	})
	if err != nil {
		return []repo.Candle{}, dberrors.Handle(err)
	}
	return candles, nil
}
func (s *svc) GetCandlesInRange(ctx context.Context, params getCandlesInRangeParams) ([]repo.Candle, error) {
	candles, err := s.repo.GetCandlesInRange(ctx, repo.GetCandlesInRangeParams{
		Symbol:      params.Symbol,
		Timeframe:   params.Timeframe,
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
	})
	if err != nil {
		return []repo.Candle{}, dberrors.Handle(err)
	}
	return candles, nil
}
func (s *svc) GetCandlesAfterTimestamp(ctx context.Context, params getCandlesAfterTimestampParams) ([]repo.Candle, error) {
	candles, err := s.repo.GetCandlesAfterTimestamp(ctx, repo.GetCandlesAfterTimestampParams{
		Timeframe: params.Timeframe,
		Symbol:    params.Symbol,
		Timestamp: params.Timestamp,
		Limit:     params.Limit,
	})
	if err != nil {
		return []repo.Candle{}, dberrors.Handle(err)
	}
	return candles, nil
}
func (s *svc) GetCandlesBeforeTimestamp(ctx context.Context, params getCandlesBeforeTimestampParams) ([]repo.Candle, error) {
	candles, err := s.repo.GetCandlesBeforeTimestamp(ctx, repo.GetCandlesBeforeTimestampParams{
		Timeframe: params.Timeframe,
		Symbol:    params.Symbol,
		Timestamp: params.Timestamp,
		Limit:     params.Limit,
	})
	if err != nil {
		return []repo.Candle{}, dberrors.Handle(err)
	}
	return candles, nil

}
func (s *svc) ListCandlesPaginated(ctx context.Context, params listCandlesPaginatedParams) ([]repo.Candle, error) {
	candles, err := s.repo.ListCandlesPaginated(ctx, repo.ListCandlesPaginatedParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
		Limit:     params.Limit,
		Offset:    params.Offset,
	})
	if err != nil {
		return []repo.Candle{}, dberrors.Handle(err)
	}
	return candles, nil

}
func (s *svc) GetCandlesBySymbol(ctx context.Context, params getCandlesBySymbolParams) ([]repo.Candle, error) {
	candles, err := s.repo.GetCandlesBySymbols(ctx, repo.GetCandlesBySymbolsParams{
		Column1:     params.Column1,
		Timeframe:   params.Timeframe,
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
	})
	if err != nil {
		return []repo.Candle{}, dberrors.Handle(err)
	}
	return candles, nil

}
func (s *svc) CountCandles(ctx context.Context, params countCandlesParams) (int64, error) {
	numOfCandles, err := s.repo.CountCandles(ctx, repo.CountCandlesParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
	})
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return numOfCandles, nil

}
func (s *svc) CountCandlesInRange(ctx context.Context, params countCandlesInRangeParams) (int64, error) {
	numOfCandles, err := s.repo.CountCandlesInRange(ctx, repo.CountCandlesInRangeParams{
		Symbol:      params.Symbol,
		Timeframe:   params.Timeframe,
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
	})
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return numOfCandles, nil

}
func (s *svc) GetCandleStats(ctx context.Context, params getCandleStatsParams) (repo.GetCandleStatsRow, error) {
	c, err := s.repo.GetCandleStats(ctx, repo.GetCandleStatsParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
	})
	if err != nil {
		return repo.GetCandleStatsRow{}, dberrors.Handle(err)
	}
	return c, nil
}
func (s *svc) GetCandleStatsInRange(ctx context.Context, params getCandleStatsInRangeParams) (repo.GetCandleStatsInRangeRow, error) {
	candles, err := s.repo.GetCandleStatsInRange(ctx, repo.GetCandleStatsInRangeParams{
		Symbol:      params.Symbol,
		Timeframe:   params.Timeframe,
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
	})
	if err != nil {
		return repo.GetCandleStatsInRangeRow{}, dberrors.Handle(err)
	}
	return candles, nil
}
func (s *svc) GetVolumeLeaders(ctx context.Context, params getVolumeLeadersParams) ([]repo.GetVolumeLeadersRow, error) {
	candles, err := s.repo.GetVolumeLeaders(ctx, repo.GetVolumeLeadersParams{
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
		Limit:       params.Limit,
	})
	if err != nil {
		return []repo.GetVolumeLeadersRow{}, dberrors.Handle(err)
	}
	return candles, nil
}
func (s *svc) UpdateCandle(ctx context.Context, params updateCandleParams) (repo.Candle, error) {
	c, err := s.repo.UpdateCandle(ctx, repo.UpdateCandleParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
		Open:      pgconverter.PGNumericConverter(params.Open),
		High:      pgconverter.PGNumericConverter(params.High),
		Low:       pgconverter.PGNumericConverter(params.Low),
		Close:     pgconverter.PGNumericConverter(params.Close),
		Volume:    pgconverter.PGNumericConverter(params.Volume),
	})
	if err != nil {
		return repo.Candle{}, dberrors.Handle(err)
	}
	return c, nil
}
func (s *svc) DeleteCandle(ctx context.Context, id int64) error {
	err := s.repo.DeleteCandle(ctx, id)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}
func (s *svc) DeleteCandlesBySymbol(ctx context.Context, symbol string) error {
	err := s.repo.DeleteCandlesBySymbol(ctx, symbol)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil

}
func (s *svc) DeleteCandlesByTimeFrame(ctx context.Context, params deleteCandlesByTimeFrameParams) error {
	err := s.repo.DeleteCandlesByTimeframe(ctx, repo.DeleteCandlesByTimeframeParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil

}
func (s *svc) DeleteOldCandles(ctx context.Context, timestamps int64) error {
	err := s.repo.DeleteOldCandles(ctx, timestamps)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}
func (s *svc) DeleteCandlesInRange(ctx context.Context, params deleteCandlesInRangeParams) error {
	err := s.repo.DeleteCandlesInRange(ctx, repo.DeleteCandlesInRangeParams{
		Symbol:      params.Symbol,
		Timeframe:   params.Timeframe,
		Timestamp_2: params.Timestamp_2,
		Timestamp:   params.Timestamp,
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil

}
func (s *svc) GetDistinctSymbols(ctx context.Context, symbol string) ([]string, error) {
	candles, err := s.repo.GetDistinctSymbols(ctx)
	if err != nil {
		return []string{}, dberrors.Handle(err)
	}
	return candles, nil

}

func (s *svc) GetDistinctTimeFrames(ctx context.Context, symbol string) ([]string, error) {
	tf, err := s.repo.GetDistinctTimeframes(ctx, symbol)
	if err != nil {
		return []string{}, dberrors.Handle(err)
	}
	return tf, nil

}
func (s *svc) GetSymbolTimeframePairs(ctx context.Context) ([]repo.GetSymbolTimeframePairsRow, error) {
	timeFramPairs, err := s.repo.GetSymbolTimeframePairs(ctx)
	if err != nil {
		return []repo.GetSymbolTimeframePairsRow{}, dberrors.Handle(err)
	}
	return timeFramPairs, nil

}
func (s *svc) CheckCandleExists(ctx context.Context, params checkCandleExistsParams) (bool, error) {
	c, err := s.repo.CheckCandleExists(ctx, repo.CheckCandleExistsParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return false, dberrors.Handle(err)
	}
	return c, nil

}
func (s *svc) GetOldestCandle(ctx context.Context, params getOldsCandleParams) (repo.Candle, error) {
	c, err := s.repo.GetOldestCandle(ctx, repo.GetOldestCandleParams{
		Symbol:    params.Symbol,
		Timeframe: params.Timeframe,
	})
	if err != nil {
		return repo.Candle{}, dberrors.Handle(err)
	}
	return c, nil
}
func (s *svc) GetCandleGaps(ctx context.Context, params getCandleGapsParams) ([]repo.GetCandleGapsRow, error) {
	candles, err := s.repo.GetCandleGaps(ctx, repo.GetCandleGapsParams{
		Symbol:      params.Symbol,
		Timeframe:   params.Timeframe,
		Timestamp:   params.Timestamp,
		Timestamp_2: params.Timestamp_2,
	})
	if err != nil {
		return []repo.GetCandleGapsRow{}, dberrors.Handle(err)
	}
	return candles, nil
}
