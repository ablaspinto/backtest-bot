-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL,
    candle_id BIGINT,
    
    -- Event identification
    event_type VARCHAR(50) NOT NULL,
    event_category VARCHAR(30) DEFAULT 'market' 
        CHECK (event_category IN ('market', 'technical', 'system', 'user', 'alert')),
    
    -- Event details
    title VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Timing
    timestamp BIGINT NOT NULL,
    duration_seconds INTEGER,
    
    -- Severity/Impact
    severity VARCHAR(20) DEFAULT 'medium' 
        CHECK (severity IN ('critical', 'high', 'medium', 'low', 'info')),
    impact_score NUMERIC(5, 4) CHECK (impact_score >= 0 AND impact_score <= 1),
    
    -- Price information at event time
    price_at_event NUMERIC(20, 8),
    price_before NUMERIC(20, 8),
    price_after NUMERIC(20, 8),
    price_change_percent NUMERIC(10, 4),
    
    -- Volume information
    volume_at_event NUMERIC(20, 8),
    volume_spike_ratio NUMERIC(10, 4),
    
    -- Event triggers
    triggered_by VARCHAR(100),
    trigger_conditions JSONB DEFAULT '{}'::jsonb,
    
    -- Event metadata
    tags VARCHAR(50)[],
    indicators_at_event JSONB DEFAULT '{}'::jsonb,
    additional_data JSONB DEFAULT '{}'::jsonb,
    
    -- User interaction
    is_user_marked BOOLEAN DEFAULT false,
    user_notes TEXT,
    is_acknowledged BOOLEAN DEFAULT false,
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    
    -- Event status
    status VARCHAR(20) DEFAULT 'active' 
        CHECK (status IN ('active', 'resolved', 'dismissed', 'archived')),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign keys
    CONSTRAINT fk_events_session 
        FOREIGN KEY (session_id) 
        REFERENCES sessions(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_events_candle 
        FOREIGN KEY (candle_id) 
        REFERENCES candles(id) 
        ON DELETE SET NULL,
    
    -- Constraints
    CONSTRAINT chk_duration_positive CHECK (duration_seconds IS NULL OR duration_seconds >= 0),
    CONSTRAINT chk_price_change CHECK (
        (price_before IS NULL AND price_after IS NULL) OR
        (price_before IS NOT NULL AND price_after IS NOT NULL)
    )
);

-- Primary lookup index (session + time)
CREATE INDEX idx_events_session_timestamp 
    ON events(session_id, timestamp DESC);

-- Index for event type queries
CREATE INDEX idx_events_type 
    ON events(event_type);

-- Index for event category
CREATE INDEX idx_events_category 
    ON events(event_category);

-- Index for severity queries
CREATE INDEX idx_events_severity 
    ON events(severity) 
    WHERE severity IN ('critical', 'high');

-- Index for timestamp-based queries
CREATE INDEX idx_events_timestamp 
    ON events(timestamp DESC);

-- Index for candle lookups
CREATE INDEX idx_events_candle 
    ON events(candle_id) 
    WHERE candle_id IS NOT NULL;

-- Index for user-marked events
CREATE INDEX idx_events_user_marked 
    ON events(is_user_marked) 
    WHERE is_user_marked = true;

-- Index for unacknowledged events
CREATE INDEX idx_events_unacknowledged 
    ON events(session_id, is_acknowledged) 
    WHERE is_acknowledged = false;

-- Index for status queries
CREATE INDEX idx_events_status 
    ON events(status) 
    WHERE status = 'active';

-- GIN index for tags
CREATE INDEX idx_events_tags 
    ON events USING GIN (tags);

-- GIN indexes for JSONB fields
CREATE INDEX idx_events_trigger_conditions 
    ON events USING GIN (trigger_conditions);

CREATE INDEX idx_events_indicators 
    ON events USING GIN (indicators_at_event);

CREATE INDEX idx_events_additional_data 
    ON events USING GIN (additional_data);

-- Composite index for common event queries
CREATE INDEX idx_events_active_by_severity 
    ON events(session_id, severity, timestamp DESC) 
    WHERE status = 'active';

-- Add table comments
COMMENT ON TABLE events IS 'Tracks significant market events, alerts, and user markers during simulations';
COMMENT ON COLUMN events.event_type IS 'Type of event (spike, crash, breakout, consolidation, volume_spike, ma_cross, signal_generated, etc.)';
COMMENT ON COLUMN events.event_category IS 'Broad category for grouping events (market, technical, system, user, alert)';
COMMENT ON COLUMN events.severity IS 'Importance level of the event';
COMMENT ON COLUMN events.impact_score IS 'Normalized score (0-1) indicating event impact magnitude';
COMMENT ON COLUMN events.triggered_by IS 'What caused this event (indicator name, system component, user action)';
COMMENT ON COLUMN events.trigger_conditions IS 'JSON object with conditions that triggered the event';
COMMENT ON COLUMN events.volume_spike_ratio IS 'Ratio of current volume to average (e.g., 2.5 = 250% of normal)';
COMMENT ON COLUMN events.indicators_at_event IS 'Snapshot of indicator values when event occurred';
COMMENT ON COLUMN events.additional_data IS 'Flexible field for event-specific data';
COMMENT ON COLUMN events.duration_seconds IS 'How long the event lasted (NULL for instantaneous events)';

-- Create a view for active high-priority events
CREATE OR REPLACE VIEW v_active_critical_events AS
SELECT 
    e.id,
    e.session_id,
    s.session_name,
    s.symbol,
    e.event_type,
    e.event_category,
    e.title,
    e.severity,
    e.price_at_event,
    e.price_change_percent,
    e.timestamp,
    e.created_at,
    e.is_acknowledged
FROM events e
JOIN sessions s ON e.session_id = s.id
WHERE e.status = 'active'
  AND e.severity IN ('critical', 'high')
  AND e.is_acknowledged = false
ORDER BY e.severity DESC, e.timestamp DESC;

COMMENT ON VIEW v_active_critical_events IS 'Quick view of unacknowledged high-priority events requiring attention';

-- Create a view for event summary by session
CREATE OR REPLACE VIEW v_event_summary_by_session AS
SELECT 
    session_id,
    COUNT(*) as total_events,
    COUNT(*) FILTER (WHERE severity = 'critical') as critical_events,
    COUNT(*) FILTER (WHERE severity = 'high') as high_events,
    COUNT(*) FILTER (WHERE severity = 'medium') as medium_events,
    COUNT(*) FILTER (WHERE severity = 'low') as low_events,
    COUNT(*) FILTER (WHERE event_category = 'market') as market_events,
    COUNT(*) FILTER (WHERE event_category = 'technical') as technical_events,
    COUNT(*) FILTER (WHERE event_category = 'user') as user_events,
    COUNT(*) FILTER (WHERE is_user_marked = true) as user_marked_count,
    AVG(impact_score) as avg_impact_score,
    MAX(timestamp) as last_event_timestamp
FROM events
GROUP BY session_id;

COMMENT ON VIEW v_event_summary_by_session IS 'Aggregated event statistics per simulation session';

-- Create a function to auto-acknowledge old events
CREATE OR REPLACE FUNCTION auto_acknowledge_old_events(days_old INTEGER DEFAULT 7)
RETURNS INTEGER AS $$
DECLARE
    rows_updated INTEGER;
BEGIN
    UPDATE events
    SET is_acknowledged = true,
        acknowledged_at = NOW(),
        status = 'archived'
    WHERE is_acknowledged = false
      AND created_at < NOW() - (days_old || ' days')::INTERVAL
      AND severity NOT IN ('critical', 'high');
    
    GET DIAGNOSTICS rows_updated = ROW_COUNT;
    RETURN rows_updated;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION auto_acknowledge_old_events IS 'Automatically acknowledges and archives non-critical events older than N days';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS auto_acknowledge_old_events;
DROP VIEW IF EXISTS v_event_summary_by_session;
DROP VIEW IF EXISTS v_active_critical_events;
DROP TABLE IF EXISTS events CASCADE;
-- +goose StatementEnd
