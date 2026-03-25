CREATE TABLE IF NOT EXISTS pipeline_stages (
    stage VARCHAR(50) PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    metadata JSONB DEFAULT '{}'
);
