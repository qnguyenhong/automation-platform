import { apiGet, apiDelete } from './client'

export interface Worker {
  id: string
  name: string
  hostname: string
  ip_address?: string
  executor_types: string[]
  status: 'online' | 'offline' | 'busy' | 'draining'
  current_job_id?: string
  max_concurrent: number
  version?: string
  labels: Record<string, any>
  last_heartbeat?: string
  registered_at: string
  updated_at: string
}

export const workersApi = {
  list: () => apiGet<Worker[]>('/workers'),
  get: (id: string) => apiGet<Worker>(`/workers/${id}`),
  delete: (id: string) => apiDelete(`/workers/${id}`),
}
