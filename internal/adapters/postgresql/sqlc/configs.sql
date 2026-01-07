-- Configs Queries
-- CRUD operations and queries for simulation configurations

-- =============================================================================
-- CREATE Operations
-- =============================================================================

-- name: CreateConfig :one
INSERT INTO configs (
    name,
    description,
    strategy,
    volatility,
    drift,
    starting_price,
    candle_interval,
    symbol,
    market_behavior,
    trend_strength,
    mean_reversion_speed,
    advanced_params,
    tags
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- =============================================================================
-- READ Operations
-- =============================================================================

-- name: GetConfigByID :one
SELECT * FROM configs
WHERE id = $1;

-- name: GetConfigByName :one
SELECT * FROM configs
WHERE name = $1
LIMIT 1;

-- name: GetDefaultConfig :one
SELECT * FROM configs
WHERE is_default = true
LIMIT 1;

-- name: ListConfigs :many
SELECT * FROM configs
ORDER BY times_used DESC, created_at DESC;

-- name: ListConfigsPaginated :many
SELECT * FROM configs
ORDER BY times_used DESC, created_at DESC
LIMIT $1
OFFSET $2;

-- name: GetFavoriteConfigs :many
SELECT * FROM configs
WHERE is_favorite = true
ORDER BY times_used DESC;

-- name: GetPublicConfigs :many
SELECT * FROM configs
WHERE is_public = true
ORDER BY times_used DESC
LIMIT $1;

-- name: GetConfigsByStrategy :many
SELECT * FROM configs
WHERE strategy = $1
ORDER BY created_at DESC;

-- name: GetConfigsByMarketBehavior :many
SELECT * FROM configs
WHERE market_behavior = $1
ORDER BY created_at DESC;

-- name: GetConfigsByCreator :many
SELECT * FROM configs
WHERE created_by = $1
ORDER BY created_at DESC;

-- name: SearchConfigsByName :many
SELECT * FROM configs
WHERE name ILIKE '%' || $1 || '%'
ORDER BY times_used DESC
LIMIT $2;

-- name: GetConfigsWithTags :many
SELECT * FROM configs
WHERE tags && $1::varchar[]
ORDER BY times_used DESC
LIMIT $2;

-- name: GetRecentlyUsedConfigs :many
SELECT * FROM configs
WHERE last_used_at IS NOT NULL
ORDER BY last_used_at DESC
LIMIT $1;

-- name: GetPopularConfigs :many
SELECT * FROM configs
WHERE times_used > 0
ORDER BY times_used DESC
LIMIT $1;

-- =============================================================================
-- UPDATE Operations
-- =============================================================================

-- name: UpdateConfig :one
UPDATE configs
SET 
    name = $2,
    description = $3,
    strategy = $4,
    volatility = $5,
    drift = $6,
    starting_price = $7,
    candle_interval = $8,
    symbol = $9,
    market_behavior = $10,
    trend_strength = $11,
    mean_reversion_speed = $12,
    advanced_params = $13,
    tags = $14,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateConfigPartial :one
UPDATE configs
SET 
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    volatility = COALESCE(sqlc.narg('volatility'), volatility),
    drift = COALESCE(sqlc.narg('drift'), drift),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: IncrementConfigUsage :exec
UPDATE configs
SET 
    times_used = times_used + 1,
    last_used_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: SetDefaultConfig :exec
UPDATE configs
SET is_default = (id = $1),
    updated_at = CURRENT_TIMESTAMP;

-- name: ToggleConfigFavorite :one
UPDATE configs
SET 
    is_favorite = NOT is_favorite,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING is_favorite;

-- name: ToggleConfigPublic :one
UPDATE configs
SET 
    is_public = NOT is_public,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING is_public;

-- name: AddConfigTag :exec
UPDATE configs
SET 
    tags = array_append(tags, $2),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND NOT ($2 = ANY(tags));

-- name: RemoveConfigTag :exec
UPDATE configs
SET 
    tags = array_remove(tags, $2),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- =============================================================================
-- DELETE Operations
-- =============================================================================

-- name: DeleteConfig :exec
DELETE FROM configs
WHERE id = $1;

-- name: DeleteUnusedConfigs :exec
DELETE FROM configs
WHERE times_used = 0
  AND is_favorite = false
  AND is_default = false
  AND created_at < $1;

-- =============================================================================
-- AGGREGATE Operations
-- =============================================================================

-- name: CountConfigs :one
SELECT COUNT(*) FROM configs;

-- name: CountConfigsByStrategy :one
SELECT COUNT(*) FROM configs
WHERE strategy = $1;

-- name: GetConfigStats :one
SELECT 
    COUNT(*) as total_configs,
    COUNT(*) FILTER (WHERE is_favorite = true) as favorite_count,
    COUNT(*) FILTER (WHERE is_public = true) as public_count,
    COUNT(*) FILTER (WHERE is_default = true) as default_count,
    AVG(times_used) as avg_times_used,
    MAX(times_used) as max_times_used,
    AVG(volatility) as avg_volatility
FROM configs;

-- name: GetStrategyDistribution :many
SELECT 
    strategy,
    COUNT(*) as config_count,
    AVG(volatility) as avg_volatility,
    AVG(drift) as avg_drift
FROM configs
GROUP BY strategy
ORDER BY config_count DESC;

-- name: GetMarketBehaviorDistribution :many
SELECT 
    market_behavior,
    COUNT(*) as config_count
FROM configs
GROUP BY market_behavior
ORDER BY config_count DESC;

-- =============================================================================
-- UTILITY Operations
-- =============================================================================

-- name: CheckConfigExists :one
SELECT EXISTS(
    SELECT 1 FROM configs WHERE id = $1
) AS exists;

-- name: CheckConfigNameExists :one
SELECT EXISTS(
    SELECT 1 FROM configs WHERE name = $1
) AS exists;

-- name: GetConfigByNameOrDefault :one
SELECT * FROM configs
WHERE name = $1
   OR (is_default = true AND $1 = '')
ORDER BY (name = $1) DESC
LIMIT 1;

-- name: DuplicateConfig :one
INSERT INTO configs (
    name,
    description,
    strategy,
    volatility,
    drift,
    starting_price,
    candle_interval,
    symbol,
    market_behavior,
    trend_strength,
    mean_reversion_speed,
    advanced_params,
    tags,
    created_by
)
SELECT 
    $2 as name,
    c.description,
    c.strategy,
    c.volatility,
    c.drift,
    c.starting_price,
    c.candle_interval,
    c.symbol,
    c.market_behavior,
    c.trend_strength,
    c.mean_reversion_speed,
    c.advanced_params,
    c.tags,
    c.created_by
FROM configs c
WHERE c.id = $1
RETURNING *;

SELECT c.* FROM configs c
WHERE c.strategy = (SELECT c2.strategy FROM configs c2 WHERE c2.id = $1)
  AND ABS(c.volatility - (SELECT c3.volatility FROM configs c3 WHERE c3.id = $1)) < $2
ORDER BY c.times_used DESC
LIMIT $3;
