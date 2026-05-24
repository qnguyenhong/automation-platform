import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { suitesApi, casesApi, TestSuite, TestCase } from '../api/suites'

export function useSuites(projectId: string) {
  return useQuery({
    queryKey: ['suites', projectId],
    queryFn: () => suitesApi.list(projectId),
    enabled: !!projectId,
  })
}

export function useSuite(projectId: string, suiteId: string) {
  return useQuery({
    queryKey: ['suites', projectId, suiteId],
    queryFn: () => suitesApi.get(projectId, suiteId),
    enabled: !!projectId && !!suiteId,
  })
}

export function useCreateSuite(projectId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<TestSuite>) => suitesApi.create(projectId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['suites', projectId] })
    },
  })
}

export function useCases(projectId: string, suiteId: string) {
  return useQuery({
    queryKey: ['cases', projectId, suiteId],
    queryFn: () => casesApi.list(projectId, suiteId),
    enabled: !!projectId && !!suiteId,
  })
}

export function useCreateCase(projectId: string, suiteId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<TestCase>) => casesApi.create(projectId, suiteId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cases', projectId, suiteId] })
    },
  })
}
