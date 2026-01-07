-- Metrics Queries
-- CRUD operations and queries for session analytics and metrics

-- =============================================================================
-- CREATE Operations
-- =============================================================================

-- name: CreateMetric :one
INSERT INTO metrics (
    session_id,
    timestamp,
    current_price,
    price_change,
    price_change_percent,
    session_high,
    session_low,
    session_open,
    realized_volatility,
    volatility_1min,
    volatility_5min,
    volatility_1hour,
    volume,
    volume_ma_5,
    volume_ma_20,
    drawdown_percent,
    max_drawdown_percent,
    return_percent,
    sma_10,
    sma_20,
    sma_50,
    sma_200,
    ema_10,
    ema_20,
    rsi_14,
    macd,
    macd_signal,
    macd_histogram,
    standard_deviation,
    variance,
    signal,
    signal_strength,
    custom_metrics
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33
)
RETURNING *;

-- name: CreateMetricBatch :copyfrom
INSERT INTO metrics (
    session_id,
    timestamp,
    current_price,
    volume,
    rsi_14,
    signal
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- =============================================================================
-- READ Operations
-- =============================================================================

-- name: GetMetricByID :one
SELECT * FROM metrics
WHERE id = $1;

-- name: GetMetricsBySession :many
SELECT * FROM metrics
WHERE session_id = $1
ORDER BY timestamp ASC;

-- name: GetRecentMetrics :many
SELECT * FROM metrics
WHERE session_id = $1
ORDER BY timestamp DESC
LIMIT $2;

-- name: GetMetricsInRange :many
SELECT * FROM metrics
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp ASC;

-- name: GetLatestMetric :one
SELECT * FROM metrics
WHERE session_id = $1
ORDER BY timestamp DESC
LIMIT 1;

-- name: GetMetricsPaginated :many
SELECT * FROM metrics
WHERE session_id = $1
ORDER BY timestamp DESC
LIMIT $2
OFFSET $3;

-- name: GetMetricsWithSignal :many
SELECT * FROM metrics
WHERE session_id = $1
  AND signal = $2
ORDER BY timestamp DESC
LIMIT $3;

-- name: GetMetricsBuySignals :many
SELECT * FROM metrics
WHERE session_id = $1
  AND signal = 'buy'
ORDER BY timestamp ASC;

-- name: GetMetricsSellSignals :many
SELECT * FROM metrics
WHERE session_id = $1
  AND signal = 'sell'
ORDER BY timestamp ASC;

-- name: GetMetricsAboveRSI :many
SELECT * FROM metrics
WHERE session_id = $1
  AND rsi_14 > $2
ORDER BY timestamp DESC
LIMIT $3;

-- name: GetMetricsBelowRSI :many
SELECT * FROM metrics
WHERE session_id = $1
  AND rsi_14 < $2
ORDER BY timestamp DESC
LIMIT $3;

-- =============================================================================
-- UPDATE Operations
-- =============================================================================

-- name: UpdateMetric :one
UPDATE metrics
SET 
    current_price = $3,
    price_change = $4,
    price_change_percent = $5,
    rsi_14 = $6,
    macd = $7,
    signal = $8
WHERE session_id = $1
  AND timestamp = $2
RETURNING *;

-- name: UpdateMetricSignal :exec
UPDATE metrics
SET 
    signal = $2,
    signal_strength = $3
WHERE id = $1;

-- =============================================================================
-- DELETE Operations
-- =============================================================================

-- name: DeleteMetric :exec
DELETE FROM metrics
WHERE id = $1;

-- name: DeleteMetricsBySession :exec
DELETE FROM metrics
WHERE session_id = $1;

-- name: DeleteOldMetrics :exec
DELETE FROM metrics
WHERE created_at < $1;

-- =============================================================================
-- AGGREGATE Operations
-- =============================================================================

-- name: CountMetrics :one
SELECT COUNT(*) FROM metrics
WHERE session_id = $1;

-- name: GetMetricsSummary :one
SELECT 
    COUNT(*) as total_metrics,
    MIN(current_price) as min_price,
    MAX(current_price) as max_price,
    AVG(current_price) as avg_price,
    MIN(rsi_14) as min_rsi,
    MAX(rsi_14) as max_rsi,
    AVG(rsi_14) as avg_rsi,
    AVG(realized_volatility) as avg_volatility,
    MAX(drawdown_percent) as max_drawdown,
    SUM(volume) as total_volume
FROM metrics
WHERE session_id = $1;

-- name: GetMetricsStatsInRange :one
SELECT 
    COUNT(*) as data_points,
    AVG(current_price) as avg_price,
    STDDEV(current_price) as price_stddev,
    AVG(rsi_14) as avg_rsi,
    AVG(realized_volatility) as avg_volatility,
    COUNT(*) FILTER (WHERE signal = 'buy') as buy_signals,
    COUNT(*) FILTER (WHERE signal = 'sell') as sell_signals
FROM metrics
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3;

-- name: GetVolatilityStats :one
SELECT 
    AVG(realized_volatility) as avg_realized_volatility,
    MAX(realized_volatility) as max_realized_volatility,
    MIN(realized_volatility) as min_realized_volatility,
    STDDEV(realized_volatility) as volatility_stddev
FROM metrics
WHERE session_id = $1;

-- name: GetDrawdownAnalysis :one
SELECT 
    MAX(drawdown_percent) as max_drawdown,
    AVG(drawdown_percent) as avg_drawdown,
    COUNT(*) FILTER (WHERE drawdown_percent > 5) as drawdowns_over_5pct,
    COUNT(*) FILTER (WHERE drawdown_percent > 10) as drawdowns_over_10pct
FROM metrics
WHERE session_id = $1;

-- name: GetMetricsSignalDistribution :many
SELECT 
    signal,
    COUNT(*) as signal_count,
    AVG(signal_strength) as avg_strength
FROM metrics
WHERE session_id = $1
  AND signal IS NOT NULL
GROUP BY signal
ORDER BY signal_count DESC;

-- name: GetPriceMovementStats :one
SELECT 
    AVG(price_change) as avg_price_change,
    STDDEV(price_change) as price_change_stddev,
    MAX(price_change) as max_price_increase,
    MIN(price_change) as max_price_decrease,
    AVG(ABS(price_change)) as avg_absolute_change
FROM metrics
WHERE session_id = $1;

-- =============================================================================
-- TIME SERIES Operations
-- =============================================================================

-- name: GetPriceSeries :many
SELECT 
    timestamp,
    current_price
FROM metrics
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp ASC;

-- name: GetRSISeries :many
SELECT 
    timestamp,
    rsi_14
FROM metrics
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp ASC;

-- name: GetVolumeSeries :many
SELECT 
    timestamp,
    volume,
    volume_ma_5,
    volume_ma_20
FROM metrics
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp ASC;

-- name: GetMovingAverages :many
SELECT 
    timestamp,
    current_price,
    sma_10,
    sma_20,
    sma_50,
    sma_200,
    ema_10,
    ema_20
FROM metrics
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp ASC;

-- name: GetMACDSeries :many
SELECT 
    timestamp,
    macd,
    macd_signal,
    macd_histogram
FROM metrics
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp ASC;

-- =============================================================================
-- ANALYSIS Operations
-- =============================================================================

-- name: GetMetricsRSIExtremes :many
SELECT * FROM metrics
WHERE session_id = $1
  AND (rsi_14 > 70 OR rsi_14 < 30)
ORDER BY timestamp ASC;

-- name: GetPriceCrossovers :many
SELECT 
    m1.timestamp,
    m1.current_price,
    m1.sma_20,
    m1.sma_50,
    CASE 
        WHEN m1.sma_20 > m1.sma_50 AND m2.sma_20 <= m2.sma_50 THEN 'golden_cross'
        WHEN m1.sma_20 < m1.sma_50 AND m2.sma_20 >= m2.sma_50 THEN 'death_cross'
    END as crossover_type
FROM metrics m1
JOIN metrics m2 ON m2.session_id = m1.session_id 
    AND m2.timestamp < m1.timestamp
WHERE m1.session_id = $1
  AND m1.sma_20 IS NOT NULL
  AND m1.sma_50 IS NOT NULL
  AND m2.sma_20 IS NOT NULL
  AND m2.sma_50 IS NOT NULL
  AND (
      (m1.sma_20 > m1.sma_50 AND m2.sma_20 <= m2.sma_50) OR
      (m1.sma_20 < m1.sma_50 AND m2.sma_20 >= m2.sma_50)
  )
ORDER BY m1.timestamp ASC;

-- name: GetHighVolatilityPeriods :many
SELECT * FROM metrics
WHERE session_id = $1
  AND realized_volatility > $2
ORDER BY realized_volatility DESC
LIMIT $3;

-- name: GetRecentPerformance :one
SELECT 
    COUNT(*) as period_count,
    (MAX(current_price) - MIN(current_price)) / MIN(current_price) * 100 as period_return,
    AVG(realized_volatility) as period_volatility,
    MAX(drawdown_percent) as period_max_drawdown
FROM metrics
WHERE session_id = $1
  AND timestamp >= $2;

-- =============================================================================
-- UTILITY Operations
-- =============================================================================

-- name: GetMetricTimestamps :many
SELECT DISTINCT timestamp
FROM metrics
WHERE session_id = $1
ORDER BY timestamp ASC;

-- name: CheckMetricExists :one
SELECT EXISTS(
    SELECT 1 FROM metrics 
    WHERE session_id = $1 
      AND timestamp = $2
) AS exists;

-- name: GetMetricCount :one
SELECT COUNT(*) 
FROM metrics
WHERE session_id = $1;

-- name: GetFirstMetric :one
SELECT * FROM metrics
WHERE session_id = $1
ORDER BY timestamp ASC
LIMIT 1;

-- name: GetLastMetric :one
SELECT * FROM metrics
WHERE session_id = $1
ORDER BY timestamp DESC
LIMIT 1;
