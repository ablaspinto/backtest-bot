-- +goose Up
-- +goose StatementBegin
CREATE TABLE indicators (
    id BIGSERIAL PRIMARY KEY,
    candle_id BIGINT NOT NULL,
    indicator_type VARCHAR(50) NOT NULL,
    
    -- Core indicator value
    value NUMERIC(20, 8),
    
    -- Multi-value indicators (e.g., Bollinger Bands, Ichimoku)
    value_upper NUMERIC(20, 8),
    value_lower NUMERIC(20, 8),
    value_middle NUMERIC(20, 8),
    
    -- Additional components
    signal_line NUMERIC(20, 8),
    histogram NUMERIC(20, 8),
    
    -- Indicator parameters used
    period INTEGER,
    period_fast INTEGER,
    period_slow INTEGER,
    period_signal INTEGER,
    standard_deviations NUMERIC(5, 2),
    
    -- Signal interpretation
    signal VARCHAR(10) CHECK (signal IN ('buy', 'sell', 'hold', 'neutral', 'overbought', 'oversold')),
    signal_strength NUMERIC(5, 4) CHECK (signal_strength >= 0 AND signal_strength <= 1),
    
    -- Crossover events
    is_crossover BOOLEAN DEFAULT false,
    crossover_type VARCHAR(20) CHECK (
        crossover_type IS NULL OR 
        crossover_type IN ('golden', 'death', 'bullish', 'bearish', 'signal')
    ),
    
    -- Additional data as JSONB
    metadata JSONB DEFAULT '{}'::jsonb,
    
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key
    CONSTRAINT fk_indicators_candle 
        FOREIGN KEY (candle_id) 
        REFERENCES candles(id) 
        ON DELETE CASCADE,
    
    -- Ensure we don't duplicate indicator calculations
    CONSTRAINT uq_indicators_candle_type_params 
        UNIQUE (candle_id, indicator_type, period, period_fast, period_slow)
);

-- Primary lookup index
CREATE INDEX idx_indicators_candle 
    ON indicators(candle_id);

-- Index for indicator type queries
CREATE INDEX idx_indicators_type 
    ON indicators(indicator_type);

-- Composite index for common queries
CREATE INDEX idx_indicators_candle_type 
    ON indicators(candle_id, indicator_type);

-- Index for signal queries
CREATE INDEX idx_indicators_signal 
    ON indicators(signal) WHERE signal IN ('buy', 'sell');

-- Index for crossover events
CREATE INDEX idx_indicators_crossover 
    ON indicators(is_crossover, crossover_type) 
    WHERE is_crossover = true;

-- GIN index for metadata
CREATE INDEX idx_indicators_metadata 
    ON indicators USING GIN (metadata);

-- Add table comments
COMMENT ON TABLE indicators IS 'Stores calculated technical indicators for each candle';
COMMENT ON COLUMN indicators.indicator_type IS 'Type of indicator (SMA, EMA, RSI, MACD, BB, STOCH, ADX, etc.)';
COMMENT ON COLUMN indicators.value IS 'Primary indicator value';
COMMENT ON COLUMN indicators.value_upper IS 'Upper band/line (e.g., Bollinger upper band)';
COMMENT ON COLUMN indicators.value_lower IS 'Lower band/line (e.g., Bollinger lower band)';
COMMENT ON COLUMN indicators.value_middle IS 'Middle line (e.g., Bollinger middle band/SMA)';
COMMENT ON COLUMN indicators.signal_line IS 'Signal line for indicators like MACD';
COMMENT ON COLUMN indicators.histogram IS 'Histogram value (e.g., MACD histogram)';
COMMENT ON COLUMN indicators.period IS 'Period used for calculation (e.g., 14 for RSI-14)';
COMMENT ON COLUMN indicators.period_fast IS 'Fast period for dual-period indicators (e.g., MACD fast)';
COMMENT ON COLUMN indicators.period_slow IS 'Slow period for dual-period indicators (e.g., MACD slow)';
COMMENT ON COLUMN indicators.standard_deviations IS 'Standard deviations for Bollinger Bands';
COMMENT ON COLUMN indicators.crossover_type IS 'Type of crossover event detected';
COMMENT ON COLUMN indicators.metadata IS 'Additional indicator-specific data in JSON format';

-- Create a view for commonly used indicators
CREATE OR REPLACE VIEW v_candles_with_indicators AS
SELECT 
    c.id as candle_id,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.open,
    c.high,
    c.low,
    c.close,
    c.volume,
    MAX(CASE WHEN i.indicator_type = 'SMA' AND i.period = 20 THEN i.value END) as sma_20,
    MAX(CASE WHEN i.indicator_type = 'SMA' AND i.period = 50 THEN i.value END) as sma_50,
    MAX(CASE WHEN i.indicator_type = 'EMA' AND i.period = 12 THEN i.value END) as ema_12,
    MAX(CASE WHEN i.indicator_type = 'EMA' AND i.period = 26 THEN i.value END) as ema_26,
    MAX(CASE WHEN i.indicator_type = 'RSI' AND i.period = 14 THEN i.value END) as rsi_14,
    MAX(CASE WHEN i.indicator_type = 'MACD' THEN i.value END) as macd,
    MAX(CASE WHEN i.indicator_type = 'MACD' THEN i.signal_line END) as macd_signal,
    MAX(CASE WHEN i.indicator_type = 'MACD' THEN i.histogram END) as macd_histogram,
    MAX(CASE WHEN i.indicator_type = 'BB' THEN i.value_upper END) as bb_upper,
    MAX(CASE WHEN i.indicator_type = 'BB' THEN i.value_middle END) as bb_middle,
    MAX(CASE WHEN i.indicator_type = 'BB' THEN i.value_lower END) as bb_lower
FROM candles c
LEFT JOIN indicators i ON c.id = i.candle_id
GROUP BY c.id, c.symbol, c.timeframe, c.timestamp, c.open, c.high, c.low, c.close, c.volume;

COMMENT ON VIEW v_candles_with_indicators IS 'Denormalized view of candles with common indicators for easy querying';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS v_candles_with_indicators;
DROP TABLE IF EXISTS indicators CASCADE;
-- +goose StatementEnd
