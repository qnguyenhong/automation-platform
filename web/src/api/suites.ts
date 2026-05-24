import { apiGet, apiPost, apiPut, apiDelete } from './client'

export interface TestSuite {
  id: string
  project_id: string
  name: string
  description?: string
  test_type: 'api' | 'e2e' | 'load' | 'unit'
  schedule_cron?: string
  config: Record<string, any>
  tags: string[]
  created_at: string
  updated_at: string
}

export interface TestCase {
  id: string
  suite_id: string
  name: string
  description?: string
  config: Record<string, any>
  tags: string[]
  sort_order: number
  enabled: boolean
  created_at: string
  updated_at: string
}

export const suitesApi = {
  list: (projectId: string) => apiGet<TestSuite[]>(`/projects/${projectId}/suites`),
  get: (projectId: string, suiteId: string) => apiGet<TestSuite>(`/projects/${projectId}/suites/${suiteId}`),
  create: (projectId: string, data: Partial<TestSuite>) => apiPost<TestSuite>(`/projects/${projectId}/suites`, data),
  update: (projectId: string, suiteId: string, data: Partial<TestSuite>) => apiPut<TestSuite>(`/projects/${projectId}/suites/${suiteId}`, data),
  delete: (projectId: string, suiteId: string) => apiDelete(`/projects/${projectId}/suites/${suiteId}`),
}

export const casesApi = {
  list: (projectId: string, suiteId: string) => apiGet<TestCase[]>(`/projects/${projectId}/suites/${suiteId}/cases`),
  get: (projectId: string, suiteId: string, caseId: string) => apiGet<TestCase>(`/projects/${projectId}/suites/${suiteId}/cases/${caseId}`),
  create: (projectId: string, suiteId: string, data: Partial<TestCase>) => apiPost<TestCase>(`/projects/${projectId}/suites/${suiteId}/cases`, data),
  update: (projectId: string, suiteId: string, caseId: string, data: Partial<TestCase>) => apiPut<TestCase>(`/projects/${projectId}/suites/${suiteId}/cases/${caseId}`, data),
  delete: (projectId: string, suiteId: string, caseId: string) => apiDelete(`/projects/${projectId}/suites/${suiteId}/cases/${caseId}`),
}
