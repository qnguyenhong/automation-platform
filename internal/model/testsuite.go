package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TestSuite struct {
	ID          uuid.UUID       `json:"id"`
	ProjectID   uuid.UUID       `json:"project_id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	TestType    TestType        `json:"test_type"`
	ScheduleCron *string        `json:"schedule_cron,omitempty"`
	Config      json.RawMessage `json:"config"`
	Tags        []string        `json:"tags"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateTestSuiteRequest struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	TestType     TestType        `json:"test_type"`
	ScheduleCron string          `json:"schedule_cron"`
	Config       json.RawMessage `json:"config"`
	Tags         []string        `json:"tags"`
}

type UpdateTestSuiteRequest struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	TestType     TestType        `json:"test_type"`
	ScheduleCron string          `json:"schedule_cron"`
	Config       json.RawMessage `json:"config"`
	Tags         []string        `json:"tags"`
}
