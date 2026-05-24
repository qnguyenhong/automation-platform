package openapi

import (
	"testing"
)

const petstoreSpec = `
openapi: "3.0.3"
info:
  title: Petstore API
  description: A sample API
  version: "1.0.0"
servers:
  - url: https://petstore.example.com/api/v1
tags:
  - name: pets
    description: Pet operations
  - name: store
    description: Store operations
paths:
  /pets:
    get:
      summary: List all pets
      tags: [pets]
      operationId: listPets
      parameters:
        - name: limit
          in: query
          required: false
          schema:
            type: integer
            format: int32
        - name: offset
          in: query
          schema:
            type: integer
      responses:
        "200":
          description: A list of pets
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: integer
                    name:
                      type: string
    post:
      summary: Create a pet
      tags: [pets]
      operationId: createPet
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [name]
              properties:
                name:
                  type: string
                  example: "Buddy"
                tag:
                  type: string
                  example: "dog"
      responses:
        "201":
          description: Pet created
  /pets/{petId}:
    get:
      summary: Get a pet by ID
      tags: [pets]
      operationId: getPet
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: integer
      responses:
        "200":
          description: A pet
    delete:
      summary: Delete a pet
      tags: [pets]
      operationId: deletePet
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: integer
      responses:
        "204":
          description: Pet deleted
`

func TestParserParse(t *testing.T) {
	parser := NewParser()

	spec, err := parser.Parse([]byte(petstoreSpec))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if spec.OpenAPI != "3.0.3" {
		t.Errorf("OpenAPI = %q, want %q", spec.OpenAPI, "3.0.3")
	}
	if spec.Info.Title != "Petstore API" {
		t.Errorf("Title = %q, want %q", spec.Info.Title, "Petstore API")
	}
	if len(spec.Paths) != 2 {
		t.Errorf("Paths count = %d, want 2", len(spec.Paths))
	}
}

func TestParserExtractEndpoints(t *testing.T) {
	parser := NewParser()

	spec, err := parser.Parse([]byte(petstoreSpec))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	endpoints := parser.ExtractEndpoints(spec)

	if len(endpoints) != 4 {
		t.Fatalf("Endpoints count = %d, want 4", len(endpoints))
	}

	// Check methods are present
	methods := make(map[string]bool)
	for _, ep := range endpoints {
		key := ep.Method + " " + ep.Path
		methods[key] = true
	}

	expected := []string{"GET /pets", "POST /pets", "GET /pets/{petId}", "DELETE /pets/{petId}"}
	for _, exp := range expected {
		if !methods[exp] {
			t.Errorf("Missing endpoint: %s", exp)
		}
	}

	// Check base URL
	for _, ep := range endpoints {
		if ep.BaseURL != "https://petstore.example.com/api/v1" {
			t.Errorf("BaseURL = %q, want %q", ep.BaseURL, "https://petstore.example.com/api/v1")
		}
	}
}

func TestParserParseAndExtract(t *testing.T) {
	parser := NewParser()

	parsed, err := parser.ParseAndExtract([]byte(petstoreSpec))
	if err != nil {
		t.Fatalf("ParseAndExtract error: %v", err)
	}

	if parsed.Title != "Petstore API" {
		t.Errorf("Title = %q, want %q", parsed.Title, "Petstore API")
	}
	if parsed.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", parsed.Version, "1.0.0")
	}
	if parsed.BaseURL != "https://petstore.example.com/api/v1" {
		t.Errorf("BaseURL = %q, want %q", parsed.BaseURL, "https://petstore.example.com/api/v1")
	}
	if len(parsed.Endpoints) != 4 {
		t.Errorf("Endpoints count = %d, want 4", len(parsed.Endpoints))
	}
}

func TestParserParseJSON(t *testing.T) {
	parser := NewParser()

	jsonSpec := `{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0"},
		"paths": {
			"/test": {
				"get": {
					"summary": "Test endpoint",
					"responses": {"200": {"description": "OK"}}
				}
			}
		}
	}`

	spec, err := parser.Parse([]byte(jsonSpec))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if spec.OpenAPI != "3.0.0" {
		t.Errorf("OpenAPI = %q, want %q", spec.OpenAPI, "3.0.0")
	}
}

func TestGenerateExampleFromSchema(t *testing.T) {
	tests := []struct {
		name   string
		schema Schema
		check  func(any) bool
	}{
		{
			name:   "string type",
			schema: Schema{Type: "string"},
			check:  func(v any) bool { _, ok := v.(string); return ok },
		},
		{
			name:   "integer type",
			schema: Schema{Type: "integer"},
			check:  func(v any) bool { _, ok := v.(int); return ok },
		},
		{
			name:   "boolean type",
			schema: Schema{Type: "boolean"},
			check:  func(v any) bool { _, ok := v.(bool); return ok },
		},
		{
			name: "object with properties",
			schema: Schema{
				Type: "object",
				Properties: map[string]Schema{
					"name": {Type: "string"},
					"age":  {Type: "integer"},
				},
			},
			check: func(v any) bool {
				m, ok := v.(map[string]any)
				return ok && len(m) == 2
			},
		},
		{
			name:   "array",
			schema: Schema{Type: "array", Items: &Schema{Type: "string"}},
			check: func(v any) bool {
				arr, ok := v.([]any)
				return ok && len(arr) == 1
			},
		},
		{
			name:   "email format",
			schema: Schema{Type: "string", Format: "email"},
			check:  func(v any) bool { return v == "user@example.com" },
		},
		{
			name:   "uuid format",
			schema: Schema{Type: "string", Format: "uuid"},
			check:  func(v any) bool { s, ok := v.(string); return ok && len(s) == 36 },
		},
		{
			name:   "enum",
			schema: Schema{Type: "string", Enum: []any{"active", "inactive"}},
			check:  func(v any) bool { return v == "active" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateExampleFromSchema(tt.schema)
			if !tt.check(result) {
				t.Errorf("GenerateExampleFromSchema check failed for %s, got: %v", tt.name, result)
			}
		})
	}
}
