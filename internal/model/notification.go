package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationConfig struct {
	ID        uuid.UUID          `json:"id"`
	ProjectID uuid.UUID          `json:"project_id"`
	Type      NotificationType   `json:"type"`
	Config    json.RawMessage    `json:"config"`
	Events    []string           `json:"events"`
	Enabled   bool               `json:"enabled"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type CreateNotificationRequest struct {
	Type    NotificationType `json:"type"`
	Config  json.RawMessage  `json:"config"`
	Events  []string         `json:"events"`
	Enabled bool             `json:"enabled"`
}

type UpdateNotificationRequest struct {
	Type    NotificationType `json:"type"`
	Config  json.RawMessage  `json:"config"`
	Events  []string         `json:"events"`
	Enabled bool             `json:"enabled"`
}
