package openapi

// Spec represents a minimal OpenAPI 3.x specification.
type Spec struct {
	OpenAPI    string                 `json:"openapi" yaml:"openapi"`
	Info       Info                   `json:"info" yaml:"info"`
	Servers    []Server               `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths      map[string]PathItem    `json:"paths" yaml:"paths"`
	Tags       []Tag                  `json:"tags,omitempty" yaml:"tags,omitempty"`
	Components *Components            `json:"components,omitempty" yaml:"components,omitempty"`
	Security   []map[string][]string  `json:"security,omitempty" yaml:"security,omitempty"`
}

// Info contains metadata about the API.
type Info struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string `json:"version" yaml:"version"`
}

// Server represents a server endpoint.
type Server struct {
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// Tag represents a grouping tag.
type Tag struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// PathItem describes the operations available on a single path.
type PathItem struct {
	Get    *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Post   *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Put    *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	Patch  *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	Head   *Operation `json:"head,omitempty" yaml:"head,omitempty"`
	Options *Operation `json:"options,omitempty" yaml:"options,omitempty"`
}

// Operation describes a single API operation.
type Operation struct {
	OperationID string              `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Summary     string              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        []string            `json:"tags,omitempty" yaml:"tags,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty" yaml:"responses,omitempty"`
	Deprecated  bool                `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name        string  `json:"name" yaml:"name"`
	In          string  `json:"in" yaml:"in"` // "query", "header", "path", "cookie"
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool    `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      Schema  `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example     any     `json:"example,omitempty" yaml:"example,omitempty"`
}

// RequestBody describes the request body.
type RequestBody struct {
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool                 `json:"required,omitempty" yaml:"required,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// MediaType describes a media type.
type MediaType struct {
	Schema  Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example any    `json:"example,omitempty" yaml:"example,omitempty"`
}

// Response describes a single response.
type Response struct {
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// Schema describes the schema of a parameter or request/response body.
type Schema struct {
	Type                 string            `json:"type,omitempty" yaml:"type,omitempty"`
	Format               string            `json:"format,omitempty" yaml:"format,omitempty"`
	Description          string            `json:"description,omitempty" yaml:"description,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items                *Schema           `json:"items,omitempty" yaml:"items,omitempty"`
	Required             []string          `json:"required,omitempty" yaml:"required,omitempty"`
	Example              any               `json:"example,omitempty" yaml:"example,omitempty"`
	Enum                 []any             `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default              any               `json:"default,omitempty" yaml:"default,omitempty"`
	Minimum              *float64          `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum              *float64          `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	MinLength            *int              `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength            *int              `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	Pattern              string            `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	AdditionalProperties any               `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	AllOf                []Schema          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	OneOf                []Schema          `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf                []Schema          `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Ref                  string            `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// Endpoint is a flattened representation of an API endpoint for import.
type Endpoint struct {
	Method      string              `json:"method"`
	Path        string              `json:"path"`
	Summary     string              `json:"summary"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"`
	OperationID string              `json:"operation_id"`
	Parameters  []Parameter         `json:"parameters"`
	RequestBody *RequestBody        `json:"request_body,omitempty"`
	Responses   map[string]Response `json:"responses"`
	Deprecated  bool                `json:"deprecated"`
	BaseURL     string              `json:"base_url"`
}

// ParsedSpec is the result of parsing an OpenAPI spec.
type ParsedSpec struct {
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	Version         string           `json:"version"`
	BaseURL         string           `json:"base_url"`
	Endpoints       []Endpoint       `json:"endpoints"`
	SecuritySchemes []SecurityScheme `json:"security_schemes,omitempty"`
}

// Components holds reusable objects for the specification.
type Components struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Responses       map[string]Response       `json:"responses,omitempty" yaml:"responses,omitempty"`
}

// SecurityScheme describes an authentication method.
type SecurityScheme struct {
	Name             string `json:"name" yaml:"name"`
	Type             string `json:"type" yaml:"type"` // apiKey, http, oauth2, openIdConnect
	Scheme           string `json:"scheme,omitempty" yaml:"scheme,omitempty"` // bearer, basic
	BearerFormat     string `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	In               string `json:"in,omitempty" yaml:"in,omitempty"` // header, query, cookie
	ParameterName    string `json:"parameterName,omitempty" yaml:"parameterName,omitempty"`
	Description      string `json:"description,omitempty" yaml:"description,omitempty"`
	HeaderTemplate   string `json:"header_template,omitempty"` // Generated: e.g., "Authorization: Bearer {{token}}"
}
