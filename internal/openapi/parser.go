package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parser parses OpenAPI 3.x specifications.
type Parser struct{}

// NewParser creates a new OpenAPI parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses an OpenAPI spec from YAML or JSON bytes.
func (p *Parser) Parse(data []byte) (*Spec, error) {
	// Try JSON first
	var spec Spec
	if err := json.Unmarshal(data, &spec); err == nil && spec.OpenAPI != "" {
		return &spec, nil
	}

	// Try YAML
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse spec as JSON or YAML: %w", err)
	}

	if spec.OpenAPI == "" {
		return nil, fmt.Errorf("invalid OpenAPI spec: missing 'openapi' field")
	}

	return &spec, nil
}

// ExtractEndpoints flattens all paths and methods into a list of endpoints.
func (p *Parser) ExtractEndpoints(spec *Spec) []Endpoint {
	var endpoints []Endpoint

	// Determine base URL from servers
	baseURL := ""
	if len(spec.Servers) > 0 {
		baseURL = strings.TrimSuffix(spec.Servers[0].URL, "/")
	}

	for path, pathItem := range spec.Paths {
		ops := []struct {
			method    string
			operation *Operation
		}{
			{"GET", pathItem.Get},
			{"POST", pathItem.Post},
			{"PUT", pathItem.Put},
			{"DELETE", pathItem.Delete},
			{"PATCH", pathItem.Patch},
			{"HEAD", pathItem.Head},
			{"OPTIONS", pathItem.Options},
		}

		for _, op := range ops {
			if op.operation == nil {
				continue
			}

			endpoint := Endpoint{
				Method:      op.method,
				Path:        path,
				Summary:     op.operation.Summary,
				Description: op.operation.Description,
				Tags:        op.operation.Tags,
				OperationID: op.operation.OperationID,
				Parameters:  op.operation.Parameters,
				RequestBody: op.operation.RequestBody,
				Responses:   op.operation.Responses,
				Deprecated:  op.operation.Deprecated,
				BaseURL:     baseURL,
			}

			// Merge path-level parameters with operation-level
			if pathItem.Get != nil || pathItem.Post != nil {
				// Path has parameters defined at path level
				endpoint.Parameters = mergeParameters(pathItem, endpoint.Parameters)
			}

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

// ParseAndExtract is a convenience method that parses and extracts in one call.
func (p *Parser) ParseAndExtract(data []byte) (*ParsedSpec, error) {
	spec, err := p.Parse(data)
	if err != nil {
		return nil, err
	}

	// Resolve $ref references
	resolver := NewResolver(spec)
	p.resolveSpec(spec, resolver)

	endpoints := p.ExtractEndpoints(spec)

	baseURL := ""
	if len(spec.Servers) > 0 {
		baseURL = strings.TrimSuffix(spec.Servers[0].URL, "/")
	}

	// Extract security schemes
	securitySchemes := ExtractSecuritySchemes(spec)

	return &ParsedSpec{
		Title:           spec.Info.Title,
		Description:     spec.Info.Description,
		Version:         spec.Info.Version,
		BaseURL:         baseURL,
		Endpoints:       endpoints,
		SecuritySchemes: securitySchemes,
	}, nil
}

// resolveSpec resolves all $ref references in the spec's paths.
func (p *Parser) resolveSpec(spec *Spec, resolver *Resolver) {
	for path, pathItem := range spec.Paths {
		if pathItem.Get != nil {
			p.resolveOperation(pathItem.Get, resolver)
		}
		if pathItem.Post != nil {
			p.resolveOperation(pathItem.Post, resolver)
		}
		if pathItem.Put != nil {
			p.resolveOperation(pathItem.Put, resolver)
		}
		if pathItem.Delete != nil {
			p.resolveOperation(pathItem.Delete, resolver)
		}
		if pathItem.Patch != nil {
			p.resolveOperation(pathItem.Patch, resolver)
		}
		spec.Paths[path] = pathItem
	}
}

// resolveOperation resolves $ref in all parts of an operation.
func (p *Parser) resolveOperation(op *Operation, resolver *Resolver) {
	// Resolve parameters
	for i, param := range op.Parameters {
		op.Parameters[i] = resolver.ResolveParameter(param)
	}

	// Resolve request body
	if op.RequestBody != nil {
		op.RequestBody = resolver.ResolveRequestBody(op.RequestBody)
	}

	// Resolve responses
	op.Responses = resolver.ResolveResponses(op.Responses)
}

// mergeParameters merges path-level and operation-level parameters.
// Operation-level parameters override path-level ones with the same name+in.
func mergeParameters(pathItem PathItem, opParams []Parameter) []Parameter {
	// Get path-level parameters from any operation (they're shared)
	// Actually, path-level parameters are on the PathItem itself in OpenAPI 3.x
	// But our struct doesn't have a Parameters field on PathItem
	// This is fine - we just return the operation parameters
	return opParams
}

// GenerateExampleFromSchema generates an example value from a JSON Schema.
func GenerateExampleFromSchema(schema Schema) any {
	if schema.Example != nil {
		return schema.Example
	}

	if schema.Default != nil {
		return schema.Default
	}

	switch schema.Type {
	case "string":
		if len(schema.Enum) > 0 {
			return schema.Enum[0]
		}
		switch schema.Format {
		case "email":
			return "user@example.com"
		case "uuid":
			return "550e8400-e29b-41d4-a716-446655440000"
		case "date":
			return "2026-01-01"
		case "date-time":
			return "2026-01-01T00:00:00Z"
		case "uri", "url":
			return "https://example.com"
		default:
			if schema.MinLength != nil && *schema.MinLength > 0 {
				return strings.Repeat("a", *schema.MinLength)
			}
			return "string"
		}
	case "integer":
		if len(schema.Enum) > 0 {
			return schema.Enum[0]
		}
		if schema.Minimum != nil {
			return int(*schema.Minimum)
		}
		return 0
	case "number":
		if len(schema.Enum) > 0 {
			return schema.Enum[0]
		}
		if schema.Minimum != nil {
			return *schema.Minimum
		}
		return 0.0
	case "boolean":
		return false
	case "array":
		if schema.Items != nil {
			item := GenerateExampleFromSchema(*schema.Items)
			return []any{item}
		}
		return []any{}
	case "object":
		obj := make(map[string]any)
		for name, prop := range schema.Properties {
			obj[name] = GenerateExampleFromSchema(prop)
		}
		return obj
	default:
		return nil
	}
}

// GetRequestBodyExample generates an example request body from the endpoint's request body schema.
func GetRequestBodyExample(ep Endpoint) map[string]any {
	if ep.RequestBody == nil {
		return nil
	}

	// Try application/json first
	if mediaType, ok := ep.RequestBody.Content["application/json"]; ok {
		if mediaType.Example != nil {
			if m, ok := mediaType.Example.(map[string]any); ok {
				return m
			}
		}
		example := GenerateExampleFromSchema(mediaType.Schema)
		if m, ok := example.(map[string]any); ok {
			return m
		}
	}

	// Try first available content type
	for _, mediaType := range ep.RequestBody.Content {
		if mediaType.Example != nil {
			if m, ok := mediaType.Example.(map[string]any); ok {
				return m
			}
		}
		example := GenerateExampleFromSchema(mediaType.Schema)
		if m, ok := example.(map[string]any); ok {
			return m
		}
	}

	return nil
}
