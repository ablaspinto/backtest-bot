-- +goose Up
-- +goose StatementBegin
CREATE TABLE metrics (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL,
    timestamp BIGINT NOT NULL,
    
    -- Price metrics
    current_price NUMERIC(20, 8) NOT NULL,
    price_change NUMERIC(20, 8),
    price_change_percent NUMERIC(10, 4),
    
    -- Session statistics
    session_high NUMERIC(20, 8),
    session_low NUMERIC(20, 8),
    session_open NUMERIC(20, 8),
    
    -- Volatility metrics
    realized_volatility NUMERIC(10, 6),
    volatility_1min NUMERIC(10, 6),
    volatility_5min NUMERIC(10, 6),
    volatility_1hour NUMERIC(10, 6),
    
    -- Volume metrics
    volume NUMERIC(20, 8) DEFAULT 0,
    volume_ma_5 NUMERIC(20, 8),
    volume_ma_20 NUMERIC(20, 8),
    
    -- Performance metrics
    drawdown_percent NUMERIC(10, 4),
    max_drawdown_percent NUMERIC(10, 4),
    return_percent NUMERIC(10, 4),
    
    -- Moving averages
    sma_10 NUMERIC(20, 8),
    sma_20 NUMERIC(20, 8),
    sma_50 NUMERIC(20, 8),
    sma_200 NUMERIC(20, 8),
    ema_10 NUMERIC(20, 8),
    ema_20 NUMERIC(20, 8),
    
    -- Momentum indicators
    rsi_14 NUMERIC(5, 2),
    macd NUMERIC(20, 8),
    macd_signal NUMERIC(20, 8),
    macd_histogram NUMERIC(20, 8),
    
    -- Statistical measures
    standard_deviation NUMERIC(20, 8),
    variance NUMERIC(20, 8),
    skewness NUMERIC(10, 6),
    kurtosis NUMERIC(10, 6),
    
    -- Trade signals (if applicable)
    signal VARCHAR(10) CHECK (signal IN ('buy', 'sell', 'hold', 'neutral')),
    signal_strength NUMERIC(5, 4) CHECK (signal_strength >= 0 AND signal_strength <= 1),
    
    -- Additional metrics as JSONB for flexibility
    custom_metrics JSONB DEFAULT '{}'::jsonb,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key
    CONSTRAINT fk_metrics_session 
        FOREIGN KEY (session_id) 
        REFERENCES sessions(id) 
        ON DELETE CASCADE,
    
    -- Constraints
    CONSTRAINT chk_price_positive CHECK (current_price > 0),
    CONSTRAINT chk_rsi_range CHECK (rsi_14 IS NULL OR (rsi_14 >= 0 AND rsi_14 <= 100))
);

-- Primary lookup index (session + time)
CREATE INDEX idx_metrics_session_timestamp 
    ON metrics(session_id, timestamp DESC);

-- Index for time-based queries
CREATE INDEX idx_metrics_timestamp 
    ON metrics(timestamp DESC);

-- Index for signal queries
CREATE INDEX idx_metrics_signal 
    ON metrics(signal) WHERE signal IS NOT NULL;

-- Partial index for buy/sell signals
CREATE INDEX idx_metrics_trade_signals 
    ON metrics(session_id, timestamp, signal) 
    WHERE signal IN ('buy', 'sell');

-- GIN index for custom metrics
CREATE INDEX idx_metrics_custom 
    ON metrics USING GIN (custom_metrics);

-- Add table comments
COMMENT ON TABLE metrics IS 'Stores calculated metrics and analytics for each point in a session';
COMMENT ON COLUMN metrics.realized_volatility IS 'Actual measured volatility over the period';
COMMENT ON COLUMN metrics.drawdown_percent IS 'Current drawdown from session high';
COMMENT ON COLUMN metrics.max_drawdown_percent IS 'Maximum drawdown experienced in session so far';
COMMENT ON COLUMN metrics.rsi_14 IS 'Relative Strength Index (14 period)';
COMMENT ON COLUMN metrics.macd IS 'Moving Average Convergence Divergence indicator';
COMMENT ON COLUMN metrics.signal_strength IS 'Confidence level of the signal (0-1)';
COMMENT ON COLUMN metrics.custom_metrics IS 'Additional user-defined or strategy-specific metrics';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics CASCADE;
-- +goose StatementEnd
