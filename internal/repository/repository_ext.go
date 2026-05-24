package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/service"
	"github.com/qnguyenhong/automation-platform/pkg/errors"
)

// Ensure interfaces are implemented
var _ service.EnvironmentRepository = (*EnvironmentRepo)(nil)
var _ service.DatasetRepository = (*DatasetRepo)(nil)
var _ service.WebhookRepository = (*WebhookRepo)(nil)
var _ service.ScheduleRepository = (*ScheduleRepo)(nil)

// EnvironmentRepo
type EnvironmentRepo struct {
	pool *pgxpool.Pool
}

func NewEnvironmentRepo(pool *pgxpool.Pool) *EnvironmentRepo {
	return &EnvironmentRepo{pool: pool}
}

func (r *EnvironmentRepo) Create(ctx context.Context, projectID uuid.UUID, name, baseURL string, variables, headers []byte, isDefault bool) (*model.Environment, error) {
	var env model.Environment
	err := r.pool.QueryRow(ctx,
		`INSERT INTO environments (project_id, name, base_url, variables, headers, is_default)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, project_id, name, base_url, variables, headers, is_default, created_at, updated_at`,
		projectID, name, baseURL, string(variables), string(headers), isDefault,
	).Scan(&env.ID, &env.ProjectID, &env.Name, &env.BaseURL, &env.Variables, &env.Headers, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (r *EnvironmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Environment, error) {
	var env model.Environment
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, base_url, variables, headers, is_default, created_at, updated_at
		 FROM environments WHERE id = $1`, id,
	).Scan(&env.ID, &env.ProjectID, &env.Name, &env.BaseURL, &env.Variables, &env.Headers, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &env, nil
}

func (r *EnvironmentRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, name, base_url, variables, headers, is_default, created_at, updated_at
		 FROM environments WHERE project_id = $1 ORDER BY name`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []model.Environment
	for rows.Next() {
		var env model.Environment
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &env.BaseURL, &env.Variables, &env.Headers, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}
	return envs, nil
}

func (r *EnvironmentRepo) Update(ctx context.Context, id uuid.UUID, name, baseURL string, variables, headers []byte, isDefault bool) (*model.Environment, error) {
	var env model.Environment
	err := r.pool.QueryRow(ctx,
		`UPDATE environments SET name = $2, base_url = $3, variables = $4, headers = $5, is_default = $6, updated_at = now()
		 WHERE id = $1
		 RETURNING id, project_id, name, base_url, variables, headers, is_default, created_at, updated_at`,
		id, name, baseURL, string(variables), string(headers), isDefault,
	).Scan(&env.ID, &env.ProjectID, &env.Name, &env.BaseURL, &env.Variables, &env.Headers, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &env, nil
}

func (r *EnvironmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM environments WHERE id = $1`, id)
	return err
}

func (r *EnvironmentRepo) GetDefault(ctx context.Context, projectID uuid.UUID) (*model.Environment, error) {
	var env model.Environment
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, base_url, variables, headers, is_default, created_at, updated_at
		 FROM environments WHERE project_id = $1 AND is_default = true LIMIT 1`, projectID,
	).Scan(&env.ID, &env.ProjectID, &env.Name, &env.BaseURL, &env.Variables, &env.Headers, &env.IsDefault, &env.CreatedAt, &env.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &env, nil
}

// DatasetRepo
type DatasetRepo struct {
	pool *pgxpool.Pool
}

func NewDatasetRepo(pool *pgxpool.Pool) *DatasetRepo {
	return &DatasetRepo{pool: pool}
}

func (r *DatasetRepo) Create(ctx context.Context, projectID uuid.UUID, name, description, format string, data []byte) (*model.Dataset, error) {
	var ds model.Dataset
	err := r.pool.QueryRow(ctx,
		`INSERT INTO datasets (project_id, name, description, format, data)
		 VALUES ($1, $2, NULLIF($3, ''), $4, $5)
		 RETURNING id, project_id, name, description, format, data, created_at, updated_at`,
		projectID, name, description, format, string(data),
	).Scan(&ds.ID, &ds.ProjectID, &ds.Name, &ds.Description, &ds.Format, &ds.Data, &ds.CreatedAt, &ds.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &ds, nil
}

func (r *DatasetRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Dataset, error) {
	var ds model.Dataset
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, description, format, data, created_at, updated_at
		 FROM datasets WHERE id = $1`, id,
	).Scan(&ds.ID, &ds.ProjectID, &ds.Name, &ds.Description, &ds.Format, &ds.Data, &ds.CreatedAt, &ds.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &ds, nil
}

func (r *DatasetRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Dataset, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, name, description, format, data, created_at, updated_at
		 FROM datasets WHERE project_id = $1 ORDER BY created_at DESC`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datasets []model.Dataset
	for rows.Next() {
		var ds model.Dataset
		if err := rows.Scan(&ds.ID, &ds.ProjectID, &ds.Name, &ds.Description, &ds.Format, &ds.Data, &ds.CreatedAt, &ds.UpdatedAt); err != nil {
			return nil, err
		}
		datasets = append(datasets, ds)
	}
	return datasets, nil
}

func (r *DatasetRepo) Update(ctx context.Context, id uuid.UUID, name, description, format string, data []byte) (*model.Dataset, error) {
	var ds model.Dataset
	err := r.pool.QueryRow(ctx,
		`UPDATE datasets SET name = $2, description = NULLIF($3, ''), format = $4, data = $5, updated_at = now()
		 WHERE id = $1
		 RETURNING id, project_id, name, description, format, data, created_at, updated_at`,
		id, name, description, format, string(data),
	).Scan(&ds.ID, &ds.ProjectID, &ds.Name, &ds.Description, &ds.Format, &ds.Data, &ds.CreatedAt, &ds.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &ds, nil
}

func (r *DatasetRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM datasets WHERE id = $1`, id)
	return err
}

// WebhookRepo
type WebhookRepo struct {
	pool *pgxpool.Pool
}

func NewWebhookRepo(pool *pgxpool.Pool) *WebhookRepo {
	return &WebhookRepo{pool: pool}
}

func generateSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (r *WebhookRepo) Create(ctx context.Context, projectID uuid.UUID, name string, suiteID *uuid.UUID, enabled bool) (*model.Webhook, error) {
	var wh model.Webhook
	secret := generateSecret()
	err := r.pool.QueryRow(ctx,
		`INSERT INTO webhooks (project_id, name, secret, suite_id, enabled)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, project_id, name, secret, suite_id, enabled, last_triggered_at, created_at, updated_at`,
		projectID, name, secret, suiteID, enabled,
	).Scan(&wh.ID, &wh.ProjectID, &wh.Name, &wh.Secret, &wh.SuiteID, &wh.Enabled, &wh.LastTriggeredAt, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

func (r *WebhookRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Webhook, error) {
	var wh model.Webhook
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, secret, suite_id, enabled, last_triggered_at, created_at, updated_at
		 FROM webhooks WHERE id = $1`, id,
	).Scan(&wh.ID, &wh.ProjectID, &wh.Name, &wh.Secret, &wh.SuiteID, &wh.Enabled, &wh.LastTriggeredAt, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &wh, nil
}

func (r *WebhookRepo) GetBySecret(ctx context.Context, secret string) (*model.Webhook, error) {
	var wh model.Webhook
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, secret, suite_id, enabled, last_triggered_at, created_at, updated_at
		 FROM webhooks WHERE secret = $1 AND enabled = true`, secret,
	).Scan(&wh.ID, &wh.ProjectID, &wh.Name, &wh.Secret, &wh.SuiteID, &wh.Enabled, &wh.LastTriggeredAt, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &wh, nil
}

func (r *WebhookRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Webhook, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, name, secret, suite_id, enabled, last_triggered_at, created_at, updated_at
		 FROM webhooks WHERE project_id = $1 ORDER BY created_at DESC`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []model.Webhook
	for rows.Next() {
		var wh model.Webhook
		if err := rows.Scan(&wh.ID, &wh.ProjectID, &wh.Name, &wh.Secret, &wh.SuiteID, &wh.Enabled, &wh.LastTriggeredAt, &wh.CreatedAt, &wh.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, wh)
	}
	return webhooks, nil
}

func (r *WebhookRepo) Update(ctx context.Context, id uuid.UUID, name string, suiteID *uuid.UUID, enabled bool) (*model.Webhook, error) {
	var wh model.Webhook
	err := r.pool.QueryRow(ctx,
		`UPDATE webhooks SET name = $2, suite_id = $3, enabled = $4, updated_at = now()
		 WHERE id = $1
		 RETURNING id, project_id, name, secret, suite_id, enabled, last_triggered_at, created_at, updated_at`,
		id, name, suiteID, enabled,
	).Scan(&wh.ID, &wh.ProjectID, &wh.Name, &wh.Secret, &wh.SuiteID, &wh.Enabled, &wh.LastTriggeredAt, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &wh, nil
}

func (r *WebhookRepo) UpdateLastTriggered(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE webhooks SET last_triggered_at = now() WHERE id = $1`, id)
	return err
}

func (r *WebhookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM webhooks WHERE id = $1`, id)
	return err
}

// ScheduleRepo
type ScheduleRepo struct {
	pool *pgxpool.Pool
}

func NewScheduleRepo(pool *pgxpool.Pool) *ScheduleRepo {
	return &ScheduleRepo{pool: pool}
}

func (r *ScheduleRepo) Create(ctx context.Context, suiteID uuid.UUID, cronExpr string, nextRunAt interface{}, enabled bool) (*model.ScheduleRun, error) {
	var sr model.ScheduleRun
	err := r.pool.QueryRow(ctx,
		`INSERT INTO schedule_runs (suite_id, cron_expr, next_run_at, enabled)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, suite_id, cron_expr, next_run_at, last_run_at, enabled, created_at`,
		suiteID, cronExpr, nextRunAt, enabled,
	).Scan(&sr.ID, &sr.SuiteID, &sr.CronExpr, &sr.NextRunAt, &sr.LastRunAt, &sr.Enabled, &sr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}

func (r *ScheduleRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.ScheduleRun, error) {
	var sr model.ScheduleRun
	err := r.pool.QueryRow(ctx,
		`SELECT id, suite_id, cron_expr, next_run_at, last_run_at, enabled, created_at
		 FROM schedule_runs WHERE id = $1`, id,
	).Scan(&sr.ID, &sr.SuiteID, &sr.CronExpr, &sr.NextRunAt, &sr.LastRunAt, &sr.Enabled, &sr.CreatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &sr, nil
}

func (r *ScheduleRepo) ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]model.ScheduleRun, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, suite_id, cron_expr, next_run_at, last_run_at, enabled, created_at
		 FROM schedule_runs WHERE suite_id = $1 ORDER BY created_at DESC`, suiteID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []model.ScheduleRun
	for rows.Next() {
		var sr model.ScheduleRun
		if err := rows.Scan(&sr.ID, &sr.SuiteID, &sr.CronExpr, &sr.NextRunAt, &sr.LastRunAt, &sr.Enabled, &sr.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, sr)
	}
	return schedules, nil
}

func (r *ScheduleRepo) ListDue(ctx context.Context) ([]model.ScheduleRun, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, suite_id, cron_expr, next_run_at, last_run_at, enabled, created_at
		 FROM schedule_runs WHERE enabled = true AND next_run_at <= now()
		 ORDER BY next_run_at`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []model.ScheduleRun
	for rows.Next() {
		var sr model.ScheduleRun
		if err := rows.Scan(&sr.ID, &sr.SuiteID, &sr.CronExpr, &sr.NextRunAt, &sr.LastRunAt, &sr.Enabled, &sr.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, sr)
	}
	return schedules, nil
}

func (r *ScheduleRepo) UpdateNextRun(ctx context.Context, id uuid.UUID, nextRunAt interface{}) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE schedule_runs SET next_run_at = $2, last_run_at = now() WHERE id = $1`,
		id, nextRunAt,
	)
	return err
}

func (r *ScheduleRepo) Update(ctx context.Context, id uuid.UUID, cronExpr string, nextRunAt interface{}, enabled bool) (*model.ScheduleRun, error) {
	var sr model.ScheduleRun
	err := r.pool.QueryRow(ctx,
		`UPDATE schedule_runs SET cron_expr = $2, next_run_at = $3, enabled = $4
		 WHERE id = $1
		 RETURNING id, suite_id, cron_expr, next_run_at, last_run_at, enabled, created_at`,
		id, cronExpr, nextRunAt, enabled,
	).Scan(&sr.ID, &sr.SuiteID, &sr.CronExpr, &sr.NextRunAt, &sr.LastRunAt, &sr.Enabled, &sr.CreatedAt)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &sr, nil
}

func (r *ScheduleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM schedule_runs WHERE id = $1`, id)
	return err
}
