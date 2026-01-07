package candles

import (
	repo "cmd/internal/adapters/postgresql/sqlc"
	"context"
)

type Service interface {
	CreateCandle(ctx context.Context, params createCandleParams) (repo.Candle, error)
	CreateCandleBatch(ctx context.Context, params []createCandleParams) (int64, error)
	GetCandleByID(ctx context.Context, id int64) (repo.Candle, error)
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}
