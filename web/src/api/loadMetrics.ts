import { apiGet } from './client'

export interface TimePoint {
  second: number
  value: number
}

export interface TimeSeriesData {
  rps: TimePoint[]
  p95: TimePoint[]
}

export interface LoadMetrics {
  id: string
  result_id: string
  run_id: string
  total_requests: number
  success_count: number
  error_count: number
  error_rate: number
  throughput_rps: number
  min_latency_ms: number
  max_latency_ms: number
  avg_latency_ms: number
  p50_latency_ms: number
  p90_latency_ms: number
  p95_latency_ms: number
  p99_latency_ms: number
  total_bytes: number
  time_series: TimeSeriesData
  status_codes: Record<string, number>
  created_at: string
}

export const loadMetricsApi = {
  getByRun: (runId: string) => apiGet<LoadMetrics>(`/runs/${runId}/load-metrics`),
  getByResult: (runId: string, resultId: string) =>
    apiGet<LoadMetrics>(`/runs/${runId}/results/${resultId}/load-metrics`),
}
