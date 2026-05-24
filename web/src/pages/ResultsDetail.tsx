import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { runsApi } from '../api/runs'
import { StatusBadge } from '../components/StatusBadge'

export default function ResultsDetail() {
  const { runId, resultId } = useParams<{ runId: string; resultId: string }>()

  const { data: result } = useQuery({
    queryKey: ['result', runId, resultId],
    queryFn: () => runsApi.result(runId!, resultId!),
    enabled: !!runId && !!resultId,
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Result Detail</h1>
          <p className="text-gray-500 font-mono">{resultId}</p>
        </div>
        <StatusBadge status={result?.status || 'pending'} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Request Data */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">Request</h3>
          <pre className="bg-gray-50 p-4 rounded-lg overflow-auto text-sm">
            {JSON.stringify(result?.request_data, null, 2) || 'No request data'}
          </pre>
        </div>

        {/* Response Data */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">Response</h3>
          <pre className="bg-gray-50 p-4 rounded-lg overflow-auto text-sm">
            {JSON.stringify(result?.response_data, null, 2) || 'No response data'}
          </pre>
        </div>
      </div>

      {/* Assertions */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold mb-4">Assertions</h3>
        <div className="space-y-2">
          {result?.assertions?.map((assertion: any, index: number) => (
            <div
              key={index}
              className={`p-3 rounded-lg ${
                assertion.passed ? 'bg-green-50' : 'bg-red-50'
              }`}
            >
              <div className="flex items-center gap-2">
                <span className={assertion.passed ? 'text-green-600' : 'text-red-600'}>
                  {assertion.passed ? '✓' : '✗'}
                </span>
                <span className="font-medium">{assertion.type}</span>
              </div>
              {assertion.message && (
                <p className="text-sm text-gray-600 mt-1">{assertion.message}</p>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Stdout/Stderr */}
      {result?.stdout && (
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">Standard Output</h3>
          <pre className="bg-gray-50 p-4 rounded-lg overflow-auto text-sm max-h-64">
            {result.stdout}
          </pre>
        </div>
      )}

      {result?.stderr && (
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">Standard Error</h3>
          <pre className="bg-red-50 p-4 rounded-lg overflow-auto text-sm max-h-64">
            {result.stderr}
          </pre>
        </div>
      )}
    </div>
  )
}
