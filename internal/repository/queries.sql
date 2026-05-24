-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, role, api_key)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users SET name = $2, role = $3, updated_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CreateProject :one
INSERT INTO projects (name, description, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id = $1;

-- name: ListProjects :many
SELECT * FROM projects ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects SET name = $2, description = $3, updated_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- name: CreateTestSuite :one
INSERT INTO test_suites (project_id, name, description, test_type, schedule_cron, config, tags, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetTestSuiteByID :one
SELECT * FROM test_suites WHERE id = $1;

-- name: ListTestSuitesByProject :many
SELECT * FROM test_suites WHERE project_id = $1 ORDER BY created_at DESC;

-- name: UpdateTestSuite :one
UPDATE test_suites SET name = $2, description = $3, test_type = $4, schedule_cron = $5, config = $6, tags = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTestSuite :exec
DELETE FROM test_suites WHERE id = $1;

-- name: CreateTestCase :one
INSERT INTO test_cases (suite_id, name, description, config, tags, sort_order, enabled)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTestCaseByID :one
SELECT * FROM test_cases WHERE id = $1;

-- name: ListTestCasesBySuite :many
SELECT * FROM test_cases WHERE suite_id = $1 ORDER BY sort_order;

-- name: UpdateTestCase :one
UPDATE test_cases SET name = $2, description = $3, config = $4, tags = $5, sort_order = $6, enabled = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTestCase :exec
DELETE FROM test_cases WHERE id = $1;

-- name: CreateTestRun :one
INSERT INTO test_runs (suite_id, project_id, status, trigger, triggered_by, total_cases, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTestRunByID :one
SELECT * FROM test_runs WHERE id = $1;

-- name: ListTestRunsByProject :many
SELECT * FROM test_runs WHERE project_id = $1 ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTestRunsBySuite :many
SELECT * FROM test_runs WHERE suite_id = $1 ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountTestRunsByProject :one
SELECT COUNT(*) FROM test_runs WHERE project_id = $1;

-- name: UpdateTestRunStatus :one
UPDATE test_runs SET status = $2, started_at = $3, finished_at = $4, duration_ms = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateTestRunCounts :one
UPDATE test_runs SET passed = $2, failed = $3, skipped = $4, errored = $5
WHERE id = $1
RETURNING *;

-- name: CreateTestResult :one
INSERT INTO test_results (run_id, case_id, worker_id, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTestResultByID :one
SELECT * FROM test_results WHERE id = $1;

-- name: ListTestResultsByRun :many
SELECT * FROM test_results WHERE run_id = $1 ORDER BY created_at;

-- name: UpdateTestResultStatus :one
UPDATE test_results SET status = $2, started_at = $3, finished_at = $4, duration_ms = $5
WHERE id = $1
RETURNING *;

-- name: UpdateTestResultData :one
UPDATE test_results SET status = $2, error_message = $3, assertions = $4, request_data = $5, response_data = $6, artifacts = $7, metrics = $8, stdout = $9, stderr = $10, duration_ms = $11, finished_at = $12
WHERE id = $1
RETURNING *;

-- name: CreateWorker :one
INSERT INTO workers (name, hostname, ip_address, executor_types, max_concurrent, version, labels)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetWorkerByID :one
SELECT * FROM workers WHERE id = $1;

-- name: ListWorkers :many
SELECT * FROM workers ORDER BY registered_at DESC;

-- name: ListWorkersByStatus :many
SELECT * FROM workers WHERE status = $1 ORDER BY registered_at DESC;

-- name: UpdateWorkerStatus :one
UPDATE workers SET status = $2, last_heartbeat = now(), updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateWorkerHeartbeat :exec
UPDATE workers SET status = $2, current_job_id = $3, last_heartbeat = now(), updated_at = now()
WHERE id = $1;

-- name: DeleteWorker :exec
DELETE FROM workers WHERE id = $1;

-- name: GetPendingResult :one
SELECT tr.* FROM test_results tr
JOIN test_runs r ON tr.run_id = r.id
WHERE tr.status = 'pending'
AND r.status IN ('pending', 'running')
ORDER BY tr.created_at ASC
FOR UPDATE SKIP LOCKED
LIMIT 1;

-- name: GetPendingResultForWorker :one
SELECT tr.* FROM test_results tr
JOIN test_runs r ON tr.run_id = r.id
JOIN test_suites ts ON r.suite_id = ts.id
WHERE tr.status = 'pending'
AND r.status IN ('pending', 'running')
AND ts.test_type = ANY($1::text[])
ORDER BY tr.created_at ASC
FOR UPDATE SKIP LOCKED
LIMIT 1;

-- name: CountPendingResults :one
SELECT COUNT(*) FROM test_results tr
JOIN test_runs r ON tr.run_id = r.id
WHERE tr.status = 'pending' AND r.status IN ('pending', 'running');

-- name: CountRunResultsByStatus :one
SELECT
    COUNT(*) FILTER (WHERE status = 'passed') AS passed,
    COUNT(*) FILTER (WHERE status = 'failed') AS failed,
    COUNT(*) FILTER (WHERE status = 'skipped') AS skipped,
    COUNT(*) FILTER (WHERE status = 'error') AS errored
FROM test_results WHERE run_id = $1;

-- name: CreateNotificationConfig :one
INSERT INTO notification_configs (project_id, type, config, events, enabled)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListNotificationConfigsByProject :many
SELECT * FROM notification_configs WHERE project_id = $1;

-- name: UpdateNotificationConfig :one
UPDATE notification_configs SET type = $2, config = $3, events = $4, enabled = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteNotificationConfig :exec
DELETE FROM notification_configs WHERE id = $1;

-- name: CreateMetricSnapshot :one
INSERT INTO metric_snapshots (project_id, suite_id, snapshot_date, total_runs, total_passed, total_failed, avg_duration_ms, p95_duration_ms, flaky_count, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (project_id, suite_id, snapshot_date)
DO UPDATE SET total_runs = $4, total_passed = $5, total_failed = $6, avg_duration_ms = $7, p95_duration_ms = $8, flaky_count = $9, metadata = $10, created_at = now()
RETURNING *;

-- name: ListMetricSnapshots :many
SELECT * FROM metric_snapshots
WHERE project_id = $1
AND ($2::uuid IS NULL OR suite_id = $2)
AND snapshot_date >= $3
AND snapshot_date <= $4
ORDER BY snapshot_date;

-- name: GetDashboardSummary :one
SELECT
    (SELECT COUNT(*) FROM test_runs WHERE project_id = $1 AND created_at >= $2) AS total_runs,
    (SELECT COUNT(*) FROM test_runs WHERE project_id = $1 AND status = 'passed' AND created_at >= $2) AS total_passed,
    (SELECT COUNT(*) FROM test_runs WHERE project_id = $1 AND status = 'failed' AND created_at >= $2) AS total_failed,
    (SELECT COUNT(*) FROM workers WHERE status = 'online') AS active_workers;

-- name: GetRecentRuns :many
SELECT tr.*, ts.name AS suite_name, p.name AS project_name
FROM test_runs tr
JOIN test_suites ts ON tr.suite_id = ts.id
JOIN projects p ON tr.project_id = p.id
WHERE ($1::uuid IS NULL OR tr.project_id = $1)
ORDER BY tr.created_at DESC
LIMIT $2;

-- name: GetFlakyTests :many
SELECT
    tc.id AS case_id,
    tc.name AS case_name,
    COUNT(*) AS total_runs,
    COUNT(*) FILTER (WHERE tr.status = 'passed') AS passed,
    COUNT(*) FILTER (WHERE tr.status = 'failed') AS failed
FROM test_results tr
JOIN test_cases tc ON tr.case_id = tc.id
JOIN test_runs r ON tr.run_id = r.id
WHERE ($1::uuid IS NULL OR r.project_id = $1)
AND r.created_at >= $2
GROUP BY tc.id, tc.name
HAVING COUNT(*) FILTER (WHERE tr.status = 'passed') > 0
   AND COUNT(*) FILTER (WHERE tr.status = 'failed') > 0
ORDER BY failed DESC
LIMIT $3;
