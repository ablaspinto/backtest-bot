package userconfigs

import (
	repo "cmd/internal/adapters/postgresql/sqlc"
	"cmd/internal/dberrors"
	"cmd/internal/pgconverter"
	"context"
)

type Service interface {
	CreateConfig(ctx context.Context, params createConfigParams) (repo.Config, error)
	GetConfigByID(ctx context.Context, id int64) (repo.Config, error)
	GetConfigByName(ctx context.Context, name string) (repo.Config, error)
	GetDefaultConfig(ctx context.Context) (repo.Config, error)
	ListConfigs(ctx context.Context) ([]repo.Config, error)
	ListConfigsPaginated(ctx context.Context, params listConfigsPaginatedParams) ([]repo.Config, error)
	GetFavoriteConfigs(ctx context.Context) ([]repo.Config, error)
	GetPublicConfigs(ctx context.Context, limit int32) ([]repo.Config, error)
	GetConfigsByStrategy(ctx context.Context, strategy string) ([]repo.Config, error)
	GetConfigsByMarketStrategy(ctx context.Context, strat string) ([]repo.Config, error)
	GetConfigsByCreator(ctx context.Context, createdAt string) ([]repo.Config, error)
	SearchConfigsByName(ctx context.Context, params searchConfigsByNameParams) ([]repo.Config, error)
	GetConfigsWithTags(ctx context.Context, params getConfigWithTagsParams) ([]repo.Config, error)
	GetRecentlyUsedConfigs(ctx context.Context, limit int32) ([]repo.Config, error)
	GetPopularConfigs(ctx context.Context, limit int32) ([]repo.Config, error)
	UpdateConfig(ctx context.Context, params updateConfigParams) (repo.Config, error)
	UpdateConfigPartial(ctx context.Context, params updateConfigPartialParams) (repo.Config, error)
	IncrementConfigUsage(ctx context.Context, id int64) error
	SetDefaultConfig(ctx context.Context, id int64) error
	ToggleFavorite(ctx context.Context, id int64) (bool, error)
	ToggleConfigPublic(ctx context.Context, id int64) (bool, error)
	AddConfigTag(ctx context.Context, params addConfigTagParams) error
	RemoveConfigTag(ctx context.Context, params removeConfigTagParams) error
	CountConfigs(ctx context.Context) (int64, error)
	CountConfigsByStrategy(ctx context.Context, strat string) (int64, error)
	GetConfigStats(ctx context.Context) (repo.GetConfigStatsRow, error)
	GetStrategyDistribution(ctx context.Context) ([]repo.GetStrategyDistributionRow, error)
	GetMarketBehaviorDistribution(ctx context.Context) ([]repo.GetMarketBehaviorDistributionRow, error)
	CheckConfigExists(ctx context.Context, id int64) (bool, error)
	CheckConfigNameExists(ctx context.Context, name string) (bool, error)
	GetConfigNameOrDefault(ctx context.Context, name string) (repo.Config, error)
	DuplicateConfig(ctx context.Context, params duplicateConfigParams) (repo.Config, error)
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateConfig(ctx context.Context, params createConfigParams) (repo.Config, error) {
	config, err := s.repo.CreateConfig(ctx, repo.CreateConfigParams{
		Description:        pgconverter.PGTextConv(params.Description),
		Drift:              pgconverter.PGNumericConverter(params.Drift),
		Name:               params.Name,
		Strategy:           params.Strategy,
		Volatility:         pgconverter.PGNumericConverter(params.Volatility),
		StartingPrice:      pgconverter.PGNumericConverter(params.StartingPrice),
		CandleInterval:     params.CandleInterval,
		Symbol:             params.Symbol,
		MarketBehavior:     pgconverter.PGTextConv(params.MarketBehavior),
		TrendStrength:      pgconverter.PGNumericConverter(params.TrendStrength),
		MeanReversionSpeed: pgconverter.PGNumericConverter(params.MeanReversionSpeed),
		AdvancedParams:     params.AdvancedParams,
		Tags:               params.Tags,
	})
	if err != nil {
		return repo.Config{}, dberrors.Handle(err)
	}
	return config, nil
}
func (s *svc) GetConfigByID(ctx context.Context, id int64) (repo.Config, error) {
	config, err := s.repo.GetConfigByID(ctx, id)
	if err != nil {
		return repo.Config{}, dberrors.Handle(err)
	}
	return config, nil
}

func (s *svc) GetConfigByName(ctx context.Context, name string) (repo.Config, error) {
	config, err := s.repo.GetConfigByName(ctx, name)
	if err != nil {
		return repo.Config{}, dberrors.Handle(err)
	}
	return config, nil
}
func (s *svc) GetDefaultConfig(ctx context.Context) (repo.Config, error) {
	config, err := s.repo.GetDefaultConfig(ctx)
	if err != nil {
		return repo.Config{}, dberrors.Handle(err)
	}
	return config, nil
}

func (s *svc) ListConfigs(ctx context.Context) ([]repo.Config, error) {
	configs, err := s.repo.ListConfigs(ctx)
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}

func (s *svc) ListConfigsPaginated(ctx context.Context, params listConfigsPaginatedParams) ([]repo.Config, error) {
	configs, err := s.repo.ListConfigsPaginated(ctx, repo.ListConfigsPaginatedParams{
		Offset: params.Offset,
		Limit:  params.Limit,
	})
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}

func (s *svc) GetFavoriteConfigs(ctx context.Context) ([]repo.Config, error) {
	configs, err := s.repo.GetFavoriteConfigs(ctx)
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}

func (s *svc) GetPublicConfigs(ctx context.Context, limit int32) ([]repo.Config, error) {
	configs, err := s.repo.GetPublicConfigs(ctx, limit)
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}

func (s *svc) GetConfigsByStrategy(ctx context.Context, strat string) ([]repo.Config, error) {
	configs, err := s.repo.GetConfigsByStrategy(ctx, strat)
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}
func (s *svc) GetConfigsByMarketStrategy(ctx context.Context, strat string) ([]repo.Config, error) {
	configs, err := s.repo.GetConfigsByStrategy(ctx, strat)
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}

func (s *svc) GetConfigsByCreator(ctx context.Context, createdAt string) ([]repo.Config, error) {
	timeStamp := pgconverter.PGTextConv(createdAt)
	configs, err := s.repo.GetConfigsByCreator(ctx, timeStamp)
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}
func (s *svc) SearchConfigsByName(ctx context.Context, params searchConfigsByNameParams) ([]repo.Config, error) {
	configs, err := s.repo.SearchConfigsByName(ctx, repo.SearchConfigsByNameParams{
		Limit:   params.Limit,
		Column1: pgconverter.PGTextConv(params.Column1),
	})
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}

func (s *svc) GetConfigsWithTags(ctx context.Context, params getConfigWithTagsParams) ([]repo.Config, error) {

	configs, err := s.repo.GetConfigsWithTags(ctx, repo.GetConfigsWithTagsParams{
		Limit:   params.Limit,
		Column1: params.Column1,
	})
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}
func (s *svc) GetRecentlyUsedConfigs(ctx context.Context, limit int32) ([]repo.Config, error) {
	configs, err := s.repo.GetRecentlyUsedConfigs(ctx, limit)
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)

	}
	return configs, nil
}

func (s *svc) GetPopularConfigs(ctx context.Context, limit int32) ([]repo.Config, error) {
	configs, err := s.repo.GetPopularConfigs(ctx, limit)
	if err != nil {
		return []repo.Config{}, dberrors.Handle(err)
	}
	return configs, nil
}

func (s *svc) UpdateConfig(ctx context.Context, params updateConfigParams) (repo.Config, error) {

	updatedConfig, err := s.repo.UpdateConfig(ctx, repo.UpdateConfigParams{
		Name:               params.Name,
		Symbol:             params.Symbol,
		Strategy:           params.Strategy,
		StartingPrice:      pgconverter.PGNumericConverter(params.StartingPrice),
		AdvancedParams:     params.AdvancedParams,
		Description:        pgconverter.PGTextConv(params.Description),
		Drift:              pgconverter.PGNumericConverter(params.Drift),
		MarketBehavior:     pgconverter.PGTextConv(params.MarketBehavior),
		MeanReversionSpeed: pgconverter.PGNumericConverter(params.MeanReversionSpeed),
		TrendStrength:      pgconverter.PGNumericConverter(params.TrendStrength),
		Tags:               params.Tags,
		Volatility:         pgconverter.PGNumericConverter(params.Volatility),
		ID:                 params.ID,
		CandleInterval:     params.CandleInterval,
	})
	if err != nil {
		return repo.Config{}, dberrors.Handle(err)
	}
	return updatedConfig, nil
}

func (s *svc) UpdateConfigPartial(ctx context.Context, params updateConfigPartialParams) (repo.Config, error) {
	updatedConfigPartial, err := s.repo.UpdateConfigPartial(ctx, repo.UpdateConfigPartialParams{
		ID:          params.ID,
		Drift:       pgconverter.PGNumericConverter(params.Drift),
		Description: pgconverter.PGTextConv(params.Description),
		Volatility:  pgconverter.PGNumericConverter(params.Volatility),
		Name:        pgconverter.PGTextConv(params.Name),
	})
	if err != nil {
		return repo.Config{}, dberrors.Handle(err)
	}
	return updatedConfigPartial, nil
}
func (s *svc) IncrementConfigUsage(ctx context.Context, id int64) error {
	err := s.repo.IncrementConfigUsage(ctx, id)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) SetDefaultConfig(ctx context.Context, id int64) error {
	err := s.repo.SetDefaultConfig(ctx, id)
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) ToggleFavorite(ctx context.Context, id int64) (bool, error) {
	boolean, err := s.repo.ToggleConfigFavorite(ctx, id)
	if err != nil {
		return false, dberrors.Handle(err)
	}
	return pgconverter.PGToBool(boolean), nil
}

func (s *svc) ToggleConfigPublic(ctx context.Context, id int64) (bool, error) {
	boolean, err := s.repo.ToggleConfigPublic(ctx, id)
	if err != nil {
		return false, dberrors.Handle(err)
	}
	return pgconverter.PGToBool(boolean), nil
}

func (s *svc) AddConfigTag(ctx context.Context, params addConfigTagParams) error {
	err := s.repo.AddConfigTag(ctx, repo.AddConfigTagParams{
		ID:          params.ID,
		ArrayAppend: params.ArrayAppend,
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) RemoveConfigTag(ctx context.Context, params removeConfigTagParams) error {
	err := s.repo.RemoveConfigTag(ctx, repo.RemoveConfigTagParams{
		ID:          params.ID,
		ArrayRemove: params.ArrayRemove,
	})
	if err != nil {
		return dberrors.Handle(err)
	}
	return nil
}

func (s *svc) CountConfigs(ctx context.Context) (int64, error) {
	num, err := s.repo.CountConfigs(ctx)
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return num, nil
}

func (s *svc) CountConfigsByStrategy(ctx context.Context, strat string) (int64, error) {
	num, err := s.repo.CountConfigsByStrategy(ctx, strat)
	if err != nil {
		return 0, dberrors.Handle(err)
	}
	return num, nil
}

func (s *svc) GetConfigStats(ctx context.Context) (repo.GetConfigStatsRow, error) {
	configStatsRow, err := s.repo.GetConfigStats(ctx)

	if err != nil {
		return repo.GetConfigStatsRow{}, dberrors.Handle(err)
	}
	return configStatsRow, nil
}

func (s *svc) GetStrategyDistribution(ctx context.Context) ([]repo.GetStrategyDistributionRow, error) {
	rows, err := s.repo.GetStrategyDistribution(ctx)
	if err != nil {
		return []repo.GetStrategyDistributionRow{}, dberrors.Handle(err)
	}
	return rows, nil
}

func (s *svc) GetMarketBehaviorDistribution(ctx context.Context) ([]repo.GetMarketBehaviorDistributionRow, error) {
	behavior, err := s.repo.GetMarketBehaviorDistribution(ctx)
	if err != nil {
		return []repo.GetMarketBehaviorDistributionRow{}, dberrors.Handle(err)
	}
	return behavior, nil
}

func (s *svc) CheckConfigExists(ctx context.Context, id int64) (bool, error) {
	boolean, err := s.repo.CheckConfigExists(ctx, id)
	if err != nil {
		return false, dberrors.Handle(err)
	}
	return boolean, nil
}

func (s *svc) CheckConfigNameExists(ctx context.Context, name string) (bool, error) {
	boolean, err := s.repo.CheckConfigNameExists(ctx, name)
	if err != nil {
		return false, dberrors.Handle(err)
	}
	return boolean, nil
}
func (s *svc) GetConfigNameOrDefault(ctx context.Context, name string) (repo.Config, error) {
	config, err := s.repo.GetConfigByNameOrDefault(ctx, name)
	if err != nil {
		return repo.Config{}, dberrors.Handle(err)
	}
	return config, nil
}

func (s svc) DuplicateConfig(ctx context.Context, params duplicateConfigParams) (repo.Config, error) {
	duplicatedConfig, err := s.repo.DuplicateConfig(ctx, repo.DuplicateConfigParams{
		ID:   params.ID,
		Name: params.Name,
	})
	if err != nil {
		return repo.Config{}, dberrors.Handle(err)
	}
	return duplicatedConfig, nil
}
