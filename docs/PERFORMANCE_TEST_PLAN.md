# Performance & Load Testing — Full Platform Improvement Plan

**Date:** 2026-05-25  
**Status:** Awaiting Approval

---

## 1. Executive Summary

The automation platform has a solid foundation — Go + Chi + PostgreSQL + distributed workers + React frontend — but **performance/load testing is entirely unimplemented** and several core quality issues block reliable everyday use. This document identifies every gap, explains its severity, and specifies exactly what to build.

---

## 2. Current State Audit

### 2.1 Critical Bugs / Missing Features

| # | Severity | Component | Issue |
|---|----------|-----------|-------|
| C1 | 🔴 Critical | `service/worker.go:90` | `GetNextJob()` always returns `nil` — **workers never receive jobs from the server**. Jobs are enqueued to an in-memory channel in the dispatcher but the handler stub ignores it. |
| C2 | 🔴 Critical | `worker/executors/` | Only `http.go` exists. `TestTypeLoad = "load"` is declared but **no load executor registered** — load tests silently fail with "no executor for test type: load". |
| C3 | 🟠 High | `service/dispatcher.go` | **In-memory channel queue** (buffer=100). All pending jobs are lost on any server restart. |
| C4 | 🟠 High | Docker networking | Worker container loses DNS resolution for `server` hostname every time the server container is recreated (fixed with `nc` healthcheck in progress). |
| C5 | 🟡 Medium | `handler/worker.go:106` | `GetNextJob` handler returns `StatusNoContent` with a JSON `null` body — the agent correctly handles 204, but the actual dispatcher dequeue is never called from this handler. |
| C6 | 🟡 Medium | `repository/repository.go` | `GetTrends` SQL scanned `DATE` into `string` and `AVG` into `int64` — causes HTTP 500 on dashboard (partially fixed). |
| C7 | 🟡 Medium | `TestSuites.tsx` | Run button had no `onClick` handler (fixed). |
| C8 | 🟡 Medium | `TestDataForm.tsx` | `ep.parameters` iteration crashed on null (fixed). |

### 2.2 Architectural Gaps for Load Testing

| # | Gap | Impact |
|---|-----|--------|
| G1 | No `LoadExecutor` | Load tests cannot run at all |
| G2 | No percentile/histogram metrics collection | No p50/p95/p99 latency data |
| G3 | No per-second time-series data | Cannot chart throughput over time |
| G4 | `test_results.metrics` is unstructured `JSONB` | No queryable load metrics |
| G5 | No load-specific assertions (error_rate, p95, throughput) | Cannot define SLOs |
| G6 | No load config UI | Users cannot configure VUs, duration, ramp-up |
| G7 | No load results dashboard | Cannot visualize test outcomes |
| G8 | WebSocket only broadcasts, no replay | Joining after run starts misses all events |
| G9 | No real-time streaming of metrics during run | No live feedback during long load tests |

---

## 3. Implementation Plan

### Phase 1 — Fix Critical Job Dispatch Bug (Blocker for Everything)

**Problem:** `GetNextJob` in `service/worker.go` is a stub. The dispatcher enqueues jobs but the handler never dequeues them.

**Fix:** Wire the handler to the dispatcher's `Dequeue()` method.

#### [MODIFY] `internal/service/worker.go`
```go
// Inject dispatcher dependency into WorkerService
type WorkerService struct {
    repo       WorkerRepository
    jwtSecret  []byte
    dispatcher interface{ Dequeue() *model.WorkerJob }
}

func (s *WorkerService) GetNextJob(ctx context.Context, workerID uuid.UUID) (*model.WorkerJob, error) {
    job := s.dispatcher.Dequeue()
    return job, nil  // nil means no jobs available — worker polls again
}
```

#### [MODIFY] `internal/server/deps.go` / server setup
Pass the dispatcher to `WorkerService` on construction.

---

### Phase 2 — Load Executor (Core Feature)

#### [NEW] `internal/worker/executors/load.go`

A concurrent HTTP load generator implementing the `worker.Executor` interface.

**Test Case Config Schema:**
```json
{
  "url": "https://api.example.com/users",
  "method": "GET",
  "headers": { "Authorization": "Bearer {{token}}" },
  "body": null,
  "load": {
    "vus": 50,
    "duration": "30s",
    "ramp_up": "10s",
    "ramp_down": "5s",
    "rate_limit_rps": 200
  },
  "assertions": [
    { "type": "p95_response_time", "expected": 500 },
    { "type": "error_rate",        "expected": 0.01 },
    { "type": "throughput_min",    "expected": 80 }
  ]
}
```

**Executor Implementation:**
1. Parse config into `LoadConfig` struct with validation.
2. Spawn goroutine pool of `vus` virtual users.
3. Ramp-up: linearly increase active VUs from 0 → vus over `ramp_up` duration using a ticker.
4. Rate limiting: per-VU token bucket (`golang.org/x/time/rate`) to enforce `rate_limit_rps`.
5. Each VU loops: make HTTP request → record `{status, duration_ms, bytes, error}` to a shared results channel.
6. Metrics collector goroutine: drains channel every second, computes per-second RPS/p95 for time series.
7. Ramp-down: reduce goroutines linearly over `ramp_down` before completing.
8. After `duration`: close VUs, aggregate all samples into `LoadMetrics`.

**Metrics Struct:**
```go
type LoadMetrics struct {
    TotalRequests  int64        `json:"total_requests"`
    SuccessCount   int64        `json:"success_count"`
    ErrorCount     int64        `json:"error_count"`
    ErrorRate      float64      `json:"error_rate"`
    ThroughputRPS  float64      `json:"throughput_rps"`
    MinLatencyMS   int64        `json:"min_latency_ms"`
    MaxLatencyMS   int64        `json:"max_latency_ms"`
    AvgLatencyMS   float64      `json:"avg_latency_ms"`
    P50LatencyMS   int64        `json:"p50_latency_ms"`
    P90LatencyMS   int64        `json:"p90_latency_ms"`
    P95LatencyMS   int64        `json:"p95_latency_ms"`
    P99LatencyMS   int64        `json:"p99_latency_ms"`
    TotalBytes     int64        `json:"total_bytes"`
    TimeSeriesRPS  []TimePoint  `json:"time_series_rps"`
    TimeSeriesP95  []TimePoint  `json:"time_series_p95"`
    StatusCodes    map[int]int  `json:"status_codes"`
}
type TimePoint struct {
    Second int     `json:"second"`
    Value  float64 `json:"value"`
}
```

**Load Assertions** (new types in `runLoadAssertions()`):
- `p95_response_time` → `p95 <= expected` ms
- `p99_response_time` → `p99 <= expected` ms
- `error_rate` → `error_rate <= expected` (0.01 = 1%)
- `throughput_min` → `throughput >= expected` RPS
- `avg_response_time` → `avg <= expected` ms

#### [MODIFY] `cmd/worker/main.go`
Register the new executor:
```go
registry.Register(executors.NewLoadExecutor())
registry.Register(executors.NewHTTPExecutor(30 * time.Second))
```

---

### Phase 3 — Database: Structured Load Metrics

#### [NEW] `migrations/000004_add_load_metrics.up.sql`
```sql
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
```

#### [MODIFY] `internal/repository/repository.go`
Add `LoadMetricsRepo` with `Create` and `GetByRun` methods.

#### [MODIFY] `internal/service/testrun.go` — `SubmitResult()`
When `result.Metrics` is non-empty and suite type is `load`, parse it as `LoadMetrics` and persist to `load_metrics` table.

---

### Phase 4 — Backend API: Load Metrics Endpoints

#### [MODIFY] `internal/server/routes.go`
```
GET /api/v1/runs/{runID}/load-metrics         → aggregate for entire run
GET /api/v1/runs/{runID}/results/{resultID}/load-metrics  → single result
```

#### [NEW] `internal/handler/loadmetrics.go`
Two handler functions for the above routes.

---

### Phase 5 — Real-time Metrics Streaming

During long load tests, the worker will push per-second metrics updates via the existing log endpoint.

#### [MODIFY] `internal/worker/agent.go`
Add a `streamMetrics(job Job, point TimePoint)` method that posts to `/api/v1/jobs/{jobID}/log` with `stream = "metrics"`.

#### [MODIFY] `internal/handler/worker.go` — `SubmitJobLog`
When `stream == "metrics"`, parse the content as a `TimePoint` and broadcast via WebSocket with `type = "load_metric"`.

#### [MODIFY] `internal/ws/hub.go`
No change needed — Hub already broadcasts to all subscribers.

---

### Phase 6 — Frontend: Load Test Configuration UI

#### [MODIFY] `web/src/pages/SuiteDetail.tsx`
For suites with `test_type === 'load'`, show a **Load Configuration** section when editing/adding test cases:
- Target URL + Method
- VUs: number input + preview label ("50 virtual users")
- Duration: text input (`30s`, `5m`, `1h`)
- Ramp Up / Ramp Down: text inputs
- Rate Limit: number input (req/s, 0 = unlimited)
- Load Assertions builder (p95, error_rate, throughput_min)

---

### Phase 7 — Frontend: Load Results Dashboard

#### [NEW] `web/src/components/load/LoadMetricsDashboard.tsx`
A rich results panel rendered inside `TestRunDetail` for load suite runs:

**Summary Cards Row:**
- Total Requests | RPS | Error Rate | Avg Latency | P95 Latency | P99 Latency

**Charts (recharts — already installed):**
- **Throughput Over Time**: AreaChart — RPS per second
- **Latency Percentiles Over Time**: LineChart — P50/P90/P95 per second
- **Status Code Distribution**: PieChart — 2xx/4xx/5xx breakdown
- **Response Time Histogram**: BarChart — bucketed by 0-100ms, 100-250ms, 250-500ms, 500ms-1s, >1s

**Live Streaming** (WebSocket):
- Subscribe to `load_metric` events while `run.status === 'running'`
- Append new `TimePoint` to chart data in real time

#### [MODIFY] `web/src/pages/TestRunDetail.tsx`
- Fetch suite details to know `test_type`
- If `test_type === 'load'`, render `<LoadMetricsDashboard runId={runId} />`
- Show polling/live badge while running (refresh every 3s)

---

## 4. Implementation Order

```
Step 1  Fix dispatcher wiring (Phase 1)             — unblocks all real test execution
Step 2  Load executor Go implementation (Phase 2)   — core load engine
Step 3  Register LoadExecutor in cmd/worker (Phase 2)
Step 4  DB migration 000004 (Phase 3)
Step 5  LoadMetricsRepo (Phase 3)
Step 6  SubmitResult integration for load metrics (Phase 3)
Step 7  New API routes + handlers (Phase 4)
Step 8  Real-time streaming support (Phase 5)
Step 9  Frontend: load config form in SuiteDetail (Phase 6)
Step 10 Frontend: LoadMetricsDashboard component (Phase 7)
Step 11 Frontend: TestRunDetail integration (Phase 7)
Step 12 Docker rebuild and end-to-end verification
```

---

## 5. Files Changed

| File | Action | Phase |
|------|--------|-------|
| `internal/service/worker.go` | MODIFY — wire dispatcher | 1 |
| `internal/service/dispatcher.go` | MODIFY — expose to WorkerService | 1 |
| `internal/worker/executors/load.go` | NEW | 2 |
| `cmd/worker/main.go` | MODIFY — register LoadExecutor | 2 |
| `internal/model/load.go` | NEW — LoadMetrics, TimePoint | 2 |
| `migrations/000004_add_load_metrics.up.sql` | NEW | 3 |
| `migrations/000004_add_load_metrics.down.sql` | NEW | 3 |
| `internal/repository/repository.go` | MODIFY — LoadMetricsRepo | 3 |
| `internal/service/testrun.go` | MODIFY — persist load metrics | 3 |
| `internal/handler/loadmetrics.go` | NEW | 4 |
| `internal/server/routes.go` | MODIFY — new routes | 4 |
| `internal/worker/agent.go` | MODIFY — stream metrics | 5 |
| `internal/handler/worker.go` | MODIFY — handle metrics stream | 5 |
| `web/src/components/load/LoadMetricsDashboard.tsx` | NEW | 7 |
| `web/src/api/loadMetrics.ts` | NEW | 7 |
| `web/src/pages/TestRunDetail.tsx` | MODIFY — load dashboard | 7 |
| `web/src/pages/SuiteDetail.tsx` | MODIFY — load config form | 6 |

---

## 6. Verification Plan

1. Create suite with `test_type = load`
2. Add test case pointing to the platform's own `/api/v1/health` endpoint
3. Config: `vus=20`, `duration=20s`, `ramp_up=5s`, assertions: `p95 <= 200ms`, `error_rate <= 0.01`
4. Click Run → redirects to run detail
5. **Verify real-time**: throughput and latency charts update every second
6. **Verify DB**: `SELECT * FROM load_metrics WHERE run_id = '...'` returns correct row
7. **Verify assertions**: p95 and error_rate pass/fail correctly
8. **Verify summary**: dashboard trends include the load run

