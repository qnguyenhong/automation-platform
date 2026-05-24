package model

// StatusType represents the status of a test run or result.
type StatusType string

const (
	StatusPending   StatusType = "pending"
	StatusRunning   StatusType = "running"
	StatusPassed    StatusType = "passed"
	StatusFailed    StatusType = "failed"
	StatusSkipped   StatusType = "skipped"
	StatusError     StatusType = "error"
	StatusCancelled StatusType = "cancelled"
)

// TestType represents the type of test suite.
type TestType string

const (
	TestTypeAPI   TestType = "api"
	TestTypeE2E   TestType = "e2e"
	TestTypeLoad  TestType = "load"
	TestTypeUnit  TestType = "unit"
)

// TriggerType represents how a test run was triggered.
type TriggerType string

const (
	TriggerManual    TriggerType = "manual"
	TriggerScheduled TriggerType = "scheduled"
	TriggerAPI       TriggerType = "api"
	TriggerCI        TriggerType = "ci"
)

// Role represents a user's role.
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleViewer Role = "viewer"
)

// WorkerStatus represents the status of a worker.
type WorkerStatus string

const (
	WorkerOnline   WorkerStatus = "online"
	WorkerOffline  WorkerStatus = "offline"
	WorkerBusy     WorkerStatus = "busy"
	WorkerDraining WorkerStatus = "draining"
)

// NotificationType represents the type of notification.
type NotificationType string

const (
	NotificationSlack   NotificationType = "slack"
	NotificationEmail   NotificationType = "email"
	NotificationWebhook NotificationType = "webhook"
)

// NotificationEvent represents events that trigger notifications.
type NotificationEvent string

const (
	EventRunFailed      NotificationEvent = "run_failed"
	EventRunPassed      NotificationEvent = "run_passed"
	EventFlakyDetected  NotificationEvent = "flaky_detected"
)
