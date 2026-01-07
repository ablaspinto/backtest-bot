-- Candles Queries
-- CRUD operations and queries for OHLC candlestick data

-- =============================================================================
-- CREATE Operations
-- =============================================================================

-- name: CreateCandle :one
INSERT INTO candles (
    symbol,
    timeframe,
    timestamp,
    open,
    high,
    low,
    close,
    volume
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (symbol, timeframe, timestamp) DO NOTHING
RETURNING *;

-- name: CreateCandleBatch :copyfrom
INSERT INTO candles (
    symbol,
    timeframe,
    timestamp,
    open,
    high,
    low,
    close,
    volume
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- =============================================================================
-- READ Operations
-- =============================================================================

-- name: GetCandleByID :one
SELECT * FROM candles
WHERE id = $1;

-- name: GetCandle :one
SELECT * FROM candles
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp = $3
LIMIT 1;

-- name: GetLatestCandle :one
SELECT * FROM candles
WHERE symbol = $1
  AND timeframe = $2
ORDER BY timestamp DESC
LIMIT 1;

-- name: GetRecentCandles :many
SELECT * FROM candles
WHERE symbol = $1
  AND timeframe = $2
ORDER BY timestamp DESC
LIMIT $3;

-- name: GetCandlesInRange :many
SELECT * FROM candles
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp >= $3
  AND timestamp <= $4
ORDER BY timestamp ASC;

-- name: GetCandlesAfterTimestamp :many
SELECT * FROM candles
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp > $3
ORDER BY timestamp ASC
LIMIT $4;

-- name: GetCandlesBeforeTimestamp :many
SELECT * FROM candles
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp < $3
ORDER BY timestamp DESC
LIMIT $4;

-- name: ListCandlesPaginated :many
SELECT * FROM candles
WHERE symbol = $1
  AND timeframe = $2
ORDER BY timestamp DESC
LIMIT $3
OFFSET $4;

-- name: GetCandlesBySymbols :many
SELECT * FROM candles
WHERE symbol = ANY($1::varchar[])
  AND timeframe = $2
  AND timestamp >= $3
  AND timestamp <= $4
ORDER BY symbol, timestamp ASC;

-- =============================================================================
-- AGGREGATE Operations
-- =============================================================================

-- name: CountCandles :one
SELECT COUNT(*) FROM candles
WHERE symbol = $1
  AND timeframe = $2;

-- name: CountCandlesInRange :one
SELECT COUNT(*) FROM candles
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp >= $3
  AND timestamp <= $4;

-- name: GetCandleStats :one
SELECT 
    COUNT(*) as total_candles,
    MIN(timestamp) as first_timestamp,
    MAX(timestamp) as last_timestamp,
    MIN(low) as lowest_price,
    MAX(high) as highest_price,
    AVG(volume) as avg_volume
FROM candles
WHERE symbol = $1
  AND timeframe = $2;

-- name: GetCandleStatsInRange :one
SELECT 
    COUNT(*) as total_candles,
    MIN(low) as lowest_price,
    MAX(high) as highest_price,
    AVG(close) as avg_close,
    SUM(volume) as total_volume,
    AVG(volume) as avg_volume
FROM candles
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp >= $3
  AND timestamp <= $4;

-- name: GetVolumeLeaders :many
SELECT 
    symbol,
    timeframe,
    SUM(volume) as total_volume,
    AVG(volume) as avg_volume,
    MAX(volume) as max_volume
FROM candles
WHERE timestamp >= $1
  AND timestamp <= $2
GROUP BY symbol, timeframe
ORDER BY total_volume DESC
LIMIT $3;

-- =============================================================================
-- UPDATE Operations
-- =============================================================================

-- name: UpdateCandle :one
UPDATE candles
SET 
    open = $4,
    high = $5,
    low = $6,
    close = $7,
    volume = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp = $3
RETURNING *;

-- =============================================================================
-- DELETE Operations
-- =============================================================================

-- name: DeleteCandle :exec
DELETE FROM candles
WHERE id = $1;

-- name: DeleteCandlesBySymbol :exec
DELETE FROM candles
WHERE symbol = $1;

-- name: DeleteCandlesByTimeframe :exec
DELETE FROM candles
WHERE symbol = $1
  AND timeframe = $2;

-- name: DeleteOldCandles :exec
DELETE FROM candles
WHERE timestamp < $1;

-- name: DeleteCandlesInRange :exec
DELETE FROM candles
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp >= $3
  AND timestamp <= $4;

-- =============================================================================
-- UTILITY Operations
-- =============================================================================

-- name: GetDistinctSymbols :many
SELECT DISTINCT symbol
FROM candles
ORDER BY symbol;

-- name: GetDistinctTimeframes :many
SELECT DISTINCT timeframe
FROM candles
WHERE symbol = $1
ORDER BY timeframe;

-- name: GetSymbolTimeframePairs :many
SELECT DISTINCT symbol, timeframe
FROM candles
ORDER BY symbol, timeframe;

-- name: CheckCandleExists :one
SELECT EXISTS(
    SELECT 1 FROM candles
    WHERE symbol = $1
      AND timeframe = $2
      AND timestamp = $3
) AS exists;

-- name: GetOldestCandle :one
SELECT * FROM candles
WHERE symbol = $1
  AND timeframe = $2
ORDER BY timestamp ASC
LIMIT 1;

-- name: GetCandleGaps :many
SELECT 
    symbol,
    timeframe,
    timestamp as gap_start,
    LEAD(timestamp) OVER (ORDER BY timestamp) as gap_end,
    LEAD(timestamp) OVER (ORDER BY timestamp) - timestamp as gap_seconds
FROM candles
WHERE symbol = $1
  AND timeframe = $2
  AND timestamp >= $3
  AND timestamp <= $4
ORDER BY timestamp;
