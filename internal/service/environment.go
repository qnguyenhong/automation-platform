package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/pkg/errors"
)

// EnvironmentRepository defines the data access interface for environments.
type EnvironmentRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, name, baseURL string, variables, headers []byte, isDefault bool) (*model.Environment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Environment, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error)
	Update(ctx context.Context, id uuid.UUID, name, baseURL string, variables, headers []byte, isDefault bool) (*model.Environment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetDefault(ctx context.Context, projectID uuid.UUID) (*model.Environment, error)
}

// DatasetRepository defines the data access interface for datasets.
type DatasetRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, name, description, format string, data []byte) (*model.Dataset, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Dataset, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Dataset, error)
	Update(ctx context.Context, id uuid.UUID, name, description, format string, data []byte) (*model.Dataset, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// WebhookRepository defines the data access interface for webhooks.
type WebhookRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, name string, suiteID *uuid.UUID, enabled bool) (*model.Webhook, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Webhook, error)
	GetBySecret(ctx context.Context, secret string) (*model.Webhook, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Webhook, error)
	Update(ctx context.Context, id uuid.UUID, name string, suiteID *uuid.UUID, enabled bool) (*model.Webhook, error)
	UpdateLastTriggered(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ScheduleRepository defines the data access interface for scheduled runs.
type ScheduleRepository interface {
	Create(ctx context.Context, suiteID uuid.UUID, cronExpr string, nextRunAt interface{}, enabled bool) (*model.ScheduleRun, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.ScheduleRun, error)
	ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]model.ScheduleRun, error)
	ListDue(ctx context.Context) ([]model.ScheduleRun, error)
	UpdateNextRun(ctx context.Context, id uuid.UUID, nextRunAt interface{}) error
	Update(ctx context.Context, id uuid.UUID, cronExpr string, nextRunAt interface{}, enabled bool) (*model.ScheduleRun, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EnvironmentService handles environment business logic.
type EnvironmentService struct {
	repo EnvironmentRepository
}

func NewEnvironmentService(repo EnvironmentRepository) *EnvironmentService {
	return &EnvironmentService{repo: repo}
}

func (s *EnvironmentService) Create(ctx context.Context, projectID uuid.UUID, req model.CreateEnvironmentRequest) (*model.Environment, error) {
	if req.Name == "" {
		return nil, errors.NewAppError("BAD_REQUEST", "name is required", nil)
	}
	if req.BaseURL == "" {
		return nil, errors.NewAppError("BAD_REQUEST", "base_url is required", nil)
	}
	variables := req.Variables
	if variables == nil {
		variables = json.RawMessage(`{}`)
	}
	headers := req.Headers
	if headers == nil {
		headers = json.RawMessage(`{}`)
	}
	return s.repo.Create(ctx, projectID, req.Name, req.BaseURL, variables, headers, req.IsDefault)
}

func (s *EnvironmentService) Get(ctx context.Context, id uuid.UUID) (*model.Environment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EnvironmentService) List(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *EnvironmentService) Update(ctx context.Context, id uuid.UUID, req model.UpdateEnvironmentRequest) (*model.Environment, error) {
	if req.Name == "" {
		return nil, errors.NewAppError("BAD_REQUEST", "name is required", nil)
	}
	variables := req.Variables
	if variables == nil {
		variables = json.RawMessage(`{}`)
	}
	headers := req.Headers
	if headers == nil {
		headers = json.RawMessage(`{}`)
	}
	return s.repo.Update(ctx, id, req.Name, req.BaseURL, variables, headers, req.IsDefault)
}

func (s *EnvironmentService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *EnvironmentService) GetDefault(ctx context.Context, projectID uuid.UUID) (*model.Environment, error) {
	return s.repo.GetDefault(ctx, projectID)
}

// DatasetService handles dataset business logic.
type DatasetService struct {
	repo DatasetRepository
}

func NewDatasetService(repo DatasetRepository) *DatasetService {
	return &DatasetService{repo: repo}
}

func (s *DatasetService) Create(ctx context.Context, projectID uuid.UUID, req model.CreateDatasetRequest) (*model.Dataset, error) {
	if req.Name == "" {
		return nil, errors.NewAppError("BAD_REQUEST", "name is required", nil)
	}
	format := req.Format
	if format == "" {
		format = "json"
	}
	data := req.Data
	if data == nil {
		data = json.RawMessage(`[]`)
	}
	return s.repo.Create(ctx, projectID, req.Name, req.Description, format, data)
}

func (s *DatasetService) Get(ctx context.Context, id uuid.UUID) (*model.Dataset, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DatasetService) List(ctx context.Context, projectID uuid.UUID) ([]model.Dataset, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *DatasetService) Update(ctx context.Context, id uuid.UUID, req model.UpdateDatasetRequest) (*model.Dataset, error) {
	format := req.Format
	if format == "" {
		format = "json"
	}
	return s.repo.Update(ctx, id, req.Name, req.Description, format, req.Data)
}

func (s *DatasetService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// WebhookService handles webhook business logic.
type WebhookService struct {
	repo   WebhookRepository
	runSvc *TestRunService
}

func NewWebhookService(repo WebhookRepository, runSvc *TestRunService) *WebhookService {
	return &WebhookService{repo: repo, runSvc: runSvc}
}

func (s *WebhookService) Create(ctx context.Context, projectID uuid.UUID, req model.CreateWebhookRequest) (*model.Webhook, error) {
	if req.Name == "" {
		return nil, errors.NewAppError("BAD_REQUEST", "name is required", nil)
	}
	return s.repo.Create(ctx, projectID, req.Name, req.SuiteID, req.Enabled)
}

func (s *WebhookService) Get(ctx context.Context, id uuid.UUID) (*model.Webhook, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WebhookService) List(ctx context.Context, projectID uuid.UUID) ([]model.Webhook, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *WebhookService) Update(ctx context.Context, id uuid.UUID, req model.UpdateWebhookRequest) (*model.Webhook, error) {
	return s.repo.Update(ctx, id, req.Name, req.SuiteID, req.Enabled)
}

func (s *WebhookService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Trigger triggers a test run via webhook secret.
func (s *WebhookService) Trigger(ctx context.Context, secret string) (*model.TestRun, error) {
	wh, err := s.repo.GetBySecret(ctx, secret)
	if err != nil {
		return nil, errors.NewAppError("UNAUTHORIZED", "invalid webhook secret", err)
	}

	if wh.SuiteID == nil {
		return nil, errors.NewAppError("BAD_REQUEST", "webhook has no suite configured", nil)
	}

	// Trigger the run
	run, err := s.runSvc.Trigger(ctx, *wh.SuiteID, "", model.TriggerRunRequest{
		Trigger:  model.TriggerAPI,
		Metadata: json.RawMessage(`{"source": "webhook", "webhook_id": "` + wh.ID.String() + `"}`),
	})
	if err != nil {
		return nil, err
	}

	// Update last triggered
	_ = s.repo.UpdateLastTriggered(ctx, wh.ID)

	return run, nil
}
