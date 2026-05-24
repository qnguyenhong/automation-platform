package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	ID            uuid.UUID       `json:"id"`
	Name          string          `json:"name"`
	Hostname      string          `json:"hostname"`
	IPAddress     *string         `json:"ip_address,omitempty"`
	ExecutorTypes []string        `json:"executor_types"`
	Status        WorkerStatus    `json:"status"`
	CurrentJobID  *uuid.UUID      `json:"current_job_id,omitempty"`
	MaxConcurrent int             `json:"max_concurrent"`
	Version       *string         `json:"version,omitempty"`
	Labels        json.RawMessage `json:"labels"`
	LastHeartbeat *time.Time      `json:"last_heartbeat,omitempty"`
	RegisteredAt  time.Time       `json:"registered_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type RegisterWorkerRequest struct {
	Name          string          `json:"name"`
	Hostname      string          `json:"hostname"`
	IPAddress     string          `json:"ip_address"`
	ExecutorTypes []string        `json:"executor_types"`
	MaxConcurrent int             `json:"max_concurrent"`
	Version       string          `json:"version"`
	Labels        json.RawMessage `json:"labels"`
}

type HeartbeatRequest struct {
	Status       WorkerStatus `json:"status"`
	CurrentJobID *string      `json:"current_job_id,omitempty"`
}

type WorkerJob struct {
	JobID      string          `json:"job_id"`
	ResultID   string          `json:"result_id"`
	TestCase   TestCase        `json:"test_case"`
	Suite      TestSuite       `json:"suite"`
	Run        TestRun         `json:"run"`
	Timeout    time.Duration   `json:"timeout"`
}

type JobResult struct {
	Status       StatusType      `json:"status"`
	ErrorMessage string          `json:"error_message"`
	Assertions   json.RawMessage `json:"assertions"`
	RequestData  json.RawMessage `json:"request_data"`
	ResponseData json.RawMessage `json:"response_data"`
	Artifacts    json.RawMessage `json:"artifacts"`
	Metrics      json.RawMessage `json:"metrics"`
	Stdout       string          `json:"stdout"`
	Stderr       string          `json:"stderr"`
	DurationMS   int64           `json:"duration_ms"`
}
