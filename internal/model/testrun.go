package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TestRun struct {
	ID         uuid.UUID       `json:"id"`
	SuiteID    uuid.UUID       `json:"suite_id"`
	ProjectID  uuid.UUID       `json:"project_id"`
	Status     StatusType      `json:"status"`
	Trigger    TriggerType     `json:"trigger"`
	TriggeredBy *uuid.UUID     `json:"triggered_by,omitempty"`
	TotalCases int             `json:"total_cases"`
	Passed     int             `json:"passed"`
	Failed     int             `json:"failed"`
	Skipped    int             `json:"skipped"`
	Errored    int             `json:"errored"`
	StartedAt  *time.Time      `json:"started_at,omitempty"`
	FinishedAt *time.Time      `json:"finished_at,omitempty"`
	DurationMS *int64          `json:"duration_ms,omitempty"`
	Metadata   json.RawMessage `json:"metadata"`
	CreatedAt  time.Time       `json:"created_at"`

	// Joined fields
	SuiteName   string `json:"suite_name,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
}

type TestResult struct {
	ID           uuid.UUID       `json:"id"`
	RunID        uuid.UUID       `json:"run_id"`
	CaseID       uuid.UUID       `json:"case_id"`
	WorkerID     *uuid.UUID      `json:"worker_id,omitempty"`
	Status       StatusType      `json:"status"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	Assertions   json.RawMessage `json:"assertions"`
	RequestData  json.RawMessage `json:"request_data,omitempty"`
	ResponseData json.RawMessage `json:"response_data,omitempty"`
	Artifacts    json.RawMessage `json:"artifacts"`
	Metrics      json.RawMessage `json:"metrics,omitempty"`
	Stdout       *string         `json:"stdout,omitempty"`
	Stderr       *string         `json:"stderr,omitempty"`
	DurationMS   *int64          `json:"duration_ms,omitempty"`
	StartedAt    *time.Time      `json:"started_at,omitempty"`
	FinishedAt   *time.Time      `json:"finished_at,omitempty"`
	RetryOf      *uuid.UUID      `json:"retry_of,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

type TriggerRunRequest struct {
	Trigger TriggerType     `json:"trigger"`
	Metadata json.RawMessage `json:"metadata"`
}

type RunResultCounts struct {
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
	Errored int `json:"errored"`
}
