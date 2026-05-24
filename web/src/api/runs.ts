import { apiGet, apiPost } from './client'

export interface TestRun {
  id: string
  suite_id: string
  project_id: string
  status: 'pending' | 'running' | 'passed' | 'failed' | 'cancelled' | 'error'
  trigger: 'manual' | 'scheduled' | 'api' | 'ci'
  triggered_by?: string
  total_cases: number
  passed: number
  failed: number
  skipped: number
  errored: number
  started_at?: string
  finished_at?: string
  duration_ms?: number
  metadata: Record<string, any>
  created_at: string
}

export interface TestResult {
  id: string
  run_id: string
  case_id: string
  worker_id?: string
  status: 'pending' | 'running' | 'passed' | 'failed' | 'skipped' | 'error'
  error_message?: string
  assertions: any[]
  request_data?: Record<string, any>
  response_data?: Record<string, any>
  artifacts: any[]
  metrics?: Record<string, any>
  stdout?: string
  stderr?: string
  duration_ms?: number
  started_at?: string
  finished_at?: string
  retry_of?: string
  created_at: string
}

export const runsApi = {
  list: (projectId: string, page = 1, perPage = 20) =>
    apiGet<TestRun[]>(`/projects/${projectId}/runs?page=${page}&per_page=${perPage}`),
  get: (runId: string) => apiGet<TestRun>(`/runs/${runId}`),
  trigger: (projectId: string, suiteId: string, data?: { trigger?: string; metadata?: any }) =>
    apiPost<TestRun>(`/projects/${projectId}/suites/${suiteId}/runs`, data),
  cancel: (runId: string) => apiPost(`/runs/${runId}/cancel`),
  results: (runId: string) => apiGet<TestResult[]>(`/runs/${runId}/results`),
  result: (runId: string, resultId: string) => apiGet<TestResult>(`/runs/${runId}/results/${resultId}`),
}
