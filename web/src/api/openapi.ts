import { apiPost } from './client'
import type {
  ParseOpenAPIResponse,
  ImportOpenAPIRequest,
  ImportOpenAPIResponse,
} from '@/types/openapi'

export async function parseSpec(content: string, url?: string): Promise<ParseOpenAPIResponse> {
  const { data } = await apiPost<{ data: ParseOpenAPIResponse }>('/openapi/parse', { content, url })
  return data.data
}

export async function importEndpoints(request: ImportOpenAPIRequest): Promise<ImportOpenAPIResponse> {
  const { data } = await apiPost<{ data: ImportOpenAPIResponse }>('/openapi/import', request)
  return data.data
}
