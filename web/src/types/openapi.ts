export interface ParsedEndpoint {
  method: string
  path: string
  summary: string
  description: string
  tags: string[]
  parameters: EndpointParameter[]
  request_body: RequestBodySchema | null
  responses: Record<string, ResponseSchema>
}

export interface EndpointParameter {
  name: string
  in: 'query' | 'header' | 'path' | 'cookie'
  required: boolean
  schema: ParameterSchema
  example: unknown
}

export interface ParameterSchema {
  type: string
  format?: string
  enum?: unknown[]
}

export interface RequestBodySchema {
  required: boolean
  content: Record<string, { schema: SchemaObject; example?: unknown }>
}

export interface ResponseSchema {
  description: string
  content?: Record<string, { schema: SchemaObject }>
}

export interface SchemaObject {
  type?: string
  properties?: Record<string, SchemaObject>
  items?: SchemaObject
  example?: unknown
  enum?: unknown[]
  required?: string[]
}

export interface ParsedSpec {
  title: string
  version: string
  base_url: string
  endpoints: ParsedEndpoint[]
}

export interface ParseOpenAPIResponse {
  spec: {
    title: string
    version: string
    description: string
  }
  endpoints: ParsedEndpoint[]
}

export interface ImportEndpointConfig {
  method: string
  path: string
  headers: Record<string, string>
  params: Record<string, string>
  body: Record<string, unknown> | unknown
  assertions: AssertionConfig[]
}

export interface AssertionConfig {
  type: 'status' | 'response_time' | 'header' | 'json_path'
  path?: string
  expected: unknown
  capture_as?: string
}

export interface ImportOpenAPIRequest {
  project_id: string
  suite_name: string
  endpoints: ImportEndpointConfig[]
}

export interface ImportOpenAPIResponse {
  suite_id: string
  cases_created: number
}

export interface VariableDefinition {
  name: string
  description: string
  example: string
  category: 'random' | 'time' | 'env' | 'captured'
}

export const BUILT_IN_VARIABLES: VariableDefinition[] = [
  { name: '$uuid', description: 'Random UUID v4', example: '550e8400-e29b-41d4-a716-446655440000', category: 'random' },
  { name: '$random_email', description: 'Random email address', example: 'a3k9@example.com', category: 'random' },
  { name: '$random_string(8)', description: 'Random alphanumeric string', example: 'aB3xK9mQ', category: 'random' },
  { name: '$random_int(1,100)', description: 'Random integer in range', example: '47', category: 'random' },
  { name: '$random_float(0,100,2)', description: 'Random float', example: '3.14', category: 'random' },
  { name: '$random_bool', description: 'Random boolean', example: 'true', category: 'random' },
  { name: '$random_name', description: 'Random full name', example: 'John Smith', category: 'random' },
  { name: '$random_phone', description: 'Random phone number', example: '+1-555-0123', category: 'random' },
  { name: '$random_ipv4', description: 'Random IPv4 address', example: '192.168.1.42', category: 'random' },
  { name: '$random_enum(a,b,c)', description: 'Random pick from values', example: 'b', category: 'random' },
  { name: '$random_hex(16)', description: 'Random hex string', example: 'a3f5b2c1d4e6f789', category: 'random' },
  { name: '$random_url', description: 'Random URL', example: 'https://example.com/path', category: 'random' },
  { name: '$timestamp', description: 'Unix timestamp (seconds)', example: '1716500000', category: 'time' },
  { name: '$timestamp_ms', description: 'Unix timestamp (milliseconds)', example: '1716500000000', category: 'time' },
  { name: '$iso_date', description: 'ISO 861 datetime', example: '2026-05-24T12:00:00Z', category: 'time' },
  { name: '$date(2006-01-02)', description: 'Custom date format', example: '2026-05-24', category: 'time' },
  { name: '$env.API_KEY', description: 'Environment variable', example: 'your-api-key', category: 'env' },
  { name: '$seq(name)', description: 'Auto-incrementing sequence', example: '1', category: 'random' },
]
