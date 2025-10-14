-- name: CreateMonitor :one
INSERT INTO monitors (
    user_id,
    name,
    url,
    method,
    interval_seconds,
    timeout_seconds,
    status,
    headers,
    body,
    expected_status_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at;

-- name: GetMonitor :one
SELECT id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at
FROM monitors
WHERE id = $1 LIMIT 1;

-- name: GetMonitorByID :one
SELECT id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at
FROM monitors
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: ListMonitors :many
SELECT id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at
FROM monitors
ORDER BY created_at DESC;

-- name: ListMonitorsByUser :many
SELECT id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at
FROM monitors
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListActiveMonitors :many
SELECT id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at
FROM monitors
WHERE status = 'active'
ORDER BY created_at DESC;

-- name: ListMonitorsByStatus :many
SELECT id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at
FROM monitors
WHERE status = $1
ORDER BY created_at DESC;

-- name: ListMonitorsByUserAndStatus :many
SELECT id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at
FROM monitors
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC;

-- name: UpdateMonitor :one
UPDATE monitors
SET name = $2,
    url = $3,
    method = $4,
    interval_seconds = $5,
    timeout_seconds = $6,
    headers = $7,
    body = $8,
    expected_status_code = $9,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $10
RETURNING id, user_id, name, url, method, interval_seconds, timeout_seconds, status, headers, body, expected_status_code, created_at, updated_at;

-- name: UpdateMonitorStatus :exec
UPDATE monitors
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteMonitor :exec
DELETE FROM monitors
WHERE id = $1 AND user_id = $2;

-- name: DeleteMonitorByID :exec
DELETE FROM monitors
WHERE id = $1;

-- name: MonitorExists :one
SELECT EXISTS(SELECT 1 FROM monitors WHERE id = $1);

-- name: UserOwnsMonitor :one
SELECT EXISTS(SELECT 1 FROM monitors WHERE id = $1 AND user_id = $2);

-- name: CountMonitorsByUser :one
SELECT COUNT(*) FROM monitors WHERE user_id = $1;

-- name: CountActiveMonitorsByUser :one
SELECT COUNT(*) FROM monitors WHERE user_id = $1 AND status = 'active';


-- name: CreateMonitorCheck :one
INSERT INTO monitor_checks (
    monitor_id,
    status,
    response_time_ms,
    status_code,
    error_message
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, monitor_id, status, response_time_ms, status_code, error_message, checked_at;

-- name: GetMonitorCheck :one
SELECT id, monitor_id, status, response_time_ms, status_code, error_message, checked_at
FROM monitor_checks
WHERE id = $1 LIMIT 1;

-- name: GetLatestMonitorCheck :one
SELECT id, monitor_id, status, response_time_ms, status_code, error_message, checked_at
FROM monitor_checks
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT 1;

-- name: ListMonitorChecks :many
SELECT id, monitor_id, status, response_time_ms, status_code, error_message, checked_at
FROM monitor_checks
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT $2 OFFSET $3;

-- name: ListMonitorChecksByDateRange :many
SELECT id, monitor_id, status, response_time_ms, status_code, error_message, checked_at
FROM monitor_checks
WHERE monitor_id = $1
    AND checked_at >= $2
    AND checked_at <= $3
ORDER BY checked_at DESC;

-- name: ListRecentMonitorChecks :many
SELECT id, monitor_id, status, response_time_ms, status_code, error_message, checked_at
FROM monitor_checks
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT $2;

-- name: ListFailedMonitorChecks :many
SELECT id, monitor_id, status, response_time_ms, status_code, error_message, checked_at
FROM monitor_checks
WHERE monitor_id = $1 AND status = 'failed'
ORDER BY checked_at DESC
LIMIT $2;

-- name: CountMonitorChecks :one
SELECT COUNT(*) FROM monitor_checks WHERE monitor_id = $1;

-- name: CountFailedMonitorChecks :one
SELECT COUNT(*) FROM monitor_checks WHERE monitor_id = $1 AND status = 'failed';

-- name: CountSuccessfulMonitorChecks :one
SELECT COUNT(*) FROM monitor_checks WHERE monitor_id = $1 AND status = 'success';

-- name: GetAverageResponseTime :one
SELECT AVG(response_time_ms) as avg_response_time
FROM monitor_checks
WHERE monitor_id = $1 AND response_time_ms IS NOT NULL;

-- name: GetAverageResponseTimeByDateRange :one
SELECT AVG(response_time_ms) as avg_response_time
FROM monitor_checks
WHERE monitor_id = $1
    AND response_time_ms IS NOT NULL
    AND checked_at >= $2
    AND checked_at <= $3;

-- name: GetMonitorUptime :one
SELECT
    COUNT(CASE WHEN status = 'success' THEN 1 END)::float / COUNT(*)::float * 100 as uptime_percentage
FROM monitor_checks
WHERE monitor_id = $1;

-- name: GetMonitorUptimeByDateRange :one
SELECT
    COUNT(CASE WHEN status = 'success' THEN 1 END)::float / COUNT(*)::float * 100 as uptime_percentage
FROM monitor_checks
WHERE monitor_id = $1
    AND checked_at >= $2
    AND checked_at <= $3;

-- name: DeleteMonitorCheck :exec
DELETE FROM monitor_checks
WHERE id = $1;

-- name: DeleteOldMonitorChecks :exec
DELETE FROM monitor_checks
WHERE checked_at < $1;

-- name: DeleteMonitorChecksByMonitorID :exec
DELETE FROM monitor_checks
WHERE monitor_id = $1;

-- name: GetMonitorStats :one
SELECT
    COUNT(*) as total_checks,
    COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_checks,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_checks,
    AVG(response_time_ms) as avg_response_time,
    MIN(response_time_ms) as min_response_time,
    MAX(response_time_ms) as max_response_time
FROM monitor_checks
WHERE monitor_id = $1;
