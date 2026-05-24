-- Environments: manage base URLs and variables per environment (dev/staging/prod)
CREATE TABLE environments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    base_url    TEXT NOT NULL,
    variables   JSONB NOT NULL DEFAULT '{}',
    headers     JSONB NOT NULL DEFAULT '{}',
    is_default  BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(project_id, name)
);

CREATE INDEX idx_environments_project ON environments(project_id);

-- Add environment_id to test_runs so we know which env was used
ALTER TABLE test_runs ADD COLUMN environment_id UUID REFERENCES environments(id);

-- Add retry configuration to test_suites and test_cases
ALTER TABLE test_suites ADD COLUMN max_retries INTEGER NOT NULL DEFAULT 0;
ALTER TABLE test_suites ADD COLUMN retry_delay_ms INTEGER NOT NULL DEFAULT 1000;
ALTER TABLE test_suites ADD COLUMN parallel BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE test_suites ADD COLUMN max_parallel INTEGER NOT NULL DEFAULT 5;
ALTER TABLE test_suites ADD COLUMN timeout_ms INTEGER NOT NULL DEFAULT 300000;

ALTER TABLE test_cases ADD COLUMN max_retries INTEGER NOT NULL DEFAULT 0;
ALTER TABLE test_cases ADD COLUMN timeout_ms INTEGER NOT NULL DEFAULT 60000;
ALTER TABLE test_cases ADD COLUMN pre_script TEXT;
ALTER TABLE test_cases ADD COLUMN post_script TEXT;

-- Datasets for data-driven testing
CREATE TABLE datasets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    format      VARCHAR(20) NOT NULL DEFAULT 'json',
    data        JSONB NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_datasets_project ON datasets(project_id);

-- Link datasets to test cases
ALTER TABLE test_cases ADD COLUMN dataset_id UUID REFERENCES datasets(id);

-- Webhook configurations for triggering runs externally
CREATE TABLE webhooks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    secret      VARCHAR(255) NOT NULL,
    suite_id    UUID REFERENCES test_suites(id) ON DELETE SET NULL,
    enabled     BOOLEAN NOT NULL DEFAULT true,
    last_triggered_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_webhooks_project ON webhooks(project_id);

-- Schedule tracking
CREATE TABLE schedule_runs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    suite_id    UUID NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
    cron_expr   VARCHAR(100) NOT NULL,
    next_run_at TIMESTAMPTZ NOT NULL,
    last_run_at TIMESTAMPTZ,
    enabled     BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_schedule_runs_next ON schedule_runs(next_run_at) WHERE enabled = true;

-- Add retry tracking to test_results
ALTER TABLE test_results ADD COLUMN retry_count INTEGER NOT NULL DEFAULT 0;
