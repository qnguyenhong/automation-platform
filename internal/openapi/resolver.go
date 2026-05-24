package openapi

import (
	"fmt"
	"strings"
)

// Resolver resolves $ref references in OpenAPI specifications.
type Resolver struct {
	spec *Spec
}

// NewResolver creates a new reference resolver for the given spec.
func NewResolver(spec *Spec) *Resolver {
	return &Resolver{spec: spec}
}

// ResolveSchema resolves $ref references in a Schema recursively.
func (r *Resolver) ResolveSchema(schema Schema) Schema {
	if schema.Ref != "" {
		resolved := r.resolveSchemaRef(schema.Ref)
		if resolved != nil {
			return *resolved
		}
		return schema
	}

	// Resolve properties
	if schema.Properties != nil {
		resolved := make(map[string]Schema, len(schema.Properties))
		for name, prop := range schema.Properties {
			resolved[name] = r.ResolveSchema(prop)
		}
		schema.Properties = resolved
	}

	// Resolve items
	if schema.Items != nil {
		resolvedItems := r.ResolveSchema(*schema.Items)
		schema.Items = &resolvedItems
	}

	// Resolve allOf
	if len(schema.AllOf) > 0 {
		resolved := make([]Schema, len(schema.AllOf))
		for i, s := range schema.AllOf {
			resolved[i] = r.ResolveSchema(s)
		}
		schema.AllOf = resolved

		// Merge allOf into a single schema
		schema = r.mergeAllOf(schema)
	}

	// Resolve oneOf
	if len(schema.OneOf) > 0 {
		resolved := make([]Schema, len(schema.OneOf))
		for i, s := range schema.OneOf {
			resolved[i] = r.ResolveSchema(s)
		}
		schema.OneOf = resolved
	}

	// Resolve anyOf
	if len(schema.AnyOf) > 0 {
		resolved := make([]Schema, len(schema.AnyOf))
		for i, s := range schema.AnyOf {
			resolved[i] = r.ResolveSchema(s)
		}
		schema.AnyOf = resolved
	}

	return schema
}

// resolveSchemaRef looks up a $ref path and returns the referenced Schema.
func (r *Resolver) resolveSchemaRef(ref string) *Schema {
	if r.spec == nil || r.spec.Components == nil {
		return nil
	}

	// Handle #/components/schemas/Name
	prefix := "#/components/schemas/"
	if strings.HasPrefix(ref, prefix) {
		name := strings.TrimPrefix(ref, prefix)
		if schema, ok := r.spec.Components.Schemas[name]; ok {
			// Recursively resolve the referenced schema
			resolved := r.ResolveSchema(schema)
			return &resolved
		}
	}

	return nil
}

// ResolveParameter resolves $ref in a Parameter.
func (r *Resolver) ResolveParameter(param Parameter) Parameter {
	param.Schema = r.ResolveSchema(param.Schema)
	return param
}

// ResolveRequestBody resolves $ref in a RequestBody.
func (r *Resolver) ResolveRequestBody(rb *RequestBody) *RequestBody {
	if rb == nil {
		return nil
	}
	for contentType, mt := range rb.Content {
		mt.Schema = r.ResolveSchema(mt.Schema)
		rb.Content[contentType] = mt
	}
	return rb
}

// ResolveResponses resolves $ref in responses.
func (r *Resolver) ResolveResponses(responses map[string]Response) map[string]Response {
	if responses == nil {
		return nil
	}
	resolved := make(map[string]Response, len(responses))
	for code, resp := range responses {
		for contentType, mt := range resp.Content {
			mt.Schema = r.ResolveSchema(mt.Schema)
			resp.Content[contentType] = mt
		}
		resolved[code] = resp
	}
	return resolved
}

// mergeAllOf merges allOf schemas into one combined schema.
func (r *Resolver) mergeAllOf(schema Schema) Schema {
	merged := Schema{
		Type:       "object",
		Properties: make(map[string]Schema),
	}

	// Copy non-allOf fields
	if schema.Type != "" {
		merged.Type = schema.Type
	}
	if schema.Description != "" {
		merged.Description = schema.Description
	}

	for _, s := range schema.AllOf {
		if s.Type != "" {
			merged.Type = s.Type
		}
		for name, prop := range s.Properties {
			merged.Properties[name] = prop
		}
		merged.Required = append(merged.Required, s.Required...)
	}

	// Also merge direct properties
	for name, prop := range schema.Properties {
		merged.Properties[name] = prop
	}
	merged.Required = append(merged.Required, schema.Required...)

	return merged
}

// ExtractSecuritySchemes extracts and enriches security schemes with header templates.
func ExtractSecuritySchemes(spec *Spec) []SecurityScheme {
	if spec.Components == nil || spec.Components.SecuritySchemes == nil {
		return nil
	}

	var schemes []SecurityScheme
	for name, scheme := range spec.Components.SecuritySchemes {
		enriched := scheme
		enriched.Name = name
		enriched.HeaderTemplate = generateHeaderTemplate(scheme)
		schemes = append(schemes, enriched)
	}
	return schemes
}

// generateHeaderTemplate creates a template string for the auth header.
func generateHeaderTemplate(scheme SecurityScheme) string {
	switch scheme.Type {
	case "http":
		switch scheme.Scheme {
		case "bearer":
			return fmt.Sprintf("Authorization: Bearer {{%s_token}}", strings.ToLower(scheme.Name))
		case "basic":
			return fmt.Sprintf("Authorization: Basic {{%s_credentials}}", strings.ToLower(scheme.Name))
		}
	case "apiKey":
		if scheme.In == "header" {
			paramName := scheme.ParameterName
			if paramName == "" {
				paramName = scheme.Name
			}
			return fmt.Sprintf("%s: {{%s_key}}", paramName, strings.ToLower(scheme.Name))
		}
	case "oauth2":
		return fmt.Sprintf("Authorization: Bearer {{%s_access_token}}", strings.ToLower(scheme.Name))
	}
	return ""
}

// GenerateNegativeTests creates negative test scenarios from an endpoint's spec.
func GenerateNegativeTests(ep Endpoint, spec *Spec) []NegativeTestCase {
	var tests []NegativeTestCase

	// Test missing required parameters
	for _, param := range ep.Parameters {
		if param.Required {
			tests = append(tests, NegativeTestCase{
				Name:           fmt.Sprintf("%s %s - Missing required param: %s", ep.Method, ep.Path, param.Name),
				Description:    fmt.Sprintf("Should return 400 when required %s parameter '%s' is missing", param.In, param.Name),
				OmitParam:      param.Name,
				ExpectedStatus: 400,
			})
		}
	}

	// Test missing required body fields
	if ep.RequestBody != nil {
		for _, mt := range ep.RequestBody.Content {
			resolver := NewResolver(spec)
			resolved := resolver.ResolveSchema(mt.Schema)
			for _, field := range resolved.Required {
				tests = append(tests, NegativeTestCase{
					Name:           fmt.Sprintf("%s %s - Missing required field: %s", ep.Method, ep.Path, field),
					Description:    fmt.Sprintf("Should return 400/422 when required body field '%s' is missing", field),
					OmitBodyField:  field,
					ExpectedStatus: 400,
				})
			}

			// Test invalid enum values
			for propName, prop := range resolved.Properties {
				if len(prop.Enum) > 0 {
					tests = append(tests, NegativeTestCase{
						Name:           fmt.Sprintf("%s %s - Invalid enum value for: %s", ep.Method, ep.Path, propName),
						Description:    fmt.Sprintf("Should return 400/422 when field '%s' has invalid enum value", propName),
						InvalidField:   propName,
						InvalidValue:   "__INVALID_ENUM_VALUE__",
						ExpectedStatus: 400,
					})
				}

				// Test invalid type
				if prop.Type == "integer" || prop.Type == "number" {
					tests = append(tests, NegativeTestCase{
						Name:           fmt.Sprintf("%s %s - Invalid type for: %s", ep.Method, ep.Path, propName),
						Description:    fmt.Sprintf("Should return 400/422 when field '%s' receives string instead of %s", propName, prop.Type),
						InvalidField:   propName,
						InvalidValue:   "not_a_number",
						ExpectedStatus: 400,
					})
				}
			}
		}
	}

	// Test unauthorized (if security defined)
	if spec.Security != nil && len(spec.Security) > 0 {
		tests = append(tests, NegativeTestCase{
			Name:           fmt.Sprintf("%s %s - Unauthorized (no auth)", ep.Method, ep.Path),
			Description:    "Should return 401 when no authentication provided",
			OmitAuth:       true,
			ExpectedStatus: 401,
		})
	}

	// Test not found with invalid ID (if path has ID parameter)
	if strings.Contains(ep.Path, "{") {
		tests = append(tests, NegativeTestCase{
			Name:           fmt.Sprintf("%s %s - Not found (invalid ID)", ep.Method, ep.Path),
			Description:    "Should return 404 when resource ID doesn't exist",
			UseInvalidID:   true,
			ExpectedStatus: 404,
		})
	}

	return tests
}

// NegativeTestCase represents a generated negative test case.
type NegativeTestCase struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	OmitParam      string `json:"omit_param,omitempty"`
	OmitBodyField  string `json:"omit_body_field,omitempty"`
	InvalidField   string `json:"invalid_field,omitempty"`
	InvalidValue   any    `json:"invalid_value,omitempty"`
	OmitAuth       bool   `json:"omit_auth,omitempty"`
	UseInvalidID   bool   `json:"use_invalid_id,omitempty"`
	ExpectedStatus int    `json:"expected_status"`
}
