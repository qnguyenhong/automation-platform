package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
)

type TestCaseRepository interface {
	Create(ctx context.Context, suiteID uuid.UUID, name, description string, config []byte, tags []string, sortOrder int, enabled bool) (*model.TestCase, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TestCase, error)
	ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]model.TestCase, error)
	Update(ctx context.Context, id uuid.UUID, name, description string, config []byte, tags []string, sortOrder int, enabled bool) (*model.TestCase, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TestCaseService struct {
	repo TestCaseRepository
}

func NewTestCaseService(repo TestCaseRepository) *TestCaseService {
	return &TestCaseService{repo: repo}
}

func (s *TestCaseService) Create(ctx context.Context, suiteID uuid.UUID, req model.CreateTestCaseRequest) (*model.TestCase, error) {
	config := req.Config
	if len(config) == 0 || string(config) == `""` || string(config) == "null" {
		config = []byte("{}")
	}
	return s.repo.Create(ctx, suiteID, req.Name, req.Description, config, req.Tags, req.SortOrder, req.Enabled)
}

func (s *TestCaseService) Get(ctx context.Context, id uuid.UUID) (*model.TestCase, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TestCaseService) ListBySuite(ctx context.Context, suiteID uuid.UUID) ([]model.TestCase, error) {
	return s.repo.ListBySuite(ctx, suiteID)
}

func (s *TestCaseService) Update(ctx context.Context, id uuid.UUID, req model.UpdateTestCaseRequest) (*model.TestCase, error) {
	config := req.Config
	if len(config) == 0 || string(config) == `""` || string(config) == "null" {
		config = []byte("{}")
	}
	return s.repo.Update(ctx, id, req.Name, req.Description, config, req.Tags, req.SortOrder, req.Enabled)
}

func (s *TestCaseService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
