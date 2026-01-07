-- Sessions Queries
-- CRUD operations and queries for simulation sessions

-- =============================================================================
-- CREATE Operations
-- =============================================================================

-- name: CreateSession :one
INSERT INTO sessions (
    session_name,
    symbol,
    timeframe,
    strategy,
    starting_price,
    parameters,
    notes,
    tags
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- =============================================================================
-- READ Operations
-- =============================================================================

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1;

-- name: GetSessionByName :one
SELECT * FROM sessions
WHERE session_name = $1
LIMIT 1;

-- name: ListSessions :many
SELECT * FROM sessions
ORDER BY started_at DESC
LIMIT $1
OFFSET $2;

-- name: ListAllSessions :many
SELECT * FROM sessions
ORDER BY started_at DESC;

-- name: GetActiveSessions :many
SELECT * FROM sessions
WHERE status = 'active'
ORDER BY started_at DESC;

-- name: GetCompletedSessions :many
SELECT * FROM sessions
WHERE status = 'completed'
ORDER BY ended_at DESC
LIMIT $1;

-- name: GetFavoriteSessions :many
SELECT * FROM sessions
WHERE is_favorite = true
ORDER BY started_at DESC;

-- name: GetSessionsBySymbol :many
SELECT * FROM sessions
WHERE symbol = $1
ORDER BY started_at DESC
LIMIT $2;

-- name: GetSessionsByStrategy :many
SELECT * FROM sessions
WHERE strategy = $1
ORDER BY started_at DESC
LIMIT $2;

-- name: GetSessionsByStatus :many
SELECT * FROM sessions
WHERE status = $1
ORDER BY started_at DESC
LIMIT $2;

-- name: GetRecentSessions :many
SELECT * FROM sessions
WHERE started_at >= $1
ORDER BY started_at DESC;

-- name: SearchSessionsByName :many
SELECT * FROM sessions
WHERE session_name ILIKE '%' || $1 || '%'
ORDER BY started_at DESC
LIMIT $2;

-- name: GetSessionsWithTags :many
SELECT * FROM sessions
WHERE tags && $1::varchar[]
ORDER BY started_at DESC
LIMIT $2;

-- =============================================================================
-- UPDATE Operations
-- =============================================================================

-- name: UpdateSession :one
UPDATE sessions
SET 
    session_name = $2,
    notes = $3,
    tags = $4,
    is_favorite = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateSessionEnd :one
UPDATE sessions
SET 
    ended_at = $2,
    ending_price = $3,
    highest_price = $4,
    lowest_price = $5,
    total_candles = $6,
    price_change_percent = $7,
    volatility = $8,
    status = $9,
    duration_seconds = EXTRACT(EPOCH FROM ($2 - started_at))::INTEGER,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateSessionStatus :exec
UPDATE sessions
SET 
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateSessionPrices :exec
UPDATE sessions
SET 
    highest_price = GREATEST(highest_price, $2),
    lowest_price = LEAST(COALESCE(lowest_price, $3), $3),
    total_candles = total_candles + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ToggleSessionFavorite :one
UPDATE sessions
SET 
    is_favorite = NOT is_favorite,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING is_favorite;

-- name: AddSessionTag :exec
UPDATE sessions
SET 
    tags = array_append(tags, $2),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND NOT ($2 = ANY(tags));

-- name: RemoveSessionTag :exec
UPDATE sessions
SET 
    tags = array_remove(tags, $2),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- =============================================================================
-- DELETE Operations
-- =============================================================================

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteOldSessions :exec
DELETE FROM sessions
WHERE ended_at < $1
  AND status = 'completed';

-- name: DeleteSessionsByStatus :exec
DELETE FROM sessions
WHERE status = $1;

-- =============================================================================
-- AGGREGATE Operations
-- =============================================================================

-- name: CountSessions :one
SELECT COUNT(*) FROM sessions;

-- name: CountSessionsByStatus :one
SELECT COUNT(*) FROM sessions
WHERE status = $1;

-- name: GetSessionStats :one
SELECT 
    COUNT(*) as total_sessions,
    COUNT(*) FILTER (WHERE status = 'completed') as completed_sessions,
    COUNT(*) FILTER (WHERE status = 'active') as active_sessions,
    COUNT(*) FILTER (WHERE status = 'paused') as paused_sessions,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_sessions,
    AVG(price_change_percent) FILTER (WHERE status = 'completed') as avg_price_change,
    AVG(volatility) FILTER (WHERE status = 'completed') as avg_volatility,
    AVG(duration_seconds) FILTER (WHERE status = 'completed') as avg_duration_seconds
FROM sessions;

-- name: GetSymbolStats :many
SELECT 
    symbol,
    COUNT(*) as session_count,
    AVG(price_change_percent) FILTER (WHERE status = 'completed') as avg_price_change,
    MAX(ending_price) as max_ending_price,
    MIN(ending_price) as min_ending_price
FROM sessions
GROUP BY symbol
ORDER BY session_count DESC;

-- name: GetStrategyPerformance :many
SELECT 
    strategy,
    COUNT(*) as total_runs,
    COUNT(*) FILTER (WHERE status = 'completed') as completed_runs,
    AVG(price_change_percent) FILTER (WHERE status = 'completed') as avg_return,
    AVG(volatility) FILTER (WHERE status = 'completed') as avg_volatility,
    MAX(price_change_percent) FILTER (WHERE status = 'completed') as max_return,
    MIN(price_change_percent) FILTER (WHERE status = 'completed') as min_return
FROM sessions
GROUP BY strategy
ORDER BY avg_return DESC;

-- name: GetDailySessionCount :many
SELECT 
    DATE(started_at) as session_date,
    COUNT(*) as session_count
FROM sessions
WHERE started_at >= $1
GROUP BY DATE(started_at)
ORDER BY session_date DESC;

-- =============================================================================
-- UTILITY Operations
-- =============================================================================

-- name: CheckSessionExists :one
SELECT EXISTS(
    SELECT 1 FROM sessions WHERE id = $1
) AS exists;

-- name: GetLatestSession :one
SELECT * FROM sessions
ORDER BY started_at DESC
LIMIT 1;

-- name: GetLongestSessions :many
SELECT * FROM sessions
WHERE status = 'completed'
  AND duration_seconds IS NOT NULL
ORDER BY duration_seconds DESC
LIMIT $1;

-- name: GetMostProfitableSessions :many
SELECT * FROM sessions
WHERE status = 'completed'
  AND price_change_percent IS NOT NULL
ORDER BY price_change_percent DESC
LIMIT $1;

-- name: GetMostVolatileSessions :many
SELECT * FROM sessions
WHERE status = 'completed'
  AND volatility IS NOT NULL
ORDER BY volatility DESC
LIMIT $1;

-- name: GetSessionDuration :one
SELECT 
    COALESCE(duration_seconds, EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - started_at))::INTEGER) as duration_seconds
FROM sessions
WHERE id = $1;
