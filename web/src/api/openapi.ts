import { apiPost } from './client'
import type {
  ParseOpenAPIResponse,
  ImportOpenAPIRequest,
  ImportOpenAPIResponse,
} from '@/types/openapi'

export async function parseSpec(content: string, url?: string): Promise<ParseOpenAPIResponse> {
  return apiPost<ParseOpenAPIResponse>('/openapi/parse', { content, url })
}

export async function importEndpoints(request: ImportOpenAPIRequest): Promise<ImportOpenAPIResponse> {
  return apiPost<ImportOpenAPIResponse>('/openapi/import', request)
}
