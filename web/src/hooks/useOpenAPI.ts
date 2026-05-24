import { useMutation } from '@tanstack/react-query'
import { parseSpec, importEndpoints } from '@/api/openapi'
import type { ImportOpenAPIRequest } from '@/types/openapi'

export function useParseSpec() {
  return useMutation({
    mutationFn: ({ content, url }: { content: string; url?: string }) =>
      parseSpec(content, url),
  })
}

export function useImportEndpoints() {
  return useMutation({
    mutationFn: (request: ImportOpenAPIRequest) => importEndpoints(request),
  })
}
