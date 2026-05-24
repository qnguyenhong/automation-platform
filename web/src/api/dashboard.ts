import { apiGet } from './client'

export interface DashboardSummary {
  total_runs: number
  total_passed: number
  total_failed: number
  active_workers: number
  pass_rate: number
}

export interface TrendPoint {
  date: string
  passed: number
  failed: number
  avg_duration_ms: number
}

export interface FlakyTest {
  case_id: string
  case_name: string
  total_runs: number
  passed: number
  failed: number
}

export const dashboardApi = {
  summary: (projectId?: string) =>
    apiGet<DashboardSummary>(`/dashboard/summary${projectId ? `?project_id=${projectId}` : ''}`),
  trends: (projectId?: string, days = 30) =>
    apiGet<TrendPoint[]>(`/dashboard/trends?days=${days}${projectId ? `&project_id=${projectId}` : ''}`),
  flaky: (projectId?: string, limit = 20) =>
    apiGet<FlakyTest[]>(`/dashboard/flaky?limit=${limit}${projectId ? `&project_id=${projectId}` : ''}`),
}
