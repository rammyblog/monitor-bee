-- +goose Up
CREATE TABLE monitors(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    url  VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL DEFAULT 'GET',
    interval_seconds INTEGER NOT NULL DEFAULT 300,
    timeout_seconds INTEGER NOT NULL DEFAULT 30,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    headers JSONB,
    body TEXT,
    expected_status_code INTEGER DEFAULT 200,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE monitor_checks(
    id SERIAL PRIMARY KEY,
    monitor_id INTEGER NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,  
    status VARCHAR(50) NOT NULL,
    response_time_ms INTEGER,
    status_code INTEGER,
    error_message TEXT,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_monitors_user_id ON monitors(user_id);
CREATE INDEX idx_monitors_status ON monitors(status);
CREATE INDEX idx_monitor_checks_monitor_id ON monitor_checks(monitor_id);
CREATE INDEX idx_monitor_checks_checked_at ON monitor_checks(checked_at);

-- +goose Down
DROP INDEX IF EXISTS idx_monitor_checks_checked_at;
DROP INDEX IF EXISTS idx_monitor_checks_monitor_id;
DROP INDEX IF EXISTS idx_monitors_status;
DROP INDEX IF EXISTS idx_monitors_user_id;
DROP TABLE IF EXISTS monitor_checks;
DROP TABLE IF EXISTS monitors;
