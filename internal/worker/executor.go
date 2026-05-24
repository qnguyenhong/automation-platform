package worker

import (
	"context"
	"time"
)

// Executor defines the interface for test execution engines.
type Executor interface {
	// Name returns the executor type identifier (e.g., "api", "e2e", "load", "unit").
	Name() string

	// CanHandle returns true if this executor supports the given test type.
	CanHandle(testType string) bool

	// Execute runs a single test case and returns the result.
	Execute(ctx context.Context, job Job) (*Result, error)
}

// Job represents a test execution job.
type Job struct {
	JobID        string
	ResultID     string
	TestCaseID   string
	SuiteID      string
	RunID        string
	TestType     string
	Config       map[string]any
	Timeout      time.Duration
	CapturedVars map[string]any    // captured values from prior results in this run
	EnvVars      map[string]string // environment variables
}

// Result represents the outcome of a test execution.
type Result struct {
	Status       string
	ErrorMessage string
	Assertions   []Assertion
	RequestData  map[string]any
	ResponseData map[string]any
	Artifacts    []Artifact
	Metrics      map[string]any
	Stdout       string
	Stderr       string
	DurationMS   int64
	CapturedVars map[string]any // values captured from assertions for chaining
}

// Assertion represents a single test assertion.
type Assertion struct {
	Type     string `json:"type"`
	Expected any    `json:"expected"`
	Actual   any    `json:"actual"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message,omitempty"`
}

// Artifact represents a file or data produced by test execution.
type Artifact struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	URL  string `json:"url,omitempty"`
}

// Registry manages available executors.
type Registry struct {
	executors []Executor
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(executor Executor) {
	r.executors = append(r.executors, executor)
}

func (r *Registry) GetExecutor(testType string) Executor {
	for _, e := range r.executors {
		if e.CanHandle(testType) {
			return e
		}
	}
	return nil
}
