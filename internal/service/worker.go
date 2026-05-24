package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/pkg/errors"
)

type WorkerRepository interface {
	Create(ctx context.Context, name, hostname, ipAddress string, executorTypes []string, maxConcurrent int, version *string, labels []byte) (*model.Worker, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Worker, error)
	List(ctx context.Context) ([]model.Worker, error)
	ListByStatus(ctx context.Context, status string) ([]model.Worker, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*model.Worker, error)
	UpdateHeartbeat(ctx context.Context, id uuid.UUID, status string, currentJobID *uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type WorkerService struct {
	repo      WorkerRepository
	jwtSecret []byte
}

func NewWorkerService(repo WorkerRepository, jwtSecret string) *WorkerService {
	return &WorkerService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *WorkerService) Register(ctx context.Context, req model.RegisterWorkerRequest) (*model.Worker, string, error) {
	version := &req.Version
	if req.Version == "" {
		version = nil
	}

	worker, err := s.repo.Create(ctx, req.Name, req.Hostname, req.IPAddress, req.ExecutorTypes, req.MaxConcurrent, version, req.Labels)
	if err != nil {
		return nil, "", err
	}

	// Generate worker token
	token, err := s.generateWorkerToken(worker.ID)
	if err != nil {
		return nil, "", errors.NewAppError("INTERNAL", "failed to generate worker token", err)
	}

	// Update status to online
	_, err = s.repo.UpdateStatus(ctx, worker.ID, string(model.WorkerOnline))
	if err != nil {
		return nil, "", err
	}

	return worker, token, nil
}

func (s *WorkerService) Get(ctx context.Context, id uuid.UUID) (*model.Worker, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WorkerService) List(ctx context.Context) ([]model.Worker, error) {
	return s.repo.List(ctx)
}

func (s *WorkerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *WorkerService) Heartbeat(ctx context.Context, workerID uuid.UUID, req model.HeartbeatRequest) error {
	var currentJobID *uuid.UUID
	if req.CurrentJobID != nil && *req.CurrentJobID != "" {
		id, err := uuid.Parse(*req.CurrentJobID)
		if err == nil {
			currentJobID = &id
		}
	}

	return s.repo.UpdateHeartbeat(ctx, workerID, string(req.Status), currentJobID)
}

func (s *WorkerService) GetNextJob(ctx context.Context, workerID uuid.UUID) (*model.WorkerJob, error) {
	// This would query for pending jobs matching worker capabilities
	// For now return nil (no jobs available)
	return nil, nil
}

func (s *WorkerService) ValidateWorkerToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.ErrUnauthorized
	}

	workerID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.ErrUnauthorized
	}

	return workerID, nil
}

func (s *WorkerService) generateWorkerToken(workerID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": workerID.String(),
		"exp": time.Now().Add(365 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
