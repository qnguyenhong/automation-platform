package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Environment represents a test environment configuration.
type Environment struct {
	ID        uuid.UUID       `json:"id"`
	ProjectID uuid.UUID       `json:"project_id"`
	Name      string          `json:"name"`
	BaseURL   string          `json:"base_url"`
	Variables json.RawMessage `json:"variables"`
	Headers   json.RawMessage `json:"headers"`
	IsDefault bool            `json:"is_default"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateEnvironmentRequest struct {
	Name      string          `json:"name" validate:"required"`
	BaseURL   string          `json:"base_url" validate:"required,url"`
	Variables json.RawMessage `json:"variables"`
	Headers   json.RawMessage `json:"headers"`
	IsDefault bool            `json:"is_default"`
}

type UpdateEnvironmentRequest struct {
	Name      string          `json:"name" validate:"required"`
	BaseURL   string          `json:"base_url" validate:"required,url"`
	Variables json.RawMessage `json:"variables"`
	Headers   json.RawMessage `json:"headers"`
	IsDefault bool            `json:"is_default"`
}

// Dataset represents a data-driven testing dataset.
type Dataset struct {
	ID          uuid.UUID       `json:"id"`
	ProjectID   uuid.UUID       `json:"project_id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Format      string          `json:"format"`
	Data        json.RawMessage `json:"data"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateDatasetRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Format      string          `json:"format"`
	Data        json.RawMessage `json:"data" validate:"required"`
}

type UpdateDatasetRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Format      string          `json:"format"`
	Data        json.RawMessage `json:"data" validate:"required"`
}

// Webhook represents a webhook configuration for external triggering.
type Webhook struct {
	ID              uuid.UUID  `json:"id"`
	ProjectID       uuid.UUID  `json:"project_id"`
	Name            string     `json:"name"`
	Secret          string     `json:"secret,omitempty"`
	SuiteID         *uuid.UUID `json:"suite_id,omitempty"`
	Enabled         bool       `json:"enabled"`
	LastTriggeredAt *time.Time `json:"last_triggered_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CreateWebhookRequest struct {
	Name    string     `json:"name" validate:"required"`
	SuiteID *uuid.UUID `json:"suite_id"`
	Enabled bool       `json:"enabled"`
}

type UpdateWebhookRequest struct {
	Name    string     `json:"name" validate:"required"`
	SuiteID *uuid.UUID `json:"suite_id"`
	Enabled bool       `json:"enabled"`
}

// ScheduleRun represents a scheduled test run configuration.
type ScheduleRun struct {
	ID        uuid.UUID  `json:"id"`
	SuiteID   uuid.UUID  `json:"suite_id"`
	CronExpr  string     `json:"cron_expr"`
	NextRunAt time.Time  `json:"next_run_at"`
	LastRunAt *time.Time `json:"last_run_at,omitempty"`
	Enabled   bool       `json:"enabled"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateScheduleRequest struct {
	CronExpr string `json:"cron_expr" validate:"required"`
	Enabled  bool   `json:"enabled"`
}

type UpdateScheduleRequest struct {
	CronExpr string `json:"cron_expr" validate:"required"`
	Enabled  bool   `json:"enabled"`
}
