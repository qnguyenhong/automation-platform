package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/service"
	"github.com/qnguyenhong/automation-platform/pkg/errors"
)

// Ensure interfaces are implemented
var _ service.UserRepository = (*UserRepo)(nil)
var _ service.ProjectRepository = (*ProjectRepo)(nil)
var _ service.TestSuiteRepository = (*TestSuiteRepo)(nil)
var _ service.TestCaseRepository = (*TestCaseRepo)(nil)
var _ service.TestRunRepository = (*TestRunRepo)(nil)
var _ service.TestResultRepository = (*TestResultRepo)(nil)
var _ service.WorkerRepository = (*WorkerRepo)(nil)
var _ service.DashboardRepository = (*DashboardRepo)(nil)
var _ service.LoadMetricsRepository = (*LoadMetricsRepo)(nil)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash, name, role string, apiKey *string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name, role, api_key)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, email, password_hash, name, role, api_key, created_at, updated_at`,
		email, passwordHash, name, role, apiKey,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.APIKey, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, role, api_key, created_at, updated_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.APIKey, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, role, api_key, created_at, updated_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.APIKey, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &u, nil
}

func (r *UserRepo) GetByAPIKey(ctx context.Context, apiKey string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, role, api_key, created_at, updated_at FROM users WHERE api_key = $1`, apiKey,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.APIKey, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context) ([]model.User, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, email, password_hash, name, role, api_key, created_at, updated_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.APIKey, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) Update(ctx context.Context, id uuid.UUID, name, role string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`UPDATE users SET name = $2, role = $3, updated_at = now() WHERE id = $1
		 RETURNING id, email, password_hash, name, role, api_key, created_at, updated_at`,
		id, name, role,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.APIKey, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &u, nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

type ProjectRepo struct {
	pool *pgxpool.Pool
}

func NewProjectRepo(pool *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{pool: pool}
}

func (r *ProjectRepo) Create(ctx context.Context, name, description string, createdBy uuid.UUID) (*model.Project, error) {
	var p model.Project
	err := r.pool.QueryRow(ctx,
		`INSERT INTO projects (name, description, created_by) VALUES ($1, $2, $3)
		 RETURNING id, name, description, created_by, created_at, updated_at`,
		name, description, createdBy,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	var p model.Project
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, description, created_by, created_at, updated_at FROM projects WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &p, nil
}

func (r *ProjectRepo) List(ctx context.Context) ([]model.Project, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, description, created_by, created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepo) Update(ctx context.Context, id uuid.UUID, name, description string) (*model.Project, error) {
	var p model.Project
	err := r.pool.QueryRow(ctx,
		`UPDATE projects SET name = $2, description = $3, updated_at = now() WHERE id = $1
		 RETURNING id, name, description, created_by, created_at, updated_at`,
		id, name, description,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &p, nil
}

func (r *ProjectRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	return err
}

type TestSuiteRepo struct {
	pool *pgxpool.Pool
}

func NewTestSuiteRepo(pool *pgxpool.Pool) *TestSuiteRepo {
	return &TestSuiteRepo{pool: pool}
}

func (r *TestSuiteRepo) Create(ctx context.Context, projectID uuid.UUID, name, description string, testType string, scheduleCron *string, config []byte, tags []string, createdBy uuid.UUID) (*model.TestSuite, error) {
	var s model.TestSuite
	var createdByVal *uuid.UUID
	if createdBy != uuid.Nil {
		createdByVal = &createdBy
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO test_suites (project_id, name, description, test_type, schedule_cron, config, tags, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, project_id, name, description, test_type, schedule_cron, config, tags, created_by, created_at, updated_at`,
		projectID, name, description, testType, scheduleCron, string(config), tags, createdByVal,
	).Scan(&s.ID, &s.ProjectID, &s.Name, &s.Description, &s.TestType, &s.ScheduleCron, &s.Config, &s.Tags, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *TestSuiteRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.TestSuite, error) {
	var s model.TestSuite
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, description, test_type, schedule_cron, config, tags, created_by, created_at, updated_at
		 FROM test_suites WHERE id = $1`, id,
	).Scan(&s.ID, &s.ProjectID, &s.Name, &s.Description, &s.TestType, &s.ScheduleCron, &s.Config, &s.Tags, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &s, nil
}

func (r *TestSuiteRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.TestSuite, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, name, description, test_type, schedule_cron, config, tags, created_by, created_at, updated_at
		 FROM test_suites WHERE project_id = $1 ORDER BY created_at DESC`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suites []model.TestSuite
	for rows.Next() {
		var s model.TestSuite
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.Description, &s.TestType, &s.ScheduleCron, &s.Config, &s.Tags, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		suites = append(suites, s)
	}
	return suites, nil
}

func (r *TestSuiteRepo) Update(ctx context.Context, id uuid.UUID, name, description string, testType string, scheduleCron *string, config []byte, tags []string) (*model.TestSuite, error) {
	var s model.TestSuite
	err := r.pool.QueryRow(ctx,
		`UPDATE test_suites SET name = $2, description = $3, test_type = $4, schedule_cron = $5, config = $6, tags = $7, updated_at = now()
		 WHERE id = $1
		 RETURNING id, project_id, name, description, test_type, schedule_cron, config, tags, created_by, created_at, updated_at`,
		id, name, description, testType, scheduleCron, string(config), tags,
	).Scan(&s.ID, &s.ProjectID, &s.Name, &s.Description, &s.TestType, &s.ScheduleCron, &s.Config, &s.Tags, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &s, nil
}

func (r *TestSuiteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM test_suites WHERE id = $1`, id)
	return err
}

type TestCaseRepo struct {
	pool *pgxpool.Pool
}

func NewTestCaseRepo(pool *pgxpool.Pool) *TestCaseRepo {
	return &TestCaseRepo{pool: pool}
}

func (r *TestCaseRepo) Create(ctx context.Context, suiteID uuid.UUID, name, description string, config []byte, tags []string, sortOrder int, enabled bool) (*model.TestCase, error) {
	var tc model.TestCase
	err := r.pool.QueryRow(ctx,
		`INSERT INTO test_cases (suite_id, name, description, config, tags, sort_order, enabled)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, suite_id, name, description, config, tags, sort_order, enabled, created_at, updated_at`,
		suiteID, name, description, string(config), tags, sortOrder, enabled,
	).Scan(&tc.ID, &tc.SuiteID, &tc.Name, &tc.Description, &tc.Config, &tc.Tags, &tc.SortOrder, &tc.Enabled, &tc.CreatedAt, &tc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *TestCaseRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.TestCase, error) {
	var tc model.TestCase
	err := r.pool.QueryRow(ctx,
		`SELECT id, suite_id, name, description, config, tags, sort_order, enabled, created_at, updated_at
		 FROM test_cases WHERE id = $1`, id,
	).Scan(&tc.ID, &tc.SuiteID, &tc.Name, &tc.Description, &tc.Config, &tc.Tags, &tc.SortOrder, &tc.Enabled, &tc.CreatedAt, &tc.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &tc, nil
}

func (r *TestCaseRepo) ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]model.TestCase, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, suite_id, name, description, config, tags, sort_order, enabled, created_at, updated_at
		 FROM test_cases WHERE suite_id = $1 ORDER BY sort_order`, suiteID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cases []model.TestCase
	for rows.Next() {
		var tc model.TestCase
		if err := rows.Scan(&tc.ID, &tc.SuiteID, &tc.Name, &tc.Description, &tc.Config, &tc.Tags, &tc.SortOrder, &tc.Enabled, &tc.CreatedAt, &tc.UpdatedAt); err != nil {
			return nil, err
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func (r *TestCaseRepo) Update(ctx context.Context, id uuid.UUID, name, description string, config []byte, tags []string, sortOrder int, enabled bool) (*model.TestCase, error) {
	var tc model.TestCase
	err := r.pool.QueryRow(ctx,
		`UPDATE test_cases SET name = $2, description = $3, config = $4, tags = $5, sort_order = $6, enabled = $7, updated_at = now()
		 WHERE id = $1
		 RETURNING id, suite_id, name, description, config, tags, sort_order, enabled, created_at, updated_at`,
		id, name, description, string(config), tags, sortOrder, enabled,
	).Scan(&tc.ID, &tc.SuiteID, &tc.Name, &tc.Description, &tc.Config, &tc.Tags, &tc.SortOrder, &tc.Enabled, &tc.CreatedAt, &tc.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &tc, nil
}

func (r *TestCaseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM test_cases WHERE id = $1`, id)
	return err
}

type TestRunRepo struct {
	pool *pgxpool.Pool
}

func NewTestRunRepo(pool *pgxpool.Pool) *TestRunRepo {
	return &TestRunRepo{pool: pool}
}

func (r *TestRunRepo) Create(ctx context.Context, suiteID, projectID uuid.UUID, status, trigger string, triggeredBy *uuid.UUID, totalCases int, metadata []byte) (*model.TestRun, error) {
	var tr model.TestRun
	err := r.pool.QueryRow(ctx,
		`INSERT INTO test_runs (suite_id, project_id, status, trigger, triggered_by, total_cases, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at`,
		suiteID, projectID, status, trigger, triggeredBy, totalCases, string(metadata),
	).Scan(&tr.ID, &tr.SuiteID, &tr.ProjectID, &tr.Status, &tr.Trigger, &tr.TriggeredBy, &tr.TotalCases, &tr.Passed, &tr.Failed, &tr.Skipped, &tr.Errored, &tr.StartedAt, &tr.FinishedAt, &tr.DurationMS, &tr.Metadata, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (r *TestRunRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.TestRun, error) {
	var tr model.TestRun
	err := r.pool.QueryRow(ctx,
		`SELECT id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at
		 FROM test_runs WHERE id = $1`, id,
	).Scan(&tr.ID, &tr.SuiteID, &tr.ProjectID, &tr.Status, &tr.Trigger, &tr.TriggeredBy, &tr.TotalCases, &tr.Passed, &tr.Failed, &tr.Skipped, &tr.Errored, &tr.StartedAt, &tr.FinishedAt, &tr.DurationMS, &tr.Metadata, &tr.CreatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &tr, nil
}

func (r *TestRunRepo) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]model.TestRun, int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM test_runs WHERE project_id = $1`, projectID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at
		 FROM test_runs WHERE project_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		projectID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var runs []model.TestRun
	for rows.Next() {
		var tr model.TestRun
		if err := rows.Scan(&tr.ID, &tr.SuiteID, &tr.ProjectID, &tr.Status, &tr.Trigger, &tr.TriggeredBy, &tr.TotalCases, &tr.Passed, &tr.Failed, &tr.Skipped, &tr.Errored, &tr.StartedAt, &tr.FinishedAt, &tr.DurationMS, &tr.Metadata, &tr.CreatedAt); err != nil {
			return nil, 0, err
		}
		runs = append(runs, tr)
	}
	return runs, total, nil
}

func (r *TestRunRepo) ListAll(ctx context.Context, limit, offset int) ([]model.TestRun, int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM test_runs`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at
		 FROM test_runs ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var runs []model.TestRun
	for rows.Next() {
		var tr model.TestRun
		if err := rows.Scan(&tr.ID, &tr.SuiteID, &tr.ProjectID, &tr.Status, &tr.Trigger, &tr.TriggeredBy, &tr.TotalCases, &tr.Passed, &tr.Failed, &tr.Skipped, &tr.Errored, &tr.StartedAt, &tr.FinishedAt, &tr.DurationMS, &tr.Metadata, &tr.CreatedAt); err != nil {
			return nil, 0, err
		}
		runs = append(runs, tr)
	}
	return runs, total, nil
}

func (r *TestRunRepo) ListBySuite(ctx context.Context, suiteID uuid.UUID, limit, offset int) ([]model.TestRun, int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM test_runs WHERE suite_id = $1`, suiteID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at
		 FROM test_runs WHERE suite_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		suiteID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var runs []model.TestRun
	for rows.Next() {
		var tr model.TestRun
		if err := rows.Scan(&tr.ID, &tr.SuiteID, &tr.ProjectID, &tr.Status, &tr.Trigger, &tr.TriggeredBy, &tr.TotalCases, &tr.Passed, &tr.Failed, &tr.Skipped, &tr.Errored, &tr.StartedAt, &tr.FinishedAt, &tr.DurationMS, &tr.Metadata, &tr.CreatedAt); err != nil {
			return nil, 0, err
		}
		runs = append(runs, tr)
	}
	return runs, total, nil
}

func (r *TestRunRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, startedAt, finishedAt *time.Time, durationMS *int64) (*model.TestRun, error) {
	var tr model.TestRun
	err := r.pool.QueryRow(ctx,
		`UPDATE test_runs SET status = $2, started_at = COALESCE($3, started_at), finished_at = $4, duration_ms = $5
		 WHERE id = $1
		 RETURNING id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at`,
		id, status, startedAt, finishedAt, durationMS,
	).Scan(&tr.ID, &tr.SuiteID, &tr.ProjectID, &tr.Status, &tr.Trigger, &tr.TriggeredBy, &tr.TotalCases, &tr.Passed, &tr.Failed, &tr.Skipped, &tr.Errored, &tr.StartedAt, &tr.FinishedAt, &tr.DurationMS, &tr.Metadata, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (r *TestRunRepo) UpdateCounts(ctx context.Context, id uuid.UUID, passed, failed, skipped, errored int) (*model.TestRun, error) {
	var tr model.TestRun
	err := r.pool.QueryRow(ctx,
		`UPDATE test_runs SET passed = $2, failed = $3, skipped = $4, errored = $5
		 WHERE id = $1
		 RETURNING id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at`,
		id, passed, failed, skipped, errored,
	).Scan(&tr.ID, &tr.SuiteID, &tr.ProjectID, &tr.Status, &tr.Trigger, &tr.TriggeredBy, &tr.TotalCases, &tr.Passed, &tr.Failed, &tr.Skipped, &tr.Errored, &tr.StartedAt, &tr.FinishedAt, &tr.DurationMS, &tr.Metadata, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

type TestResultRepo struct {
	pool *pgxpool.Pool
}

func NewTestResultRepo(pool *pgxpool.Pool) *TestResultRepo {
	return &TestResultRepo{pool: pool}
}

func (r *TestResultRepo) Create(ctx context.Context, runID, caseID uuid.UUID, workerID *uuid.UUID, status string) (*model.TestResult, error) {
	var tr model.TestResult
	err := r.pool.QueryRow(ctx,
		`INSERT INTO test_results (run_id, case_id, worker_id, status)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, run_id, case_id, worker_id, status, created_at`,
		runID, caseID, workerID, status,
	).Scan(&tr.ID, &tr.RunID, &tr.CaseID, &tr.WorkerID, &tr.Status, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (r *TestResultRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.TestResult, error) {
	var tr model.TestResult
	err := r.pool.QueryRow(ctx,
		`SELECT id, run_id, case_id, worker_id, status, error_message, assertions, request_data, response_data, artifacts, metrics, stdout, stderr, duration_ms, started_at, finished_at, retry_of, created_at
		 FROM test_results WHERE id = $1`, id,
	).Scan(&tr.ID, &tr.RunID, &tr.CaseID, &tr.WorkerID, &tr.Status, &tr.ErrorMessage, &tr.Assertions, &tr.RequestData, &tr.ResponseData, &tr.Artifacts, &tr.Metrics, &tr.Stdout, &tr.Stderr, &tr.DurationMS, &tr.StartedAt, &tr.FinishedAt, &tr.RetryOf, &tr.CreatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &tr, nil
}

func (r *TestResultRepo) ListByRun(ctx context.Context, runID uuid.UUID) ([]model.TestResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, run_id, case_id, worker_id, status, error_message, assertions, request_data, response_data, artifacts, metrics, stdout, stderr, duration_ms, started_at, finished_at, retry_of, created_at
		 FROM test_results WHERE run_id = $1 ORDER BY created_at`, runID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.TestResult
	for rows.Next() {
		var tr model.TestResult
		if err := rows.Scan(&tr.ID, &tr.RunID, &tr.CaseID, &tr.WorkerID, &tr.Status, &tr.ErrorMessage, &tr.Assertions, &tr.RequestData, &tr.ResponseData, &tr.Artifacts, &tr.Metrics, &tr.Stdout, &tr.Stderr, &tr.DurationMS, &tr.StartedAt, &tr.FinishedAt, &tr.RetryOf, &tr.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, tr)
	}
	return results, nil
}

func (r *TestResultRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, startedAt, finishedAt *time.Time, durationMS *int64) (*model.TestResult, error) {
	var tr model.TestResult
	err := r.pool.QueryRow(ctx,
		`UPDATE test_results SET status = $2, started_at = $3, finished_at = $4, duration_ms = $5
		 WHERE id = $1
		 RETURNING id, run_id, case_id, worker_id, status, created_at`,
		id, status, startedAt, finishedAt, durationMS,
	).Scan(&tr.ID, &tr.RunID, &tr.CaseID, &tr.WorkerID, &tr.Status, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (r *TestResultRepo) UpdateData(ctx context.Context, id uuid.UUID, status string, errorMessage *string, assertions, requestData, responseData, artifacts, metrics []byte, stdout, stderr *string, durationMS *int64, finishedAt *time.Time) (*model.TestResult, error) {
	var tr model.TestResult
	err := r.pool.QueryRow(ctx,
		`UPDATE test_results SET status = $2, error_message = $3, assertions = $4, request_data = $5, response_data = $6, artifacts = $7, metrics = $8, stdout = $9, stderr = $10, duration_ms = $11, finished_at = $12
		 WHERE id = $1
		 RETURNING id, run_id, case_id, worker_id, status, created_at`,
		id, status, errorMessage, string(assertions), string(requestData), string(responseData), string(artifacts), string(metrics), stdout, stderr, durationMS, finishedAt,
	).Scan(&tr.ID, &tr.RunID, &tr.CaseID, &tr.WorkerID, &tr.Status, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

type WorkerRepo struct {
	pool *pgxpool.Pool
}

func NewWorkerRepo(pool *pgxpool.Pool) *WorkerRepo {
	return &WorkerRepo{pool: pool}
}

func (r *WorkerRepo) Create(ctx context.Context, name, hostname, ipAddress string, executorTypes []string, maxConcurrent int, version *string, labels []byte) (*model.Worker, error) {
	var w model.Worker
	err := r.pool.QueryRow(ctx,
		`INSERT INTO workers (name, hostname, ip_address, executor_types, max_concurrent, version, labels)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, name, hostname, ip_address, executor_types, status, current_job_id, max_concurrent, version, labels, last_heartbeat, registered_at, updated_at`,
		name, hostname, ipAddress, executorTypes, maxConcurrent, version, string(labels),
	).Scan(&w.ID, &w.Name, &w.Hostname, &w.IPAddress, &w.ExecutorTypes, &w.Status, &w.CurrentJobID, &w.MaxConcurrent, &w.Version, &w.Labels, &w.LastHeartbeat, &w.RegisteredAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WorkerRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Worker, error) {
	var w model.Worker
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, hostname, ip_address, executor_types, status, current_job_id, max_concurrent, version, labels, last_heartbeat, registered_at, updated_at
		 FROM workers WHERE id = $1`, id,
	).Scan(&w.ID, &w.Name, &w.Hostname, &w.IPAddress, &w.ExecutorTypes, &w.Status, &w.CurrentJobID, &w.MaxConcurrent, &w.Version, &w.Labels, &w.LastHeartbeat, &w.RegisteredAt, &w.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &w, nil
}

func (r *WorkerRepo) List(ctx context.Context) ([]model.Worker, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, hostname, ip_address, executor_types, status, current_job_id, max_concurrent, version, labels, last_heartbeat, registered_at, updated_at
		 FROM workers ORDER BY registered_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workers []model.Worker
	for rows.Next() {
		var w model.Worker
		if err := rows.Scan(&w.ID, &w.Name, &w.Hostname, &w.IPAddress, &w.ExecutorTypes, &w.Status, &w.CurrentJobID, &w.MaxConcurrent, &w.Version, &w.Labels, &w.LastHeartbeat, &w.RegisteredAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}

func (r *WorkerRepo) ListByStatus(ctx context.Context, status string) ([]model.Worker, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, hostname, ip_address, executor_types, status, current_job_id, max_concurrent, version, labels, last_heartbeat, registered_at, updated_at
		 FROM workers WHERE status = $1 ORDER BY registered_at DESC`, status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workers []model.Worker
	for rows.Next() {
		var w model.Worker
		if err := rows.Scan(&w.ID, &w.Name, &w.Hostname, &w.IPAddress, &w.ExecutorTypes, &w.Status, &w.CurrentJobID, &w.MaxConcurrent, &w.Version, &w.Labels, &w.LastHeartbeat, &w.RegisteredAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}

func (r *WorkerRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*model.Worker, error) {
	var w model.Worker
	err := r.pool.QueryRow(ctx,
		`UPDATE workers SET status = $2, last_heartbeat = now(), updated_at = now()
		 WHERE id = $1
		 RETURNING id, name, hostname, ip_address, executor_types, status, current_job_id, max_concurrent, version, labels, last_heartbeat, registered_at, updated_at`,
		id, status,
	).Scan(&w.ID, &w.Name, &w.Hostname, &w.IPAddress, &w.ExecutorTypes, &w.Status, &w.CurrentJobID, &w.MaxConcurrent, &w.Version, &w.Labels, &w.LastHeartbeat, &w.RegisteredAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WorkerRepo) UpdateHeartbeat(ctx context.Context, id uuid.UUID, status string, currentJobID *uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE workers SET status = $2, current_job_id = $3, last_heartbeat = now(), updated_at = now() WHERE id = $1`,
		id, status, currentJobID,
	)
	return err
}

func (r *WorkerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM workers WHERE id = $1`, id)
	return err
}

type DashboardRepo struct {
	pool *pgxpool.Pool
}

func NewDashboardRepo(pool *pgxpool.Pool) *DashboardRepo {
	return &DashboardRepo{pool: pool}
}

func (r *DashboardRepo) GetSummary(ctx context.Context, projectID *uuid.UUID, since time.Time) (*service.DashboardSummary, error) {
	var s service.DashboardSummary
	if projectID != nil {
		err := r.pool.QueryRow(ctx,
			`SELECT
				COALESCE((SELECT COUNT(*) FROM test_runs WHERE project_id = $1 AND created_at >= $2), 0),
				COALESCE((SELECT COUNT(*) FROM test_runs WHERE project_id = $1 AND status = 'passed' AND created_at >= $2), 0),
				COALESCE((SELECT COUNT(*) FROM test_runs WHERE project_id = $1 AND status = 'failed' AND created_at >= $2), 0),
				COALESCE((SELECT COUNT(*) FROM workers WHERE status = 'online'), 0)`,
			*projectID, since,
		).Scan(&s.TotalRuns, &s.TotalPassed, &s.TotalFailed, &s.ActiveWorkers)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.pool.QueryRow(ctx,
			`SELECT
				COALESCE((SELECT COUNT(*) FROM test_runs WHERE created_at >= $1), 0),
				COALESCE((SELECT COUNT(*) FROM test_runs WHERE status = 'passed' AND created_at >= $1), 0),
				COALESCE((SELECT COUNT(*) FROM test_runs WHERE status = 'failed' AND created_at >= $1), 0),
				COALESCE((SELECT COUNT(*) FROM workers WHERE status = 'online'), 0)`,
			since,
		).Scan(&s.TotalRuns, &s.TotalPassed, &s.TotalFailed, &s.ActiveWorkers)
		if err != nil {
			return nil, err
		}
	}

	if s.TotalRuns > 0 {
		s.PassRate = float64(s.TotalPassed) / float64(s.TotalRuns) * 100
	}

	return &s, nil
}

func (r *DashboardRepo) GetRecentRuns(ctx context.Context, projectID *uuid.UUID, limit int) ([]model.TestRun, error) {
	var rows_query string
	var args []any

	if projectID != nil {
		rows_query = `SELECT id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at
		 FROM test_runs WHERE project_id = $1 ORDER BY created_at DESC LIMIT $2`
		args = []any{*projectID, limit}
	} else {
		rows_query = `SELECT id, suite_id, project_id, status, trigger, triggered_by, total_cases, passed, failed, skipped, errored, started_at, finished_at, duration_ms, metadata, created_at
		 FROM test_runs ORDER BY created_at DESC LIMIT $1`
		args = []any{limit}
	}

	rows, err := r.pool.Query(ctx, rows_query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []model.TestRun
	for rows.Next() {
		var tr model.TestRun
		if err := rows.Scan(&tr.ID, &tr.SuiteID, &tr.ProjectID, &tr.Status, &tr.Trigger, &tr.TriggeredBy, &tr.TotalCases, &tr.Passed, &tr.Failed, &tr.Skipped, &tr.Errored, &tr.StartedAt, &tr.FinishedAt, &tr.DurationMS, &tr.Metadata, &tr.CreatedAt); err != nil {
			return nil, err
		}
		runs = append(runs, tr)
	}
	return runs, nil
}

func (r *DashboardRepo) GetFlakyTests(ctx context.Context, projectID *uuid.UUID, since time.Time, limit int) ([]service.FlakyTest, error) {
	var flaky []service.FlakyTest

	if projectID != nil {
		rows, err := r.pool.Query(ctx,
			`SELECT tc.id, tc.name, COUNT(*), COUNT(*) FILTER (WHERE tr.status = 'passed'), COUNT(*) FILTER (WHERE tr.status = 'failed')
			 FROM test_results tr
			 JOIN test_cases tc ON tr.case_id = tc.id
			 JOIN test_runs r ON tr.run_id = r.id
			 WHERE r.project_id = $1 AND r.created_at >= $2
			 GROUP BY tc.id, tc.name
			 HAVING COUNT(*) FILTER (WHERE tr.status = 'passed') > 0 AND COUNT(*) FILTER (WHERE tr.status = 'failed') > 0
			 ORDER BY COUNT(*) FILTER (WHERE tr.status = 'failed') DESC
			 LIMIT $3`,
			*projectID, since, limit,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var f service.FlakyTest
			if err := rows.Scan(&f.CaseID, &f.CaseName, &f.Total, &f.Passed, &f.Failed); err != nil {
				return nil, err
			}
			flaky = append(flaky, f)
		}
	} else {
		rows, err := r.pool.Query(ctx,
			`SELECT tc.id, tc.name, COUNT(*), COUNT(*) FILTER (WHERE tr.status = 'passed'), COUNT(*) FILTER (WHERE tr.status = 'failed')
			 FROM test_results tr
			 JOIN test_cases tc ON tr.case_id = tc.id
			 JOIN test_runs r ON tr.run_id = r.id
			 WHERE r.created_at >= $1
			 GROUP BY tc.id, tc.name
			 HAVING COUNT(*) FILTER (WHERE tr.status = 'passed') > 0 AND COUNT(*) FILTER (WHERE tr.status = 'failed') > 0
			 ORDER BY COUNT(*) FILTER (WHERE tr.status = 'failed') DESC
			 LIMIT $2`,
			since, limit,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var f service.FlakyTest
			if err := rows.Scan(&f.CaseID, &f.CaseName, &f.Total, &f.Passed, &f.Failed); err != nil {
				return nil, err
			}
			flaky = append(flaky, f)
		}
	}

	return flaky, nil
}

func (r *DashboardRepo) GetTrends(ctx context.Context, projectID *uuid.UUID, from, to time.Time) ([]service.TrendPoint, error) {
	var trends []service.TrendPoint

	if projectID != nil {
		rows, err := r.pool.Query(ctx,
			`SELECT
				(created_at::date)::text as date,
				COUNT(*) FILTER (WHERE status = 'passed'),
				COUNT(*) FILTER (WHERE status = 'failed'),
				COALESCE(AVG(duration_ms), 0)::bigint
			 FROM test_runs
			 WHERE project_id = $1 AND created_at >= $2 AND created_at <= $3
			 GROUP BY created_at::date
			 ORDER BY date`,
			*projectID, from, to,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var t service.TrendPoint
			if err := rows.Scan(&t.Date, &t.Passed, &t.Failed, &t.AvgDuration); err != nil {
				return nil, err
			}
			trends = append(trends, t)
		}
	} else {
		rows, err := r.pool.Query(ctx,
			`SELECT
				(created_at::date)::text as date,
				COUNT(*) FILTER (WHERE status = 'passed'),
				COUNT(*) FILTER (WHERE status = 'failed'),
				COALESCE(AVG(duration_ms), 0)::bigint
			 FROM test_runs
			 WHERE created_at >= $1 AND created_at <= $2
			 GROUP BY created_at::date
			 ORDER BY date`,
			from, to,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var t service.TrendPoint
			if err := rows.Scan(&t.Date, &t.Passed, &t.Failed, &t.AvgDuration); err != nil {
				return nil, err
			}
			trends = append(trends, t)
		}
	}

	return trends, nil
}

// OpenAPIRepo implements service.OpenAPIRepository.
type OpenAPIRepo struct {
	pool *pgxpool.Pool
}

var _ service.OpenAPIRepository = (*OpenAPIRepo)(nil)

func NewOpenAPIRepo(pool *pgxpool.Pool) *OpenAPIRepo {
	return &OpenAPIRepo{pool: pool}
}

func (r *OpenAPIRepo) ListImports(ctx context.Context, projectID uuid.UUID) ([]service.OpenAPIImport, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, filename, spec_url, version, imported_by, imported_at
		 FROM openapi_imports WHERE project_id = $1 ORDER BY imported_at DESC`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var imports []service.OpenAPIImport
	for rows.Next() {
		var imp service.OpenAPIImport
		if err := rows.Scan(&imp.ID, &imp.ProjectID, &imp.Filename, &imp.SpecURL, &imp.Version, &imp.ImportedBy, &imp.ImportedAt); err != nil {
			return nil, err
		}
		imports = append(imports, imp)
	}
	return imports, nil
}

type LoadMetricsRepo struct {
	pool *pgxpool.Pool
}

func NewLoadMetricsRepo(pool *pgxpool.Pool) *LoadMetricsRepo {
	return &LoadMetricsRepo{pool: pool}
}

func (r *LoadMetricsRepo) Create(ctx context.Context, m *model.LoadMetricsDb) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO load_metrics (
			result_id, run_id, total_requests, success_count, error_count, error_rate, throughput_rps,
			min_latency_ms, max_latency_ms, avg_latency_ms, p50_latency_ms, p90_latency_ms, p95_latency_ms, p99_latency_ms,
			total_bytes, time_series, status_codes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, created_at`,
		m.ResultID, m.RunID, m.TotalRequests, m.SuccessCount, m.ErrorCount, m.ErrorRate, m.ThroughputRPS,
		m.MinLatencyMS, m.MaxLatencyMS, m.AvgLatencyMS, m.P50LatencyMS, m.P90LatencyMS, m.P95LatencyMS, m.P99LatencyMS,
		m.TotalBytes, string(m.TimeSeries), string(m.StatusCodes),
	).Scan(&m.ID, &m.CreatedAt)
	return err
}

func (r *LoadMetricsRepo) GetByRun(ctx context.Context, runID uuid.UUID) (*model.LoadMetricsDb, error) {
	var m model.LoadMetricsDb
	err := r.pool.QueryRow(ctx,
		`SELECT id, result_id, run_id, total_requests, success_count, error_count, error_rate, throughput_rps,
			min_latency_ms, max_latency_ms, avg_latency_ms, p50_latency_ms, p90_latency_ms, p95_latency_ms, p99_latency_ms,
			total_bytes, time_series, status_codes, created_at
		 FROM load_metrics WHERE run_id = $1`, runID,
	).Scan(
		&m.ID, &m.ResultID, &m.RunID, &m.TotalRequests, &m.SuccessCount, &m.ErrorCount, &m.ErrorRate, &m.ThroughputRPS,
		&m.MinLatencyMS, &m.MaxLatencyMS, &m.AvgLatencyMS, &m.P50LatencyMS, &m.P90LatencyMS, &m.P95LatencyMS, &m.P99LatencyMS,
		&m.TotalBytes, &m.TimeSeries, &m.StatusCodes, &m.CreatedAt,
	)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &m, nil
}

func (r *OpenAPIRepo) CreateImport(ctx context.Context, projectID uuid.UUID, filename, specURL string, content []byte, version string, importedBy uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO openapi_imports (project_id, filename, spec_url, content, version, imported_by)
		 VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6)`,
		projectID, filename, specURL, string(content), version, importedBy,
	)
	return err
}

func (r *LoadMetricsRepo) GetByResult(ctx context.Context, resultID uuid.UUID) (*model.LoadMetricsDb, error) {
	var m model.LoadMetricsDb
	err := r.pool.QueryRow(ctx,
		`SELECT id, result_id, run_id, total_requests, success_count, error_count, error_rate, throughput_rps,
			min_latency_ms, max_latency_ms, avg_latency_ms, p50_latency_ms, p90_latency_ms, p95_latency_ms, p99_latency_ms,
			total_bytes, time_series, status_codes, created_at
		 FROM load_metrics WHERE result_id = $1`, resultID,
	).Scan(
		&m.ID, &m.ResultID, &m.RunID, &m.TotalRequests, &m.SuccessCount, &m.ErrorCount, &m.ErrorRate, &m.ThroughputRPS,
		&m.MinLatencyMS, &m.MaxLatencyMS, &m.AvgLatencyMS, &m.P50LatencyMS, &m.P90LatencyMS, &m.P95LatencyMS, &m.P99LatencyMS,
		&m.TotalBytes, &m.TimeSeries, &m.StatusCodes, &m.CreatedAt,
	)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &m, nil
}
