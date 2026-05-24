import { apiGet, apiPost, apiPut, apiDelete } from './client'

export interface Project {
  id: string
  name: string
  description?: string
  created_by?: string
  created_at: string
  updated_at: string
}

export const projectsApi = {
  list: () => apiGet<Project[]>('/projects'),
  get: (id: string) => apiGet<Project>(`/projects/${id}`),
  create: (data: { name: string; description?: string }) => apiPost<Project>('/projects', data),
  update: (id: string, data: { name: string; description?: string }) => apiPut<Project>(`/projects/${id}`, data),
  delete: (id: string) => apiDelete(`/projects/${id}`),
}
