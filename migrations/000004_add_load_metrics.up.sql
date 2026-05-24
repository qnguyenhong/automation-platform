CREATE TABLE load_metrics (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    result_id       UUID NOT NULL REFERENCES test_results(id) ON DELETE CASCADE,
    run_id          UUID NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    total_requests  BIGINT NOT NULL DEFAULT 0,
    success_count   BIGINT NOT NULL DEFAULT 0,
    error_count     BIGINT NOT NULL DEFAULT 0,
    error_rate      FLOAT NOT NULL DEFAULT 0,
    throughput_rps  FLOAT NOT NULL DEFAULT 0,
    min_latency_ms  BIGINT NOT NULL DEFAULT 0,
    max_latency_ms  BIGINT NOT NULL DEFAULT 0,
    avg_latency_ms  FLOAT NOT NULL DEFAULT 0,
    p50_latency_ms  BIGINT NOT NULL DEFAULT 0,
    p90_latency_ms  BIGINT NOT NULL DEFAULT 0,
    p95_latency_ms  BIGINT NOT NULL DEFAULT 0,
    p99_latency_ms  BIGINT NOT NULL DEFAULT 0,
    total_bytes     BIGINT NOT NULL DEFAULT 0,
    time_series     JSONB NOT NULL DEFAULT '[]',
    status_codes    JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_load_metrics_run ON load_metrics(run_id);
