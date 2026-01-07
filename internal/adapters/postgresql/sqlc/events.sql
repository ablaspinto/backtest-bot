-- Events Queries
-- CRUD operations and queries for market events and alerts

-- =============================================================================
-- CREATE Operations
-- =============================================================================

-- name: CreateEvent :one
INSERT INTO events (
    session_id,
    candle_id,
    event_type,
    event_category,
    title,
    description,
    timestamp,
    duration_seconds,
    severity,
    impact_score,
    price_at_event,
    price_before,
    price_after,
    price_change_percent,
    volume_at_event,
    volume_spike_ratio,
    triggered_by,
    trigger_conditions,
    tags,
    indicators_at_event,
    additional_data,
    is_user_marked,
    user_notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23
)
RETURNING *;

-- name: CreateEventSimple :one
INSERT INTO events (
    session_id,
    event_type,
    event_category,
    title,
    timestamp,
    severity,
    price_at_event
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- =============================================================================
-- READ Operations
-- =============================================================================

-- name: GetEventByID :one
SELECT * FROM events
WHERE id = $1;

-- name: GetEventsBySession :many
SELECT * FROM events
WHERE session_id = $1
ORDER BY timestamp DESC;

-- name: GetEventsBySessionPaginated :many
SELECT * FROM events
WHERE session_id = $1
ORDER BY timestamp DESC
LIMIT $2
OFFSET $3;

-- name: GetRecentEvents :many
SELECT * FROM events
WHERE session_id = $1
ORDER BY timestamp DESC
LIMIT $2;

-- name: GetEventsInRange :many
SELECT * FROM events
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp ASC;

-- name: GetEventsByType :many
SELECT * FROM events
WHERE session_id = $1
  AND event_type = $2
ORDER BY timestamp DESC
LIMIT $3;

-- name: GetEventsByCategory :many
SELECT * FROM events
WHERE session_id = $1
  AND event_category = $2
ORDER BY timestamp DESC
LIMIT $3;

-- name: GetEventsBySeverity :many
SELECT * FROM events
WHERE session_id = $1
  AND severity = $2
ORDER BY timestamp DESC
LIMIT $3;

-- name: GetCriticalEvents :many
SELECT * FROM events
WHERE session_id = $1
  AND severity IN ('critical', 'high')
  AND status = 'active'
ORDER BY severity DESC, timestamp DESC;

-- name: GetUnacknowledgedEvents :many
SELECT * FROM events
WHERE session_id = $1
  AND is_acknowledged = false
ORDER BY severity DESC, timestamp DESC;

-- name: GetUserMarkedEvents :many
SELECT * FROM events
WHERE session_id = $1
  AND is_user_marked = true
ORDER BY timestamp DESC;

-- name: GetActiveEvents :many
SELECT * FROM events
WHERE session_id = $1
  AND status = 'active'
ORDER BY timestamp DESC;

-- name: GetEventsWithTags :many
SELECT * FROM events
WHERE session_id = $1
  AND tags && $2::varchar[]
ORDER BY timestamp DESC
LIMIT $3;

-- name: SearchEvents :many
SELECT * FROM events
WHERE session_id = $1
  AND (
    title ILIKE '%' || $2 || '%' OR
    description ILIKE '%' || $2 || '%'
  )
ORDER BY timestamp DESC
LIMIT $3;

-- =============================================================================
-- UPDATE Operations
-- =============================================================================

-- name: UpdateEvent :one
UPDATE events
SET 
    title = $2,
    description = $3,
    severity = $4,
    user_notes = $5,
    tags = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: AcknowledgeEvent :exec
UPDATE events
SET 
    is_acknowledged = true,
    acknowledged_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: AcknowledgeAllEvents :exec
UPDATE events
SET 
    is_acknowledged = true,
    acknowledged_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE session_id = $1
  AND is_acknowledged = false;

-- name: UpdateEventStatus :exec
UPDATE events
SET 
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ResolveEvent :exec
UPDATE events
SET 
    status = 'resolved',
    is_acknowledged = true,
    acknowledged_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DismissEvent :exec
UPDATE events
SET 
    status = 'dismissed',
    is_acknowledged = true,
    acknowledged_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ArchiveOldEvents :execrows
UPDATE events
SET 
    status = 'archived',
    is_acknowledged = true,
    acknowledged_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE created_at < $1
  AND severity NOT IN ('critical', 'high')
  AND is_acknowledged = false;

-- name: AddEventTag :exec
UPDATE events
SET 
    tags = array_append(tags, $2),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND NOT ($2 = ANY(tags));

-- name: RemoveEventTag :exec
UPDATE events
SET 
    tags = array_remove(tags, $2),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ToggleUserMark :one
UPDATE events
SET 
    is_user_marked = NOT is_user_marked,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING is_user_marked;

-- name: UpdateEventNotes :exec
UPDATE events
SET 
    user_notes = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- =============================================================================
-- DELETE Operations
-- =============================================================================

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;

-- name: DeleteEventsBySession :exec
DELETE FROM events
WHERE session_id = $1;

-- name: DeleteArchivedEvents :exec
DELETE FROM events
WHERE status = 'archived'
  AND created_at < $1;

-- name: DeleteDismissedEvents :exec
DELETE FROM events
WHERE status = 'dismissed'
  AND created_at < $1;

-- =============================================================================
-- AGGREGATE Operations
-- =============================================================================

-- name: CountEvents :one
SELECT COUNT(*) FROM events
WHERE session_id = $1;

-- name: CountEventsByType :one
SELECT COUNT(*) FROM events
WHERE session_id = $1
  AND event_type = $2;

-- name: CountEventsBySeverity :one
SELECT COUNT(*) FROM events
WHERE session_id = $1
  AND severity = $2;

-- name: CountUnacknowledgedEvents :one
SELECT COUNT(*) FROM events
WHERE session_id = $1
  AND is_acknowledged = false;

-- name: GetEventsSummary :one
SELECT 
    COUNT(*) as total_events,
    COUNT(*) FILTER (WHERE severity = 'critical') as critical_count,
    COUNT(*) FILTER (WHERE severity = 'high') as high_count,
    COUNT(*) FILTER (WHERE severity = 'medium') as medium_count,
    COUNT(*) FILTER (WHERE severity = 'low') as low_count,
    COUNT(*) FILTER (WHERE is_acknowledged = false) as unacknowledged_count,
    COUNT(*) FILTER (WHERE is_user_marked = true) as user_marked_count,
    AVG(impact_score) as avg_impact_score,
    MAX(impact_score) as max_impact_score
FROM events
WHERE session_id = $1;

-- name: GetEventTypeDistribution :many
SELECT 
    event_type,
    COUNT(*) as event_count,
    AVG(impact_score) as avg_impact
FROM events
WHERE session_id = $1
GROUP BY event_type
ORDER BY event_count DESC;

-- name: GetEventCategoryDistribution :many
SELECT 
    event_category,
    COUNT(*) as event_count,
    COUNT(*) FILTER (WHERE severity IN ('critical', 'high')) as high_severity_count
FROM events
WHERE session_id = $1
GROUP BY event_category
ORDER BY event_count DESC;

-- name: GetSeverityDistribution :many
SELECT 
    severity,
    COUNT(*) as event_count,
    AVG(impact_score) as avg_impact
FROM events
WHERE session_id = $1
GROUP BY severity
ORDER BY 
    CASE severity
        WHEN 'critical' THEN 1
        WHEN 'high' THEN 2
        WHEN 'medium' THEN 3
        WHEN 'low' THEN 4
        WHEN 'info' THEN 5
    END;

-- name: GetEventFrequencyByHour :many
SELECT 
    EXTRACT(HOUR FROM TO_TIMESTAMP(timestamp)) as hour,
    COUNT(*) as event_count
FROM events
WHERE session_id = $1
GROUP BY hour
ORDER BY hour;

-- name: GetTopEventTriggers :many
SELECT 
    triggered_by,
    COUNT(*) as trigger_count,
    AVG(impact_score) as avg_impact
FROM events
WHERE session_id = $1
  AND triggered_by IS NOT NULL
GROUP BY triggered_by
ORDER BY trigger_count DESC
LIMIT $2;

-- =============================================================================
-- ANALYSIS Operations
-- =============================================================================

-- name: GetHighImpactEvents :many
SELECT * FROM events
WHERE session_id = $1
  AND impact_score >= $2
ORDER BY impact_score DESC, timestamp DESC
LIMIT $3;

-- name: GetPriceSpikes :many
SELECT * FROM events
WHERE session_id = $1
  AND event_type IN ('price_spike', 'price_crash')
  AND ABS(price_change_percent) >= $2
ORDER BY ABS(price_change_percent) DESC
LIMIT $3;

-- name: GetVolumeSpikes :many
SELECT * FROM events
WHERE session_id = $1
  AND event_type = 'volume_spike'
  AND volume_spike_ratio >= $2
ORDER BY volume_spike_ratio DESC
LIMIT $3;

-- name: GetBreakoutEvents :many
SELECT * FROM events
WHERE session_id = $1
  AND event_type IN ('breakout', 'breakdown')
ORDER BY timestamp DESC
LIMIT $2;

-- name: GetEventsCausingDrawdown :many
SELECT * FROM events
WHERE session_id = $1
  AND price_change_percent < 0
  AND ABS(price_change_percent) >= $2
ORDER BY price_change_percent ASC
LIMIT $3;

-- name: GetCriticalEventsNeedingAction :many
SELECT * FROM events
WHERE session_id = $1
  AND severity = 'critical'
  AND status = 'active'
  AND is_acknowledged = false
ORDER BY timestamp DESC;

-- name: GetEventCorrelations :many
SELECT 
    e1.event_type as event_type_1,
    e2.event_type as event_type_2,
    COUNT(*) as correlation_count,
    AVG(e2.timestamp - e1.timestamp) as avg_time_diff
FROM events e1
JOIN events e2 ON e2.session_id = e1.session_id
    AND e2.timestamp > e1.timestamp
    AND e2.timestamp <= e1.timestamp + $2
WHERE e1.session_id = $1
  AND e1.event_type != e2.event_type
GROUP BY e1.event_type, e2.event_type
HAVING COUNT(*) >= $3
ORDER BY correlation_count DESC;

-- =============================================================================
-- TIME SERIES Operations
-- =============================================================================

-- name: GetEventTimeline :many
SELECT 
    timestamp,
    event_type,
    severity,
    title,
    price_at_event
FROM events
WHERE session_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp ASC;

-- name: GetEventDensity :many
SELECT 
    (timestamp / $2) * $2 as time_bucket,
    COUNT(*) as event_count
FROM events
WHERE session_id = $1
GROUP BY time_bucket
ORDER BY time_bucket ASC;

-- =============================================================================
-- UTILITY Operations
-- =============================================================================

-- name: CheckEventExists :one
SELECT EXISTS(
    SELECT 1 FROM events
    WHERE id = $1
) AS exists;

-- name: GetDistinctEventTypes :many
SELECT DISTINCT event_type
FROM events
WHERE session_id = $1
ORDER BY event_type;

-- name: GetDistinctEventCategories :many
SELECT DISTINCT event_category
FROM events
ORDER BY event_category;

-- name: GetLatestEvent :one
SELECT * FROM events
WHERE session_id = $1
ORDER BY timestamp DESC
LIMIT 1;

-- name: GetFirstEvent :one
SELECT * FROM events
WHERE session_id = $1
ORDER BY timestamp ASC
LIMIT 1;

-- name: GetEventsByCandle :many
SELECT * FROM events
WHERE candle_id = $1
ORDER BY timestamp DESC;

-- name: GetEventsNearTimestamp :many
SELECT * FROM events
WHERE session_id = $1
  AND timestamp BETWEEN $2 - $3 AND $2 + $3
ORDER BY ABS(timestamp - $2) ASC
LIMIT $4;
