package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/qnguyenhong/automation-platform/pkg/config"
)

// Agent represents a worker agent that polls for and executes test jobs.
type Agent struct {
	config   *config.WorkerConfig
	logger   *slog.Logger
	client   *http.Client
	registry *Registry

	workerID string
	token    string
	status   string

	mu       sync.Mutex
	running  map[string]context.CancelFunc
}

func NewAgent(cfg *config.WorkerConfig, logger *slog.Logger, registry *Registry) *Agent {
	return &Agent{
		config:   cfg,
		logger:   logger,
		client:   &http.Client{Timeout: 30 * time.Second},
		registry: registry,
		status:   "online",
		running:  make(map[string]context.CancelFunc),
	}
}

// Start begins the worker agent lifecycle: register, heartbeat, and job polling.
func (a *Agent) Start(ctx context.Context) error {
	// Register with server
	if err := a.register(ctx); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	a.logger.Info("worker registered", "id", a.workerID, "name", a.config.Name)

	// Start heartbeat goroutine
	go a.heartbeatLoop(ctx)

	// Start job polling
	a.pollLoop(ctx)

	return nil
}

func (a *Agent) register(ctx context.Context) error {
	payload := map[string]any{
		"name":           a.config.Name,
		"hostname":       getHostname(),
		"ip_address":     getLocalIP(),
		"executor_types": a.config.ExecutorTypes,
		"max_concurrent": a.config.MaxConcurrent,
		"version":        "0.1.0",
		"labels":         map[string]string{},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", a.config.ServerURL+"/api/v1/workers/register", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("register request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("register failed: status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Worker struct {
				ID string `json:"id"`
			} `json:"worker"`
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode register response: %w", err)
	}

	a.workerID = result.Data.Worker.ID
	a.token = result.Data.Token
	a.registry.SetToken(a.token)
	return nil
}

func (a *Agent) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(a.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.sendHeartbeat(ctx)
		}
	}
}

func (a *Agent) sendHeartbeat(ctx context.Context) {
	a.mu.Lock()
	status := a.status
	runningCount := len(a.running)
	a.mu.Unlock()

	if runningCount >= a.config.MaxConcurrent {
		status = "busy"
	}

	payload := map[string]any{
		"status": status,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "PUT",
		fmt.Sprintf("%s/api/v1/workers/%s/heartbeat", a.config.ServerURL, a.workerID),
		bytes.NewReader(body))
	if err != nil {
		a.logger.Error("heartbeat request creation failed", "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Worker-Token", a.token)

	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Error("heartbeat failed", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.logger.Warn("heartbeat returned unexpected status", "status", resp.StatusCode)
	}
}

func (a *Agent) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(a.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.mu.Lock()
			runningCount := len(a.running)
			a.mu.Unlock()

			if runningCount >= a.config.MaxConcurrent {
				continue
			}

			job, err := a.pollNextJob(ctx)
			if err != nil {
				a.logger.Error("poll failed", "error", err)
				continue
			}

			if job == nil {
				continue
			}

			go a.executeJob(ctx, *job)
		}
	}
}

func (a *Agent) pollNextJob(ctx context.Context) (*Job, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.config.ServerURL+"/api/v1/jobs/next", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Worker-Token", a.token)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			JobID    string `json:"job_id"`
			ResultID string `json:"result_id"`
			TestCase struct {
				ID     string          `json:"id"`
				Config json.RawMessage `json:"config"`
			} `json:"test_case"`
			Suite struct {
				ID       string `json:"id"`
				TestType string `json:"test_type"`
			} `json:"suite"`
			Run struct {
				ID string `json:"id"`
			} `json:"run"`
			Timeout time.Duration `json:"timeout"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}

	var config map[string]any
	if len(envelope.Data.TestCase.Config) > 0 {
		_ = json.Unmarshal(envelope.Data.TestCase.Config, &config)
	}

	job := Job{
		JobID:      envelope.Data.JobID,
		ResultID:   envelope.Data.ResultID,
		TestCaseID: envelope.Data.TestCase.ID,
		SuiteID:    envelope.Data.Suite.ID,
		RunID:      envelope.Data.Run.ID,
		TestType:   envelope.Data.Suite.TestType,
		Config:     config,
		Timeout:    envelope.Data.Timeout,
	}

	return &job, nil
}

func (a *Agent) executeJob(ctx context.Context, job Job) {
	a.mu.Lock()
	jobCtx, cancel := context.WithTimeout(ctx, job.Timeout)
	a.running[job.JobID] = cancel
	a.mu.Unlock()

	defer func() {
		cancel()
		a.mu.Lock()
		delete(a.running, job.JobID)
		a.mu.Unlock()
	}()

	a.logger.Info("executing job", "job_id", job.JobID, "case", job.TestCaseID, "type", job.TestType)

	executor := a.registry.GetExecutor(job.TestType)
	if executor == nil {
		a.reportError(job, "no executor for test type: "+job.TestType)
		return
	}

	result, err := executor.Execute(jobCtx, job)
	if err != nil {
		a.reportError(job, err.Error())
		return
	}

	a.reportResult(job, result)
}

func (a *Agent) reportResult(job Job, result *Result) {
	payload := map[string]any{
		"status":        result.Status,
		"error_message": result.ErrorMessage,
		"assertions":    result.Assertions,
		"request_data":  result.RequestData,
		"response_data": result.ResponseData,
		"artifacts":     result.Artifacts,
		"metrics":       result.Metrics,
		"stdout":        result.Stdout,
		"stderr":        result.Stderr,
		"duration_ms":   result.DurationMS,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/jobs/%s/result", a.config.ServerURL, job.JobID),
		bytes.NewReader(body))
	if err != nil {
		a.logger.Error("create result request failed", "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Worker-Token", a.token)

	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Error("report result failed", "error", err)
		return
	}
	defer resp.Body.Close()

	a.logger.Info("result reported", "job_id", job.JobID, "status", result.Status)
}

func (a *Agent) reportError(job Job, message string) {
	a.reportResult(job, &Result{
		Status:       "error",
		ErrorMessage: message,
	})
}

func getHostname() string {
	hostname, _ := getHostnameFromOS()
	return hostname
}

func getLocalIP() string {
	return "127.0.0.1"
}
