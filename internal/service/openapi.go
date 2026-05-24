package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/openapi"
	"github.com/qnguyenhong/automation-platform/pkg/errors"
)

type OpenAPIRepository interface {
	CreateImport(ctx context.Context, projectID uuid.UUID, filename, specURL string, content []byte, version string, importedBy uuid.UUID) error
	ListImports(ctx context.Context, projectID uuid.UUID) ([]OpenAPIImport, error)
}

type OpenAPIImport struct {
	ID         uuid.UUID `json:"id"`
	ProjectID  uuid.UUID `json:"project_id"`
	Filename   string    `json:"filename"`
	SpecURL    *string   `json:"spec_url,omitempty"`
	Version    string    `json:"version"`
	ImportedBy uuid.UUID `json:"imported_by"`
	ImportedAt string    `json:"imported_at"`
}

type ParseOpenAPIRequest struct {
	Content string `json:"content"`
	URL     string `json:"url"`
}

type ImportEndpointConfig struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Headers     map[string]any    `json:"headers"`
	Params      map[string]any    `json:"params"`
	Body        map[string]any    `json:"body"`
	Assertions  []AssertionConfig `json:"assertions"`
	Variables   map[string]any    `json:"variables"`
}

type AssertionConfig struct {
	Type      string `json:"type"`
	Path      string `json:"path,omitempty"`
	Expected  any    `json:"expected"`
	CaptureAs string `json:"capture_as,omitempty"`
}

type ImportOpenAPIRequest struct {
	ProjectID  uuid.UUID              `json:"project_id"`
	SuiteName  string                 `json:"suite_name"`
	Tags       []string               `json:"tags"`
	Endpoints  []ImportEndpointConfig `json:"endpoints"`
}

type ImportOpenAPIResponse struct {
	SuiteID      uuid.UUID `json:"suite_id"`
	CasesCreated int       `json:"cases_created"`
}

type OpenAPIService struct {
	parser     *openapi.Parser
	repo       OpenAPIRepository
	suiteSvc   *TestSuiteService
	caseSvc    *TestCaseService
}

func NewOpenAPIService(repo OpenAPIRepository, suiteSvc *TestSuiteService, caseSvc *TestCaseService) *OpenAPIService {
	return &OpenAPIService{
		parser:   openapi.NewParser(),
		repo:     repo,
		suiteSvc: suiteSvc,
		caseSvc:  caseSvc,
	}
}

// ParseSpec parses an OpenAPI spec and returns the parsed result.
func (s *OpenAPIService) ParseSpec(ctx context.Context, content string) (*openapi.ParsedSpec, error) {
	if content == "" {
		return nil, errors.NewAppError("BAD_REQUEST", "spec content is required", nil)
	}

	parsed, err := s.parser.ParseAndExtract([]byte(content))
	if err != nil {
		return nil, errors.NewAppError("BAD_REQUEST", fmt.Sprintf("failed to parse spec: %v", err), err)
	}

	return parsed, nil
}

// ImportEndpoints creates a test suite and test cases from imported OpenAPI endpoints.
func (s *OpenAPIService) ImportEndpoints(ctx context.Context, req ImportOpenAPIRequest, userID string) (*ImportOpenAPIResponse, error) {
	if req.ProjectID == uuid.Nil {
		return nil, errors.NewAppError("BAD_REQUEST", "project_id is required", nil)
	}
	if req.SuiteName == "" {
		return nil, errors.NewAppError("BAD_REQUEST", "suite_name is required", nil)
	}
	if len(req.Endpoints) == 0 {
		return nil, errors.NewAppError("BAD_REQUEST", "at least one endpoint is required", nil)
	}

	// Create a test suite
	suite, err := s.suiteSvc.Create(ctx, req.ProjectID, model.CreateTestSuiteRequest{
		Name:     req.SuiteName,
		TestType: model.TestTypeAPI,
		Tags:     req.Tags,
		Config:   json.RawMessage(`{}`),
	}, userID)
	if err != nil {
		return nil, fmt.Errorf("create suite: %w", err)
	}

	// Create test cases for each endpoint
	casesCreated := 0
	for i, ep := range req.Endpoints {
		config := buildTestCaseConfig(ep)

		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("marshal config for %s %s: %w", ep.Method, ep.Path, err)
		}

		caseName := fmt.Sprintf("%s %s", ep.Method, ep.Path)
		if ep.Path == "" {
			continue
		}

		_, err = s.caseSvc.Create(ctx, suite.ID, model.CreateTestCaseRequest{
			Name:      caseName,
			Config:    json.RawMessage(configJSON),
			SortOrder: i,
			Enabled:   true,
		})
		if err != nil {
			return nil, fmt.Errorf("create case for %s %s: %w", ep.Method, ep.Path, err)
		}

		casesCreated++
	}

	return &ImportOpenAPIResponse{
		SuiteID:      suite.ID,
		CasesCreated: casesCreated,
	}, nil
}

// buildTestCaseConfig builds the test case config JSON from an import endpoint config.
func buildTestCaseConfig(ep ImportEndpointConfig) map[string]any {
	config := map[string]any{
		"method": ep.Method,
		"url":    ep.Path,
	}

	if len(ep.Headers) > 0 {
		config["headers"] = ep.Headers
	}
	if len(ep.Params) > 0 {
		config["params"] = ep.Params
	}
	if len(ep.Body) > 0 {
		config["body"] = ep.Body
	}
	if len(ep.Variables) > 0 {
		config["variables"] = ep.Variables
	}

	// Build assertions with capture_as support
	if len(ep.Assertions) > 0 {
		assertions := make([]map[string]any, len(ep.Assertions))
		for i, a := range ep.Assertions {
			assertion := map[string]any{
				"type":     a.Type,
				"expected": a.Expected,
			}
			if a.Path != "" {
				assertion["path"] = a.Path
			}
			if a.CaptureAs != "" {
				assertion["capture_as"] = a.CaptureAs
			}
			assertions[i] = assertion
		}
		config["assertions"] = assertions
	}

	return config
}
