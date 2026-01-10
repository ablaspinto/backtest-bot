package sessions

import (
	repo "cmd/internal/adapters/postgresql/sqlc"
	"cmd/internal/dberrors"
	"cmd/internal/pgconverter"
	"context"
	"runtime/debug"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/benchmark"
)

type Service interface {
	CreateSession(ctx context.Context, params createSessionParams) (repo.Session, error)
	GetSessionByID(ctx context.Context, id int64) (repo.Session, error)
	GetSessionByName(ctx context.Context, sessionName string) (repo.Session, error)
	ListAllSessions(ctx context.Context) ([]repo.Session, error)
	ListSessionsWithLimit(ctx context.Context, limit, offset string) ([]repo.Session, error)
	GetActiveSessions(ctx context.Context)
	GetCompletedSessions(ctx context.Context)
	GetFavoriteSessions(ctx context.Context)
	GetSessionsBySymbol(ctx context.Context)
	GetSessionByStrategy(ctx context.Context)
	GetSessionByStatus(ctx context.Context)
	GetRecentSessions(ctx context.Context)
	SearchSessionByName(ctx context.Context)
	GetSessionWithTags(ctx context.Context) ([]repo.Session, error)
	UpdateSession(ctx context.Context) ([]repo.Session, error)
	UpdateSessionEnd(ctx context.Context) ([]repo.Session, error)
	UpdateSessionStatus(ctx context.Context) ([]repo.Session, error)
	UpdateSessionPrices(ctx context.Context) ([]repo.Session, error)
	ToggleSessionFavorite(ctx context.Context) ([]repo.Session, error)
	AddSessionTag(ctx context.Context) ([]repo.Session, error)
	RemoveSessionTag(ctx context.Context) ([]repo.Session, error)
	DeleteSession(ctx context.Context) ([]repo.Session, error)
	DeleteOldSessions(ctx context.Context) ([]repo.Session, error)
	DeleteSessionsByStatus(ctx context.Context) ([]repo.Session, error)
	CountSessions(ctx context.Context) ([]repo.Session, error)
	CountSessionsByStatus(ctx context.Context) ([]repo.Session, error)
	GetSessionStats(ctx context.Context) ([]repo.Session, error)
	GetSymbolStats(ctx context.Context) ([]repo.Session, error)
	GetStrategyPerformance(ctx context.Context) ([]repo.Session, error)
	GetDailySessionCount(ctx context.Context) ([]repo.Session, error)
	CheckSessionExists(ctx context.Context) ([]repo.Session, error)
	GetLatestSession(ctx context.Context) ([]repo.Session, error)
	GetLongestSession(ctx context.Context) ([]repo.Session, error)
	GetMostProfitableSession(ctx context.Context) ([]repo.Session, error)
	GetMostVolatileSession(ctx context.Context) ([]repo.Session, error)
	GetSessionDuration(ctx context.Context) ([]repo.Session, error)
}
type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateSession(ctx context.Context, params createSessionParams) (repo.Session, error) {
	createdSession, err := s.repo.CreateSession(ctx, repo.CreateSessionParams{
		SessionName:   params.SessionName,
		Parameters:    params.Parameters,
		Symbol:        params.Symbol,
		Timeframe:     params.Timeframe,
		Strategy:      params.Strategy,
		StartingPrice: pgconverter.PGNumericConverter(params.StartingPrice),
		Notes:         pgconverter.PGTextConv(params.Notes),
		Tags:          params.Tags,
	})
	if err != nil {
		return repo.Session{}, dberrors.Handle(err)
	}
	return createdSession, nil

}
func (s *svc) GetSessionByID(ctx context.Context, id int64) (repo.Session, error) {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		return repo.Session{}, dberrors.Handle(err)
	}
	return session, nil

}
func (s *svc) GetSessionByName(ctx context.Context, sessionName string) (repo.Session, error) {
	session, err := s.repo.GetSessionByName(ctx, sessionName)
	if err != nil {
		return repo.Session{}, dberrors.Handle(err)
	}
	return session, nil
}
func (s *svc) ListAllSessions(ctx context.Context) ([]repo.Session, error) {
	sessions, err := s.repo.ListAllSessions(ctx)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil

}
func (s *svc) ListSessionsWithLimit(ctx context.Context, params listSessionParams) ([]repo.Session, error) {
	sessions, err := s.repo.ListSessions(ctx, repo.ListSessionsParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil
}
func (s *svc) GetActiveSessions(ctx context.Context) ([]repo.Session, error) {
	activeSessions, err := s.repo.GetActiveSessions(ctx)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return activeSessions, nil

}
func (s *svc) GetCompletedSessions(ctx context.Context, limit int32) ([]repo.Session, error) {
	completedSessions, err := s.repo.GetCompletedSessions(ctx, limit)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return completedSessions, nil
}
func (s *svc) GetFavoriteSessions(ctx context.Context) ([]repo.Session, error) {
	favoriteSessions, err := s.repo.GetFavoriteSessions(ctx)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return favoriteSessions, nil
}
func (s *svc) GetSessionsBySymbol(ctx context.Context, params getSessionBySymbolsParams) ([]repo.Session, error) {
	sessions, err := s.repo.GetSessionsBySymbol(ctx, repo.GetSessionsBySymbolParams{
		Symbol: params.Symbol,
		Limit:  params.Limit,
	})
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil
}
func (s *svc) GetSessionByStrategy(ctx context.Context, params getSessionByStrategyParams) ([]repo.Session, error) {
	sessions, err := s.repo.GetSessionsByStrategy(ctx, repo.GetSessionsByStrategyParams{
		Strategy: params.Strategy,
		Limit:    params.Limit,
	})
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil

}
func (s *svc) GetSessionByStatus(ctx context.Context, params getSessionByStatusParams) ([]repo.Session, error) {
	sessions, err := s.repo.GetSessionsByStatus(ctx, repo.GetSessionsByStatusParams{
		Status: pgconverter.PGTextConv(params.Status),
		Limit:  params.Limit,
	})
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil
}
func (s *svc) GetRecentSessions(ctx context.Context, timeStamp string) ([]repo.Session, error) {
	newTimeStamp := pgconverter.PGTimeStampConverter(timeStamp)
	sessions, err := s.repo.GetRecentSessions(ctx, newTimeStamp)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil

}
func (s *svc) SearchSessionByName(ctx context.Context, params searchSessionByNameParams) ([]repo.Session, error) {
	sessions, err := s.repo.SearchSessionsByName(ctx, repo.SearchSessionsByNameParams{
		Limit:   params.Limit,
		Column1: pgconverter.PGTextConv(params.Column1),
	})
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil
}
func (s *svc) GetSessionWithTags(ctx context.Context, params getSessionWithTagsParams) ([]repo.Session, error) {
	sessions, err := s.repo.GetSessionsWithTags(ctx, repo.GetSessionsWithTagsParams{
		Limit:   params.Limit,
		Column1: params.Column1,
	})
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil
}
func (s *svc) UpdateSession(ctx context.Context, params updateSessionParams) (repo.Session, error) {
	updatedSession, err := s.repo.UpdateSession(ctx, repo.UpdateSessionParams{
		ID:          params.ID,
		Notes:       pgconverter.PGTextConv(params.Notes),
		SessionName: params.SessionName,
		Tags:        params.Tags,
		IsFavorite:  pgconverter.PGBoolConv(params.IsFavorite),
	})
	if err != nil {
		return repo.Session{}, dberrors.Handle(err)
	}
	return updatedSession, nil

}
func (s *svc) UpdateSessionEnd(ctx context.Context, params updateSessionEndParams) (repo.Session, error) {
	sessionEnd, err := s.repo.UpdateSessionEnd(ctx, repo.UpdateSessionEndParams{
		ID:                 params.ID,
		EndedAt:            pgconverter.PGTimeStampConverter(params.EndedAt),
		Status:             pgconverter.PGTextConv(params.Status),
		EndingPrice:        pgconverter.PGNumericConverter(params.EndingPrice),
		HighestPrice:       pgconverter.PGNumericConverter(params.HighestPrice),
		LowestPrice:        pgconverter.PGNumericConverter(params.LowestPrice),
		TotalCandles:       pgconverter.PGInt4Converter(params.TotalCandles),
		PriceChangePercent: pgconverter.PGNumericConverter(params.PriceChangePercent),
		Volatility:         pgconverter.PGNumericConverter(params.Volatility),
	})
	if err != nil {
		return repo.Session{}, dberrors.Handle(err)
	}
	return sessionEnd, nil
}

func (s *svc) UpdateSessionStatus(ctx context.Context, params updateSessionStatus) error {
	err := s.repo.UpdateSessionStatus(ctx, repo.UpdateSessionStatusParams{
		ID:     params.ID,
		Status: pgconverter.PGTextConv(params.Status),
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}
func (s *svc) UpdateSessionPrices(ctx context.Context, params updateSessionPricesParams) error {
	err := s.repo.UpdateSessionPrices(ctx, repo.UpdateSessionPricesParams{
		ID:           params.ID,
		HighestPrice: pgconverter.PGNumericConverter(params.HighestPrice),
		LowestPrice:  pgconverter.PGNumericConverter(params.LowestPrice),
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) ToggleSessionFavorite(ctx context.Context, id int64) (bool, error) {
	boolean, err := s.repo.ToggleSessionFavorite(ctx, id)
	if err != nil {
		return false, dberrors.Handle(err)
	}
	return pgconverter.PGToBool(boolean), nil
}

func (s *svc) AddSessionTag(ctx context.Context, params addSessionTagParams) error {
	err := s.repo.AddSessionTag(ctx, repo.AddSessionTagParams{
		ID:          params.ID,
		ArrayAppend: params.ArrayAppend,
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil

}
func (s *svc) RemoveSessionTag(ctx context.Context, params removeSessionTagParams) error {
	err := s.repo.RemoveSessionTag(ctx, repo.RemoveSessionTagParams{
		ID:          params.ID,
		ArrayRemove: params.ArrayRemove,
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) DeleteSession(ctx context.Context, id int64) error {
	err := s.repo.DeleteSession(ctx, id)
	if err != nil {
		dberrors.Handle(err)
	}
	return nil
}

func (s *svc) DeleteOldSessions(ctx context.Context, endedAt string) error {
	str := pgconverter.PGTimeStampConverter(endedAt)
	err := s.repo.DeleteOldSessions(ctx, str)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}
func (s *svc) DeleteSessionsByStatus(ctx context.Context, stringStatus string) error {
	status := pgconverter.PGTextConv(stringStatus)
	err := s.repo.DeleteSessionsByStatus(ctx, status)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) CountSessions(ctx context.Context) (int64, error) {
	num, err := s.repo.CountSessions(ctx)
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return num, nil
}

func (s *svc) CountSessionsByStatus(ctx context.Context, stringStatus string) (int64, error) {
	status := pgconverter.PGTextConv(stringStatus)
	num, err := s.repo.CountSessionsByStatus(ctx, status)
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return num, nil

}
func (s *svc) GetSessionStats(ctx context.Context) (repo.GetSessionStatsRow, error) {
	sessionRow, err := s.repo.GetSessionStats(ctx)
	if err != nil {
		return repo.GetSessionStatsRow{}, dberrors.Handle(err)
	}
	return sessionRow, nil
}

func (s *svc) GetSymbolStats(ctx context.Context) ([]repo.Session, error) {
	sessionStats, err := s.GetSymbolStats(ctx)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessionStats, nil
}
func (s *svc) GetStrategyPerformance(ctx context.Context) ([]repo.GetStrategyPerformanceRow, error) {
	strategyPerformance, err := s.repo.GetStrategyPerformance(ctx)
	if err != nil {
		return []repo.GetStrategyPerformanceRow{}, dberrors.Handle(err)
	}
	return strategyPerformance, nil

}
func (s *svc) GetDailySessionCount(ctx context.Context, startedAt string) ([]repo.GetDailySessionCountRow, error) {
	timeStamp := pgconverter.PGTimeStampConverter(startedAt)
	dailySessionCountRows, err := s.repo.GetDailySessionCount(ctx, timeStamp)
	if err != nil {
		return []repo.GetDailySessionCountRow{}, dberrors.Handle(err)
	}
	return dailySessionCountRows, nil
}
func (s *svc) CheckSessionExists(ctx context.Context, id int64) (bool, error) {
	boolean, err := s.repo.CheckSessionExists(ctx, id)
	if err != nil {
		return false, dberrors.Handle(err)
	}
	return boolean, nil
}

func (s *svc) GetLatestSession(ctx context.Context) (repo.Session, error) {
	session, err := s.repo.GetLatestSession(ctx)
	if err != nil {
		return repo.Session{}, dberrors.Handle(err)
	}
	return session, nil
}
func (s *svc) GetLongestSession(ctx context.Context, limit int32) ([]repo.Session, error) {
	sessions, err := s.repo.GetLongestSessions(ctx, limit)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return sessions, nil
}

func (s *svc) GetMostProfitableSession(ctx context.Context, limit int32) ([]repo.Session, error) {
	mostProfitableSessions, err := s.repo.GetMostProfitableSessions(ctx, limit)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return mostProfitableSessions, nil
}

func (s *svc) GetMostVolatileSession(ctx context.Context, limit int32) ([]repo.Session, error) {
	mostVolatileSessions, err := s.repo.GetMostVolatileSessions(ctx, limit)
	if err != nil {
		return []repo.Session{}, dberrors.Handle(err)
	}
	return mostVolatileSessions, nil
}
func (s *svc) GetSessionDuration(ctx context.Context, id int64) (int32, error) {
	duration, err := s.repo.GetSessionDuration(ctx, id)
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return pgconverter.Int4ToInt32(duration), nil
}
