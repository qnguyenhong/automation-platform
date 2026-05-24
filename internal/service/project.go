package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
)

type ProjectRepository interface {
	Create(ctx context.Context, name, description string, createdBy uuid.UUID) (*model.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	List(ctx context.Context) ([]model.Project, error)
	Update(ctx context.Context, id uuid.UUID, name, description string) (*model.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, req model.CreateProjectRequest, userID string) (*model.Project, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, req.Name, req.Description, uid)
}

func (s *ProjectService) Get(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) List(ctx context.Context) ([]model.Project, error) {
	return s.repo.List(ctx)
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, req model.UpdateProjectRequest) (*model.Project, error) {
	return s.repo.Update(ctx, id, req.Name, req.Description)
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
