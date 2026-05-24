import axios from 'axios'
import { useAuthStore } from '../store/auth'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

const client = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

client.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default client

export interface ApiResponse<T> {
  data?: T
  error?: string
  meta?: any
}

export async function apiGet<T>(url: string): Promise<T> {
  const response = await client.get<ApiResponse<T>>(url)
  return response.data.data as T
}

export async function apiPost<T>(url: string, data?: any): Promise<T> {
  const response = await client.post<ApiResponse<T>>(url, data)
  return response.data.data as T
}

export async function apiPut<T>(url: string, data?: any): Promise<T> {
  const response = await client.put<ApiResponse<T>>(url, data)
  return response.data.data as T
}

export async function apiDelete(url: string): Promise<void> {
  await client.delete(url)
}
