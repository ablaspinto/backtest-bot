-- Indicators Queries
-- CRUD operations and queries for technical indicators

-- =============================================================================
-- CREATE Operations
-- =============================================================================

-- name: CreateIndicator :one
INSERT INTO indicators (
    candle_id,
    indicator_type,
    value,
    value_upper,
    value_lower,
    value_middle,
    signal_line,
    histogram,
    period,
    period_fast,
    period_slow,
    period_signal,
    standard_deviations,
    signal,
    signal_strength,
    is_crossover,
    crossover_type,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18
)
RETURNING *;

-- name: CreateIndicatorBatch :copyfrom
INSERT INTO indicators (
    candle_id,
    indicator_type,
    value,
    period,
    signal
) VALUES (
    $1, $2, $3, $4, $5
);

-- =============================================================================
-- READ Operations
-- =============================================================================

-- name: GetIndicatorByID :one
SELECT * FROM indicators
WHERE id = $1;

-- name: GetIndicatorsByCandle :many
SELECT * FROM indicators
WHERE candle_id = $1
ORDER BY indicator_type, period;

-- name: GetIndicatorByType :one
SELECT * FROM indicators
WHERE candle_id = $1
  AND indicator_type = $2
  AND period = COALESCE($3, period)
LIMIT 1;

-- name: GetIndicatorsByType :many
SELECT * FROM indicators
WHERE indicator_type = $1
ORDER BY calculated_at DESC
LIMIT $2;

-- name: GetIndicatorsForCandles :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND c.timestamp >= $3
  AND c.timestamp <= $4
ORDER BY c.timestamp ASC, i.indicator_type;

-- name: GetSMAIndicators :many
SELECT 
    i.*,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = 'SMA'
  AND i.period = $3
  AND c.timestamp >= $4
  AND c.timestamp <= $5
ORDER BY c.timestamp ASC;

-- name: GetEMAIndicators :many
SELECT 
    i.*,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = 'EMA'
  AND i.period = $3
  AND c.timestamp >= $4
  AND c.timestamp <= $5
ORDER BY c.timestamp ASC;

-- name: GetRSIIndicators :many
SELECT 
    i.*,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = 'RSI'
  AND i.period = $3
  AND c.timestamp >= $4
  AND c.timestamp <= $5
ORDER BY c.timestamp ASC;

-- name: GetMACDIndicators :many
SELECT 
    i.*,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = 'MACD'
  AND c.timestamp >= $3
  AND c.timestamp <= $4
ORDER BY c.timestamp ASC;

-- name: GetBollingerBands :many
SELECT 
    i.*,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = 'BB'
  AND i.period = $3
  AND c.timestamp >= $4
  AND c.timestamp <= $5
ORDER BY c.timestamp ASC;

-- =============================================================================
-- UPDATE Operations
-- =============================================================================

-- name: UpdateIndicator :one
UPDATE indicators
SET 
    value = $2,
    value_upper = $3,
    value_lower = $4,
    signal = $5,
    signal_strength = $6
WHERE id = $1
RETURNING *;

-- name: UpdateIndicatorSignal :exec
UPDATE indicators
SET 
    signal = $2,
    signal_strength = $3
WHERE id = $1;

-- name: MarkCrossover :exec
UPDATE indicators
SET 
    is_crossover = true,
    crossover_type = $2
WHERE id = $1;

-- =============================================================================
-- DELETE Operations
-- =============================================================================

-- name: DeleteIndicator :exec
DELETE FROM indicators
WHERE id = $1;

-- name: DeleteIndicatorsByCandle :exec
DELETE FROM indicators
WHERE candle_id = $1;

-- name: DeleteIndicatorsByType :exec
DELETE FROM indicators
WHERE indicator_type = $1;

-- name: DeleteOldIndicators :exec
DELETE FROM indicators
WHERE calculated_at < $1;

-- =============================================================================
-- SIGNAL Operations
-- =============================================================================

-- name: GetBuySignals :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.signal = 'buy'
  AND c.timestamp >= $3
ORDER BY c.timestamp DESC
LIMIT $4;

-- name: GetSellSignals :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.signal = 'sell'
  AND c.timestamp >= $3
ORDER BY c.timestamp DESC
LIMIT $4;

-- name: GetOverboughtIndicators :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.signal = 'overbought'
  AND c.timestamp >= $3
ORDER BY c.timestamp DESC
LIMIT $4;

-- name: GetOversoldIndicators :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.signal = 'oversold'
  AND c.timestamp >= $3
ORDER BY c.timestamp DESC
LIMIT $4;

-- name: GetSignalsByStrength :many
SELECT 
    i.*,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.signal_strength >= $3
  AND i.signal IN ('buy', 'sell')
ORDER BY i.signal_strength DESC, c.timestamp DESC
LIMIT $4;

-- =============================================================================
-- CROSSOVER Operations
-- =============================================================================

-- name: GetCrossovers :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.is_crossover = true
  AND c.timestamp >= $3
ORDER BY c.timestamp DESC;

-- name: GetGoldenCrosses :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.crossover_type = 'golden'
  AND c.timestamp >= $3
ORDER BY c.timestamp DESC
LIMIT $4;

-- name: GetDeathCrosses :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.crossover_type = 'death'
  AND c.timestamp >= $3
ORDER BY c.timestamp DESC
LIMIT $4;

-- name: GetRecentCrossovers :many
SELECT 
    i.*,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.is_crossover = true
ORDER BY c.timestamp DESC
LIMIT $3;

-- =============================================================================
-- AGGREGATE Operations
-- =============================================================================

-- name: CountIndicators :one
SELECT COUNT(*) FROM indicators
WHERE candle_id = $1;

-- name: CountIndicatorsByType :one
SELECT COUNT(*) FROM indicators
WHERE indicator_type = $1;

-- name: GetIndicatorStats :one
SELECT 
    indicator_type,
    COUNT(*) as total_count,
    AVG(value) as avg_value,
    MIN(value) as min_value,
    MAX(value) as max_value,
    COUNT(*) FILTER (WHERE signal = 'buy') as buy_signals,
    COUNT(*) FILTER (WHERE signal = 'sell') as sell_signals,
    COUNT(*) FILTER (WHERE is_crossover = true) as crossovers
FROM indicators
WHERE indicator_type = $1
GROUP BY indicator_type;

-- name: GetSignalDistribution :many
SELECT 
    indicator_type,
    signal,
    COUNT(*) as signal_count,
    AVG(signal_strength) as avg_strength
FROM indicators
WHERE signal IS NOT NULL
GROUP BY indicator_type, signal
ORDER BY indicator_type, signal_count DESC;

-- name: GetCrossoverFrequency :many
SELECT 
    indicator_type,
    crossover_type,
    COUNT(*) as crossover_count
FROM indicators
WHERE is_crossover = true
  AND crossover_type IS NOT NULL
GROUP BY indicator_type, crossover_type
ORDER BY crossover_count DESC;

-- =============================================================================
-- ANALYSIS Operations
-- =============================================================================

-- name: GetRSIExtremes :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = 'RSI'
  AND (i.value > $3 OR i.value < $4)
  AND c.timestamp >= $5
ORDER BY c.timestamp DESC;

-- name: GetBollingerBreakouts :many
SELECT 
    i.*,
    c.symbol,
    c.timeframe,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = 'BB'
  AND (c.close > i.value_upper OR c.close < i.value_lower)
  AND c.timestamp >= $3
ORDER BY c.timestamp DESC
LIMIT $4;

-- name: GetMAConvergence :many
SELECT 
    i1.candle_id,
    i1.value as fast_ma,
    i2.value as slow_ma,
    ABS(i1.value - i2.value) as convergence,
    c.timestamp,
    c.close
FROM indicators i1
JOIN indicators i2 ON i2.candle_id = i1.candle_id
JOIN candles c ON c.id = i1.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i1.indicator_type = 'SMA'
  AND i2.indicator_type = 'SMA'
  AND i1.period = $3
  AND i2.period = $4
  AND c.timestamp >= $5
  AND c.timestamp <= $6
ORDER BY c.timestamp ASC;

-- =============================================================================
-- UTILITY Operations
-- =============================================================================

-- name: GetDistinctIndicatorTypes :many
SELECT DISTINCT indicator_type
FROM indicators
ORDER BY indicator_type;

-- name: GetIndicatorPeriods :many
SELECT DISTINCT period
FROM indicators
WHERE indicator_type = $1
  AND period IS NOT NULL
ORDER BY period;

-- name: CheckIndicatorExists :one
SELECT EXISTS(
    SELECT 1 FROM indicators
    WHERE candle_id = $1
      AND indicator_type = $2
      AND period = $3
) AS exists;

-- name: GetLatestIndicator :one
SELECT 
    i.*,
    c.timestamp
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = $3
  AND i.period = $4
ORDER BY c.timestamp DESC
LIMIT 1;

-- name: GetIndicatorHistory :many
SELECT 
    i.*,
    c.timestamp,
    c.close
FROM indicators i
JOIN candles c ON c.id = i.candle_id
WHERE c.symbol = $1
  AND c.timeframe = $2
  AND i.indicator_type = $3
ORDER BY c.timestamp DESC
LIMIT $4;
