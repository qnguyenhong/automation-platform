# Automation Platform - Improvement Plan

## Current State Summary

The platform currently provides:
- **Backend**: Go + Chi router + PostgreSQL + JWT auth + WebSocket
- **Frontend**: React + Vite + TailwindCSS + Zustand
- **Core Features**: Projects, Test Suites, Test Cases, OpenAPI import (parse + endpoints), Distributed Workers, Job Dispatcher, Variable Engine (`$uuid`, `$random_*`, `$env`), HTTP Executor with assertions (status, json_path, header, response_time), Dashboard with trends, Notifications, Real-time WebSocket updates

---

## Phase 1: Core Testing Experience (High Priority)

| # | Area | Improvement | Details |
|---|------|-------------|---------|
| 1.1 | Environment Management | Add environments (dev/staging/prod) with base URLs and variables | New `environments` table; UI to switch env before run; resolves `{{base_url}}` per env |
| 1.2 | Test Case Editor | Rich visual editor for creating/editing test cases | Form-based editor: method selector, URL builder, headers/params/body editors, assertion builder with autocomplete |
| 1.3 | Test Chaining | Sequential execution with variable capture between steps | Pass `captured_vars` from step N to step N+1; UI to show variable flow |
| 1.4 | Data-Driven Testing | Parameterized test runs with CSV/JSON datasets | New `datasets` table; loop test case over rows; report per-row results |
| 1.5 | Pre/Post Request Scripts | JavaScript-like scripting (or Go template) for setup/teardown | Pre-request: set auth tokens, compute HMAC; Post-request: extract values, conditional assertions |

---

## Phase 2: OpenAPI Integration Enhancement

| # | Area | Improvement | Details |
|---|------|-------------|---------|
| 2.1 | $ref Resolution | Fully resolve `$ref` pointers in schemas | Currently schemas may have unresolved `$ref`; implement recursive resolution |
| 2.2 | Auth Scheme Detection | Auto-detect security schemes from spec | Parse `securitySchemes`, auto-populate Authorization header templates |
| 2.3 | Smart Test Generation | Auto-generate negative tests from spec | Generate 400/401/404 scenarios from required fields, enum validations |
| 2.4 | Spec Sync | Re-import/update when spec changes | Track spec version; diff endpoints; mark deprecated/new endpoints |
| 2.5 | Request Body Validation | Validate request/response against schema | Assert response matches OpenAPI schema definition (schema-based assertion type) |

---

## Phase 3: Execution & Reliability

| # | Area | Improvement | Details |
|---|------|-------------|---------|
| 3.1 | Retry Logic | Auto-retry failed/flaky tests | Config `max_retries` per suite/case; `retry_of` field already exists in schema |
| 3.2 | Scheduled Runs | Implement cron scheduling | `schedule_cron` field exists in test_suites; add scheduler service using cron library |
| 3.3 | Parallel Execution | Control concurrency within a suite | Config `parallel: true/false` + `max_parallel` at suite level |
| 3.4 | Timeout Configuration | Per-case and per-suite timeout settings | Currently hardcoded 5min; make configurable in suite/case config |
| 3.5 | Persistent Queue | Replace in-memory channel with Redis/DB queue | Current dispatcher loses jobs on restart; use Redis or DB-backed queue |

---

## Phase 4: CI/CD & External Integration

| # | Area | Improvement | Details |
|---|------|-------------|---------|
| 4.1 | API Key Auth | Allow triggering runs via API key | `api_key` field exists on users; implement header-based auth for CI pipelines |
| 4.2 | CLI Enhancements | CLI commands to trigger runs, check status | `cmd/cli` exists; add `run`, `status`, `results`, `export` commands |
| 4.3 | Webhook Triggers | Trigger runs from external events | New endpoint `/api/v1/webhooks/trigger` with project+suite filters |
| 4.4 | CI Plugins | GitHub Actions / GitLab CI integration | Provide YAML templates using CLI or API to run tests in pipelines |
| 4.5 | Export Formats | Export test cases to Postman/cURL/k6 | Generate Postman collections, curl commands, or k6 scripts from suites |

---

## Phase 5: Observability & Reporting

| # | Area | Improvement | Details |
|---|------|-------------|---------|
| 5.1 | Detailed Reporting | Rich HTML/PDF reports with charts | Generate per-run reports with pass/fail charts, response time histograms |
| 5.2 | History Comparison | Compare two runs side-by-side | Show regressions, new failures, fixed tests between runs |
| 5.3 | Request/Response Viewer | Pretty-printed req/resp in results page | Syntax-highlighted JSON, headers table, timing breakdown |
| 5.4 | Real-time Logs | Stream execution logs via WebSocket | Workers push logs; UI shows live output during test execution |
| 5.5 | Flaky Test Detection | Improve flaky detection algorithm | Track pass/fail patterns over N runs; auto-tag as flaky; suggest quarantine |

---

## Phase 6: Frontend UX Improvements

| # | Area | Improvement | Details |
|---|------|-------------|---------|
| 6.1 | Dark Mode | Full theme support | Complete dark mode with theme toggle in header |
| 6.2 | Drag & Drop Ordering | Reorder test cases with DnD | Update `sort_order` field; use `@dnd-kit` library |
| 6.3 | Bulk Operations | Multi-select cases for enable/disable/delete/tag | Checkbox selection with action toolbar |
| 6.4 | Search & Filter | Full-text search across cases, filter by tag/status | Backend search endpoint; frontend debounced search input |
| 6.5 | Keyboard Shortcuts | Quick actions (Ctrl+R = run, Ctrl+N = new case) | Global shortcut handler |

---

## Recommended Execution Order

```
Phase 1 (Weeks 1-3): Core testing experience — most user-facing value
  → 1.2 Test Case Editor
  → 1.1 Environments
  → 1.3 Test Chaining
  → 1.5 Pre/Post Scripts

Phase 2 (Weeks 3-4): OpenAPI enhancements — differentiator
  → 2.1 $ref Resolution
  → 2.2 Auth Scheme Detection
  → 2.3 Smart Test Generation

Phase 3 (Weeks 4-6): Reliability
  → 3.5 Persistent Queue
  → 3.1 Retry Logic
  → 3.2 Scheduled Runs
  → 3.3 Parallel Execution

Phase 4 (Weeks 6-7): CI/CD
  → 4.1 API Key Auth
  → 4.2 CLI Enhancements
  → 4.3 Webhook Triggers

Phase 5 (Weeks 7-8): Observability
  → 5.3 Request/Response Viewer
  → 5.4 Real-time Logs
  → 5.1 Detailed Reports

Phase 6 (Ongoing): UX polish
  → 6.1 Dark Mode
  → 6.2 Drag & Drop
  → 6.3 Bulk Operations
```

---

## Quick Wins (Immediate Implementation)

These can be implemented with minimal effort since the groundwork already exists:

1. **Environment table + API** — Simple migration + CRUD; `{{base_url}}` variable already supported by engine
2. **Test Case Editor form** — Replace "Add Case" button with a modal/drawer form
3. **$ref resolution** in the OpenAPI parser — recursive schema resolver
4. **Retry logic** — `retry_of` column already exists in `test_results`
5. **API key authentication middleware** — `api_key` field already on `users` table

---

## Technical Debt to Address

- [ ] Replace in-memory job queue with persistent storage (Redis/PostgreSQL)
- [ ] Add request/response size limits to HTTP executor
- [ ] Implement proper graceful shutdown for workers
- [ ] Add database connection pooling configuration
- [ ] Add structured error codes across all API endpoints
- [ ] Add request validation using the existing validator package
- [ ] Add comprehensive API documentation (Swagger/OpenAPI for the platform itself)
- [ ] Add integration tests for critical paths
- [ ] Add rate limiting middleware
