-- +goose Up
-- +goose StatementBegin
CREATE TABLE candles (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    timeframe VARCHAR(10) NOT NULL,
    timestamp BIGINT NOT NULL,
    open NUMERIC(20, 8) NOT NULL,
    high NUMERIC(20, 8) NOT NULL,
    low NUMERIC(20, 8) NOT NULL,
    close NUMERIC(20, 8) NOT NULL,
    volume NUMERIC(20, 8) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Constraints
    CONSTRAINT chk_ohlc_valid CHECK (
        high >= open AND 
        high >= close AND 
        high >= low AND
        low <= open AND 
        low <= close
    ),
    CONSTRAINT chk_volume_positive CHECK (volume >= 0)
);

-- Unique constraint to prevent duplicate candles
CREATE UNIQUE INDEX idx_candles_symbol_timeframe_timestamp 
    ON candles(symbol, timeframe, timestamp);

-- Index for time-based queries (most common query pattern)
CREATE INDEX idx_candles_timestamp 
    ON candles(timestamp DESC);

-- Index for symbol + timeframe queries
CREATE INDEX idx_candles_symbol_timeframe 
    ON candles(symbol, timeframe);

-- Composite index for the most common query pattern
CREATE INDEX idx_candles_lookup 
    ON candles(symbol, timeframe, timestamp DESC);

-- Add table comment
COMMENT ON TABLE candles IS 'Stores OHLC candlestick data for all symbols and timeframes';
COMMENT ON COLUMN candles.symbol IS 'Trading symbol (e.g., BTC, AAPL, EURUSD)';
COMMENT ON COLUMN candles.timeframe IS 'Candle interval (e.g., 1m, 5m, 15m, 1h, 1d)';
COMMENT ON COLUMN candles.timestamp IS 'Unix timestamp in seconds representing candle open time';
COMMENT ON COLUMN candles.open IS 'Opening price for the time period';
COMMENT ON COLUMN candles.high IS 'Highest price during the time period';
COMMENT ON COLUMN candles.low IS 'Lowest price during the time period';
COMMENT ON COLUMN candles.close IS 'Closing price for the time period';
COMMENT ON COLUMN candles.volume IS 'Total volume traded during the time period';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS candles CASCADE;
-- +goose StatementEnd
