package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
	repo       WorkerRepository
	jwtSecret  []byte
	dispatcher interface {
		Dequeue() *model.WorkerJob
	}
}

func NewWorkerService(repo WorkerRepository, jwtSecret string) *WorkerService {
	return &WorkerService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *WorkerService) SetDispatcher(dispatcher interface{ Dequeue() *model.WorkerJob }) {
	s.dispatcher = dispatcher
}

func pingWorker(hostname, ipAddress string) error {
	port := "9090"
	if envPort := os.Getenv("WORKER_PING_PORT"); envPort != "" {
		port = envPort
	}

	targets := []string{ipAddress, hostname}
	var lastErr error

	for _, target := range targets {
		if target == "" {
			continue
		}
		pingURL := fmt.Sprintf("http://%s:%s/ping", target, port)
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get(pingURL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			lastErr = fmt.Errorf("status code: %d", resp.StatusCode)
		} else {
			lastErr = err
		}
	}

	if lastErr != nil {
		return fmt.Errorf("failed to ping worker on targets %v: %w", targets, lastErr)
	}
	return fmt.Errorf("no target to ping")
}

func (s *WorkerService) StartHealthCheck(ctx context.Context, logger *slog.Logger) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				s.checkWorkers(ctx, logger)
			}
		}
	}()
}

func (s *WorkerService) checkWorkers(ctx context.Context, logger *slog.Logger) {
	workers, err := s.repo.List(ctx)
	if err != nil {
		logger.Error("failed to list workers for health check", "error", err)
		return
	}

	for _, w := range workers {
		workerAddr := w.Hostname
		ipAddr := ""
		if w.IPAddress != nil {
			ipAddr = *w.IPAddress
		}

		isLive := false
		err := pingWorker(workerAddr, ipAddr)
		if err == nil {
			isLive = true
		}

		// Fallback check: last heartbeat within 30 seconds
		if !isLive && w.LastHeartbeat != nil {
			if time.Since(*w.LastHeartbeat) < 30*time.Second {
				isLive = true
			}
		}

		if !isLive {
			logger.Warn("worker is not live, removing it", "id", w.ID, "name", w.Name, "hostname", workerAddr)
			if err := s.repo.Delete(ctx, w.ID); err != nil {
				logger.Error("failed to delete non-live worker", "id", w.ID, "error", err)
			}
		}
	}
}

func (s *WorkerService) Register(ctx context.Context, req model.RegisterWorkerRequest) (*model.Worker, string, error) {
	// Active ping-pong check before registering
	if err := pingWorker(req.Hostname, req.IPAddress); err != nil {
		return nil, "", errors.NewAppError("VALIDATION", fmt.Sprintf("worker is not active or did not respond to ping: %v", err), nil)
	}

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
	if s.dispatcher == nil {
		return nil, nil
	}
	return s.dispatcher.Dequeue(), nil
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
