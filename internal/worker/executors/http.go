package executors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qnguyenhong/automation-platform/internal/worker"
	"github.com/qnguyenhong/automation-platform/pkg/variable"
)

// HTTPExecutor executes API/HTTP tests.
type HTTPExecutor struct {
	client *http.Client
}

func NewHTTPExecutor(timeout time.Duration) *HTTPExecutor {
	return &HTTPExecutor{
		client: &http.Client{Timeout: timeout},
	}
}

func (e *HTTPExecutor) Name() string {
	return "api"
}

func (e *HTTPExecutor) CanHandle(testType string) bool {
	return testType == "api"
}

func (e *HTTPExecutor) Execute(ctx context.Context, job worker.Job) (*worker.Result, error) {
	start := time.Now()

	result := &worker.Result{
		Status:    "passed",
		Assertions: []worker.Assertion{},
	}

	// Resolve variables in config
	engine := variable.NewEngine(job.CapturedVars, job.EnvVars)
	resolvedConfig, err := engine.ResolveMap(job.Config)
	if err != nil {
		return &worker.Result{
			Status:       "error",
			ErrorMessage: fmt.Sprintf("variable resolution failed: %v", err),
			DurationMS:   time.Since(start).Milliseconds(),
		}, nil
	}

	// Parse config
	method := getString(resolvedConfig, "method", "GET")
	url := getString(resolvedConfig, "url", "")
	headers := getMap(resolvedConfig, "headers")
	body := resolvedConfig["body"]
	assertions := getAssertions(resolvedConfig)

	if url == "" {
		return &worker.Result{
			Status:       "error",
			ErrorMessage: "url is required",
			DurationMS:   time.Since(start).Milliseconds(),
		}, nil
	}

	// Build request
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return &worker.Result{
				Status:       "error",
				ErrorMessage: fmt.Sprintf("failed to marshal body: %v", err),
				DurationMS:   time.Since(start).Milliseconds(),
			}, nil
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return &worker.Result{
			Status:       "error",
			ErrorMessage: fmt.Sprintf("failed to create request: %v", err),
			DurationMS:   time.Since(start).Milliseconds(),
		}, nil
	}

	// Set headers
	for k, v := range headers {
		if s, ok := v.(string); ok {
			req.Header.Set(k, s)
		}
	}

	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Record request data
	result.RequestData = map[string]any{
		"method":  method,
		"url":     url,
		"headers": headers,
		"body":    body,
	}

	// Execute request
	resp, err := e.client.Do(req)
	if err != nil {
		return &worker.Result{
			Status:       "error",
			ErrorMessage: fmt.Sprintf("request failed: %v", err),
			RequestData:  result.RequestData,
			DurationMS:   time.Since(start).Milliseconds(),
		}, nil
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &worker.Result{
			Status:       "error",
			ErrorMessage: fmt.Sprintf("failed to read response: %v", err),
			RequestData:  result.RequestData,
			DurationMS:   time.Since(start).Milliseconds(),
		}, nil
	}

	duration := time.Since(start).Milliseconds()

	// Parse response body as JSON if possible
	var respJSON any
	if err := json.Unmarshal(respBody, &respJSON); err != nil {
		respJSON = string(respBody)
	}

	result.ResponseData = map[string]any{
		"status": resp.StatusCode,
		"headers": resp.Header,
		"body": respJSON,
	}

	// Run assertions
	allPassed := true
	for _, assertion := range assertions {
		workerAssertion := runAssertion(assertion, resp.StatusCode, resp.Header, respJSON, duration)
		result.Assertions = append(result.Assertions, workerAssertion)
		if !workerAssertion.Passed {
			allPassed = false
		}
		// Capture value for chaining if capture_as is set
		if assertion.CaptureAs != "" && workerAssertion.Actual != nil {
			engine.Capture(assertion.CaptureAs, workerAssertion.Actual)
		}
	}

	if !allPassed {
		result.Status = "failed"
		var failedMessages []string
		for _, a := range result.Assertions {
			if !a.Passed {
				failedMessages = append(failedMessages, a.Message)
			}
		}
		result.ErrorMessage = strings.Join(failedMessages, "; ")
	}

	result.DurationMS = duration
	result.CapturedVars = engine.GetCaptured()
	return result, nil
}

type assertionDef struct {
	Type      string `json:"type"`
	Path      string `json:"path,omitempty"`
	Expected  any    `json:"expected"`
	CaptureAs string `json:"capture_as,omitempty"`
}

func getAssertions(config map[string]any) []assertionDef {
	var assertions []assertionDef
	raw, ok := config["assertions"]
	if !ok {
		return assertions
	}

	switch v := raw.(type) {
	case []any:
		for _, a := range v {
			if m, ok := a.(map[string]any); ok {
				assertions = append(assertions, assertionDef{
					Type:      getString(m, "type", ""),
					Path:      getString(m, "path", ""),
					Expected:  m["expected"],
					CaptureAs: getString(m, "capture_as", ""),
				})
			}
		}
	}
	return assertions
}

func runAssertion(assertion assertionDef, statusCode int, headers http.Header, body any, duration int64) worker.Assertion {
	a := worker.Assertion{
		Type:     assertion.Type,
		Expected: assertion.Expected,
	}

	switch assertion.Type {
	case "status":
		a.Actual = statusCode
		if expected, ok := toInt(assertion.Expected); ok {
			a.Passed = statusCode == expected
		}
		if !a.Passed {
			a.Message = fmt.Sprintf("expected status %v, got %d", assertion.Expected, statusCode)
		}

	case "response_time":
		a.Actual = duration
		if expected, ok := toInt64(assertion.Expected); ok {
			a.Passed = duration <= expected
		}
		if !a.Passed {
			a.Message = fmt.Sprintf("expected response time <= %vms, got %dms", assertion.Expected, duration)
		}

	case "header":
		actual := headers.Get(assertion.Path)
		a.Actual = actual
		a.Passed = actual == fmt.Sprintf("%v", assertion.Expected)
		if !a.Passed {
			a.Message = fmt.Sprintf("expected header %s = %v, got %s", assertion.Path, assertion.Expected, actual)
		}

	case "json_path":
		actual := getJSONPath(body, assertion.Path)
		a.Actual = actual
		a.Passed = fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", assertion.Expected)
		if !a.Passed {
			a.Message = fmt.Sprintf("expected %s = %v, got %v", assertion.Path, assertion.Expected, actual)
		}

	default:
		a.Passed = true
		a.Message = fmt.Sprintf("unknown assertion type: %s", assertion.Type)
	}

	return a
}

func getString(m map[string]any, key, defaultVal string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func getMap(m map[string]any, key string) map[string]any {
	if v, ok := m[key]; ok {
		if m, ok := v.(map[string]any); ok {
			return m
		}
	}
	return nil
}

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	}
	return 0, false
}

func toInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	case json.Number:
		i, err := n.Int64()
		return i, err == nil
	}
	return 0, false
}

func getJSONPath(body any, path string) any {
	if path == "" || path == "$" {
		return body
	}

	// Simple JSON path implementation - supports $.key.subkey
	parts := strings.Split(strings.TrimPrefix(path, "$."), ".")
	current := body

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			current = v[part]
		default:
			return nil
		}
	}

	return current
}
