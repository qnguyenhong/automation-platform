-- Users
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL DEFAULT 'member',
    api_key         VARCHAR(64) UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Projects
CREATE TABLE projects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Test Suites
CREATE TABLE test_suites (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    test_type       VARCHAR(50) NOT NULL,
    schedule_cron   VARCHAR(100),
    config          JSONB NOT NULL DEFAULT '{}',
    tags            TEXT[] DEFAULT '{}',
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Test Cases
CREATE TABLE test_cases (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    suite_id        UUID NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    config          JSONB NOT NULL DEFAULT '{}',
    tags            TEXT[] DEFAULT '{}',
    sort_order      INTEGER NOT NULL DEFAULT 0,
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Workers
CREATE TABLE workers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    hostname        VARCHAR(255) NOT NULL,
    ip_address      VARCHAR(45),
    executor_types  TEXT[] NOT NULL DEFAULT '{}',
    status          VARCHAR(50) NOT NULL DEFAULT 'offline',
    current_job_id  UUID,
    max_concurrent  INTEGER NOT NULL DEFAULT 1,
    version         VARCHAR(50),
    labels          JSONB DEFAULT '{}',
    last_heartbeat  TIMESTAMPTZ,
    registered_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Test Runs
CREATE TABLE test_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    suite_id        UUID NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    status          VARCHAR(50) NOT NULL DEFAULT 'pending',
    trigger         VARCHAR(50) NOT NULL DEFAULT 'manual',
    triggered_by    UUID REFERENCES users(id),
    total_cases     INTEGER NOT NULL DEFAULT 0,
    passed          INTEGER NOT NULL DEFAULT 0,
    failed          INTEGER NOT NULL DEFAULT 0,
    skipped         INTEGER NOT NULL DEFAULT 0,
    errored         INTEGER NOT NULL DEFAULT 0,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    duration_ms     BIGINT,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Test Results
CREATE TABLE test_results (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id          UUID NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    case_id         UUID NOT NULL REFERENCES test_cases(id),
    worker_id       UUID REFERENCES workers(id),
    status          VARCHAR(50) NOT NULL DEFAULT 'pending',
    error_message   TEXT,
    assertions      JSONB DEFAULT '[]',
    request_data    JSONB,
    response_data   JSONB,
    artifacts       JSONB DEFAULT '[]',
    metrics         JSONB,
    stdout          TEXT,
    stderr          TEXT,
    duration_ms     BIGINT,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    retry_of        UUID REFERENCES test_results(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Notification Configs
CREATE TABLE notification_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL,
    config          JSONB NOT NULL,
    events          TEXT[] NOT NULL DEFAULT '{run_failed}',
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Metric Snapshots
CREATE TABLE metric_snapshots (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    suite_id        UUID REFERENCES test_suites(id) ON DELETE CASCADE,
    snapshot_date   DATE NOT NULL,
    total_runs      INTEGER NOT NULL DEFAULT 0,
    total_passed    INTEGER NOT NULL DEFAULT 0,
    total_failed    INTEGER NOT NULL DEFAULT 0,
    avg_duration_ms BIGINT,
    p95_duration_ms BIGINT,
    flaky_count     INTEGER NOT NULL DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(project_id, suite_id, snapshot_date)
);

-- Indexes
CREATE INDEX idx_test_suites_project ON test_suites(project_id);
CREATE INDEX idx_test_cases_suite ON test_cases(suite_id);
CREATE INDEX idx_test_runs_project_created ON test_runs(project_id, created_at DESC);
CREATE INDEX idx_test_runs_suite_created ON test_runs(suite_id, created_at DESC);
CREATE INDEX idx_test_runs_status ON test_runs(status) WHERE status IN ('pending', 'running');
CREATE INDEX idx_test_results_run ON test_results(run_id);
CREATE INDEX idx_test_results_case ON test_results(case_id);
CREATE INDEX idx_test_results_status ON test_results(status);
CREATE INDEX idx_workers_status ON workers(status);
CREATE INDEX idx_metric_snapshots_project_date ON metric_snapshots(project_id, snapshot_date DESC);
