package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TestCase struct {
	ID          uuid.UUID       `json:"id"`
	SuiteID     uuid.UUID       `json:"suite_id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Config      json.RawMessage `json:"config"`
	Tags        []string        `json:"tags"`
	SortOrder   int             `json:"sort_order"`
	Enabled     bool            `json:"enabled"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateTestCaseRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Config      json.RawMessage `json:"config"`
	Tags        []string        `json:"tags"`
	SortOrder   int             `json:"sort_order"`
	Enabled     bool            `json:"enabled"`
}

type UpdateTestCaseRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Config      json.RawMessage `json:"config"`
	Tags        []string        `json:"tags"`
	SortOrder   int             `json:"sort_order"`
	Enabled     bool            `json:"enabled"`
}

// API test case config examples:
// {
//   "method": "GET",
//   "url": "https://api.example.com/users",
//   "headers": {"Authorization": "Bearer xxx"},
//   "assertions": [
//     {"type": "status", "expected": 200},
//     {"type": "json_path", "path": "$.length", "expected": 10}
//   ]
// }
