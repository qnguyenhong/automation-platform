package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
)

type TestSuiteRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, name, description string, testType string, scheduleCron *string, config []byte, tags []string, createdBy uuid.UUID) (*model.TestSuite, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TestSuite, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.TestSuite, error)
	Update(ctx context.Context, id uuid.UUID, name, description string, testType string, scheduleCron *string, config []byte, tags []string) (*model.TestSuite, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TestSuiteService struct {
	repo TestSuiteRepository
}

func NewTestSuiteService(repo TestSuiteRepository) *TestSuiteService {
	return &TestSuiteService{repo: repo}
}

func (s *TestSuiteService) Create(ctx context.Context, projectID uuid.UUID, req model.CreateTestSuiteRequest, userID string) (*model.TestSuite, error) {
	var scheduleCron *string
	if req.ScheduleCron != "" {
		scheduleCron = &req.ScheduleCron
	}

	createdBy := uuid.Nil
	if userID != "" {
		if u, err := uuid.Parse(userID); err == nil {
			createdBy = u
		}
	}

	config := req.Config
	if len(config) == 0 || string(config) == `""` || string(config) == "null" {
		config = []byte("{}")
	}

	return s.repo.Create(ctx, projectID, req.Name, req.Description, string(req.TestType), scheduleCron, config, req.Tags, createdBy)
}

func (s *TestSuiteService) Get(ctx context.Context, id uuid.UUID) (*model.TestSuite, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TestSuiteService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.TestSuite, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *TestSuiteService) Update(ctx context.Context, id uuid.UUID, req model.UpdateTestSuiteRequest) (*model.TestSuite, error) {
	var scheduleCron *string
	if req.ScheduleCron != "" {
		scheduleCron = &req.ScheduleCron
	}

	config := req.Config
	if len(config) == 0 || string(config) == `""` || string(config) == "null" {
		config = []byte("{}")
	}

	return s.repo.Update(ctx, id, req.Name, req.Description, string(req.TestType), scheduleCron, config, req.Tags)
}

func (s *TestSuiteService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
