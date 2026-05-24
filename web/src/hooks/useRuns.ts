import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { runsApi } from '../api/runs'

export function useRuns(projectId: string, page = 1, perPage = 20) {
  return useQuery({
    queryKey: ['runs', projectId, page, perPage],
    queryFn: () => runsApi.list(projectId, page, perPage),
    enabled: !!projectId,
  })
}

export function useAllRuns(page = 1, perPage = 20) {
  return useQuery({
    queryKey: ['runs', 'all', page, perPage],
    queryFn: () => runsApi.listAll(page, perPage),
  })
}

export function useRun(runId: string) {
  return useQuery({
    queryKey: ['runs', runId],
    queryFn: () => runsApi.get(runId),
    enabled: !!runId,
  })
}

export function useRunResults(runId: string) {
  return useQuery({
    queryKey: ['results', runId],
    queryFn: () => runsApi.results(runId),
    enabled: !!runId,
  })
}

export function useTriggerRun(projectId: string, suiteId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data?: { trigger?: string; metadata?: any }) =>
      runsApi.trigger(projectId, suiteId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['runs'] })
    },
  })
}

export function useCancelRun() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (runId: string) => runsApi.cancel(runId),
    onSuccess: (_, runId) => {
      queryClient.invalidateQueries({ queryKey: ['runs'] })
      queryClient.invalidateQueries({ queryKey: ['runs', runId] })
    },
  })
}
