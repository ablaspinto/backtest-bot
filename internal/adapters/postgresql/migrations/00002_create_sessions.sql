-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    session_name VARCHAR(255) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    timeframe VARCHAR(10) NOT NULL,
    strategy VARCHAR(50) NOT NULL DEFAULT 'random_walk',
    
    -- Session timing
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER,
    
    -- Price data
    starting_price NUMERIC(20, 8) NOT NULL,
    ending_price NUMERIC(20, 8),
    highest_price NUMERIC(20, 8),
    lowest_price NUMERIC(20, 8),
    
    -- Statistics
    total_candles INTEGER DEFAULT 0,
    price_change_percent NUMERIC(10, 4),
    volatility NUMERIC(10, 6),
    
    -- Configuration (stored as JSONB for flexibility)
    parameters JSONB DEFAULT '{}'::jsonb,
    
    -- Session status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed', 'failed')),
    
    -- Metadata
    notes TEXT,
    tags VARCHAR(50)[],
    is_favorite BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_ended_after_started CHECK (ended_at IS NULL OR ended_at >= started_at),
    CONSTRAINT chk_prices_valid CHECK (
        (ending_price IS NULL OR ending_price > 0) AND
        (highest_price IS NULL OR highest_price >= starting_price) AND
        (lowest_price IS NULL OR lowest_price <= starting_price)
    )
);

-- Index for finding recent sessions
CREATE INDEX idx_sessions_started_at 
    ON sessions(started_at DESC);

-- Index for status queries
CREATE INDEX idx_sessions_status 
    ON sessions(status) WHERE status = 'active';

-- Index for symbol lookups
CREATE INDEX idx_sessions_symbol 
    ON sessions(symbol);

-- Index for finding favorites
CREATE INDEX idx_sessions_favorites 
    ON sessions(is_favorite) WHERE is_favorite = true;

-- GIN index for JSONB parameters (enables fast queries on config)
CREATE INDEX idx_sessions_parameters 
    ON sessions USING GIN (parameters);

-- Index for tag searches
CREATE INDEX idx_sessions_tags 
    ON sessions USING GIN (tags);

-- Add table comments
COMMENT ON TABLE sessions IS 'Tracks individual simulation sessions and their outcomes';
COMMENT ON COLUMN sessions.strategy IS 'Simulation strategy used (random_walk, trending, volatile, mean_reverting)';
COMMENT ON COLUMN sessions.parameters IS 'JSON object containing strategy-specific parameters (volatility, drift, etc.)';
COMMENT ON COLUMN sessions.duration_seconds IS 'Total duration of the session in seconds';
COMMENT ON COLUMN sessions.tags IS 'Array of tags for categorizing sessions';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions CASCADE;
-- +goose StatementEnd
