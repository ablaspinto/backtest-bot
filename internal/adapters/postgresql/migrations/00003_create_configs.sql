-- +goose Up
-- +goose StatementBegin
CREATE TABLE configs (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Market behavior settings
    strategy VARCHAR(50) NOT NULL DEFAULT 'random_walk',
    volatility NUMERIC(10, 6) NOT NULL DEFAULT 0.02,
    drift NUMERIC(10, 6) DEFAULT 0.0,
    
    -- Simulation settings
    starting_price NUMERIC(20, 8) NOT NULL DEFAULT 100.0,
    candle_interval VARCHAR(10) NOT NULL DEFAULT '1m',
    symbol VARCHAR(20) NOT NULL DEFAULT 'SIM',
    
    -- Market conditions
    market_behavior VARCHAR(30) DEFAULT 'neutral' 
        CHECK (market_behavior IN ('bullish', 'bearish', 'neutral', 'volatile', 'sideways')),
    trend_strength NUMERIC(5, 4) DEFAULT 0.5 
        CHECK (trend_strength >= 0 AND trend_strength <= 1),
    mean_reversion_speed NUMERIC(5, 4) DEFAULT 0.1,
    
    -- Advanced parameters (stored as JSONB for flexibility)
    advanced_params JSONB DEFAULT '{}'::jsonb,
    
    -- User preferences
    is_default BOOLEAN DEFAULT false,
    is_favorite BOOLEAN DEFAULT false,
    is_public BOOLEAN DEFAULT false,
    
    -- Usage tracking
    times_used INTEGER DEFAULT 0,
    last_used_at TIMESTAMP WITH TIME ZONE,
    
    -- Metadata
    created_by VARCHAR(100),
    tags VARCHAR(50)[],
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_volatility_positive CHECK (volatility >= 0),
    CONSTRAINT chk_starting_price_positive CHECK (starting_price > 0)
);

-- Index for finding default config
CREATE UNIQUE INDEX idx_configs_default 
    ON configs(is_default) WHERE is_default = true;

-- Index for favorites
CREATE INDEX idx_configs_favorites 
    ON configs(is_favorite) WHERE is_favorite = true;

-- Index for strategy lookups
CREATE INDEX idx_configs_strategy 
    ON configs(strategy);

-- Index for market behavior
CREATE INDEX idx_configs_market_behavior 
    ON configs(market_behavior);

-- Index for most used configs
CREATE INDEX idx_configs_usage 
    ON configs(times_used DESC);

-- GIN index for JSONB advanced params
CREATE INDEX idx_configs_advanced_params 
    ON configs USING GIN (advanced_params);

-- Index for tags
CREATE INDEX idx_configs_tags 
    ON configs USING GIN (tags);

-- Add table comments
COMMENT ON TABLE configs IS 'Stores simulation configuration presets and user preferences';
COMMENT ON COLUMN configs.volatility IS 'Price volatility coefficient (0.02 = 2% standard deviation)';
COMMENT ON COLUMN configs.drift IS 'Price drift/trend direction (-1 to 1, negative=bearish, positive=bullish)';
COMMENT ON COLUMN configs.mean_reversion_speed IS 'How quickly prices revert to mean (0-1, higher = faster)';
COMMENT ON COLUMN configs.trend_strength IS 'Strength of trending behavior (0-1)';
COMMENT ON COLUMN configs.advanced_params IS 'Additional parameters in JSON format for extensibility';
COMMENT ON COLUMN configs.is_default IS 'Whether this is the default configuration (only one can be true)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS configs CASCADE;
-- +goose StatementEnd
