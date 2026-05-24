package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
)

// ExportService handles exporting test suites to various formats.
type ExportService struct {
	suiteSvc *TestSuiteService
	caseSvc  *TestCaseService
	envSvc   *EnvironmentService
}

func NewExportService(suiteSvc *TestSuiteService, caseSvc *TestCaseService, envSvc *EnvironmentService) *ExportService {
	return &ExportService{
		suiteSvc: suiteSvc,
		caseSvc:  caseSvc,
		envSvc:   envSvc,
	}
}

// ExportFormat represents supported export formats.
type ExportFormat string

const (
	ExportPostman ExportFormat = "postman"
	ExportCurl    ExportFormat = "curl"
	ExportK6      ExportFormat = "k6"
)

// ExportRequest represents an export request.
type ExportRequest struct {
	SuiteID       uuid.UUID    `json:"suite_id"`
	Format        ExportFormat `json:"format"`
	EnvironmentID *uuid.UUID   `json:"environment_id,omitempty"`
}

// ExportResponse holds the exported content.
type ExportResponse struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	MimeType string `json:"mime_type"`
}

func (s *ExportService) Export(ctx context.Context, req ExportRequest) (*ExportResponse, error) {
	suite, err := s.suiteSvc.Get(ctx, req.SuiteID)
	if err != nil {
		return nil, err
	}

	cases, err := s.caseSvc.ListBySuite(ctx, req.SuiteID)
	if err != nil {
		return nil, err
	}

	// Get environment if specified
	var env *model.Environment
	if req.EnvironmentID != nil {
		env, _ = s.envSvc.Get(ctx, *req.EnvironmentID)
	}

	switch req.Format {
	case ExportPostman:
		return s.exportPostman(suite, cases, env)
	case ExportCurl:
		return s.exportCurl(suite, cases, env)
	case ExportK6:
		return s.exportK6(suite, cases, env)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", req.Format)
	}
}

func (s *ExportService) exportPostman(suite *model.TestSuite, cases []model.TestCase, env *model.Environment) (*ExportResponse, error) {
	baseURL := ""
	if env != nil {
		baseURL = env.BaseURL
	}

	collection := map[string]any{
		"info": map[string]any{
			"name":        suite.Name,
			"description": suite.Description,
			"schema":      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"item": buildPostmanItems(cases, baseURL),
	}

	content, _ := json.MarshalIndent(collection, "", "  ")
	return &ExportResponse{
		Filename: sanitizeFilename(suite.Name) + ".postman_collection.json",
		Content:  string(content),
		MimeType: "application/json",
	}, nil
}

func buildPostmanItems(cases []model.TestCase, baseURL string) []map[string]any {
	items := make([]map[string]any, 0, len(cases))
	for _, tc := range cases {
		var config map[string]any
		json.Unmarshal(tc.Config, &config)

		method, _ := config["method"].(string)
		url, _ := config["url"].(string)
		headers, _ := config["headers"].(map[string]any)
		body := config["body"]

		if baseURL != "" && !strings.HasPrefix(url, "http") {
			url = baseURL + url
		}

		item := map[string]any{
			"name": tc.Name,
			"request": map[string]any{
				"method": method,
				"url":    url,
				"header": buildPostmanHeaders(headers),
			},
		}

		if body != nil {
			bodyJSON, _ := json.Marshal(body)
			item["request"].(map[string]any)["body"] = map[string]any{
				"mode": "raw",
				"raw":  string(bodyJSON),
				"options": map[string]any{
					"raw": map[string]any{
						"language": "json",
					},
				},
			}
		}

		items = append(items, item)
	}
	return items
}

func buildPostmanHeaders(headers map[string]any) []map[string]any {
	result := make([]map[string]any, 0)
	for k, v := range headers {
		result = append(result, map[string]any{
			"key":   k,
			"value": fmt.Sprintf("%v", v),
		})
	}
	return result
}

func (s *ExportService) exportCurl(suite *model.TestSuite, cases []model.TestCase, env *model.Environment) (*ExportResponse, error) {
	baseURL := ""
	if env != nil {
		baseURL = env.BaseURL
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n# Generated from automation platform\n\n", suite.Name))

	for _, tc := range cases {
		var config map[string]any
		json.Unmarshal(tc.Config, &config)

		method, _ := config["method"].(string)
		url, _ := config["url"].(string)
		headers, _ := config["headers"].(map[string]any)
		body := config["body"]

		if baseURL != "" && !strings.HasPrefix(url, "http") {
			url = baseURL + url
		}

		sb.WriteString(fmt.Sprintf("# %s\n", tc.Name))
		sb.WriteString(fmt.Sprintf("curl -X %s '%s'", method, url))

		for k, v := range headers {
			sb.WriteString(fmt.Sprintf(" \\\n  -H '%s: %v'", k, v))
		}

		if body != nil {
			bodyJSON, _ := json.Marshal(body)
			sb.WriteString(fmt.Sprintf(" \\\n  -d '%s'", string(bodyJSON)))
		}

		sb.WriteString("\n\n")
	}

	return &ExportResponse{
		Filename: sanitizeFilename(suite.Name) + ".sh",
		Content:  sb.String(),
		MimeType: "text/x-shellscript",
	}, nil
}

func (s *ExportService) exportK6(suite *model.TestSuite, cases []model.TestCase, env *model.Environment) (*ExportResponse, error) {
	baseURL := ""
	if env != nil {
		baseURL = env.BaseURL
	}

	var sb strings.Builder
	sb.WriteString("import http from 'k6/http';\n")
	sb.WriteString("import { check, sleep } from 'k6';\n\n")
	sb.WriteString("export const options = {\n")
	sb.WriteString("  vus: 10,\n")
	sb.WriteString("  duration: '30s',\n")
	sb.WriteString("};\n\n")

	if baseURL != "" {
		sb.WriteString(fmt.Sprintf("const BASE_URL = '%s';\n\n", baseURL))
	}

	sb.WriteString("export default function () {\n")

	for _, tc := range cases {
		var config map[string]any
		json.Unmarshal(tc.Config, &config)

		method, _ := config["method"].(string)
		url, _ := config["url"].(string)
		headers, _ := config["headers"].(map[string]any)
		body := config["body"]

		if baseURL != "" && !strings.HasPrefix(url, "http") {
			url = "BASE_URL + '" + url + "'"
		} else {
			url = "'" + url + "'"
		}

		sb.WriteString(fmt.Sprintf("  // %s\n", tc.Name))

		headersJSON, _ := json.Marshal(headers)

		if body != nil {
			bodyJSON, _ := json.Marshal(body)
			sb.WriteString(fmt.Sprintf("  let res_%s = http.%s(%s, JSON.stringify(%s), { headers: %s });\n",
				sanitizeVarName(tc.Name), strings.ToLower(method), url, string(bodyJSON), string(headersJSON)))
		} else {
			sb.WriteString(fmt.Sprintf("  let res_%s = http.%s(%s, { headers: %s });\n",
				sanitizeVarName(tc.Name), strings.ToLower(method), url, string(headersJSON)))
		}

		sb.WriteString(fmt.Sprintf("  check(res_%s, { '%s status 200': (r) => r.status === 200 });\n\n",
			sanitizeVarName(tc.Name), tc.Name))
	}

	sb.WriteString("  sleep(1);\n")
	sb.WriteString("}\n")

	return &ExportResponse{
		Filename: sanitizeFilename(suite.Name) + ".k6.js",
		Content:  sb.String(),
		MimeType: "application/javascript",
	}, nil
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(" ", "_", "/", "_", "\\", "_", ":", "_")
	return strings.ToLower(replacer.Replace(name))
}

func sanitizeVarName(name string) string {
	replacer := strings.NewReplacer(" ", "_", "/", "_", "-", "_", ".", "_")
	result := strings.ToLower(replacer.Replace(name))
	// Remove non-alphanumeric characters
	var sb strings.Builder
	for _, c := range result {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}
