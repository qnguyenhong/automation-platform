import { useParams } from 'react-router-dom'
import { useRun, useRunResults } from '../hooks/useRuns'
import { StatusBadge } from '../components/StatusBadge'
import { Link } from 'react-router-dom'

export default function TestRunDetail() {
  const { runId } = useParams<{ runId: string }>()
  const { data: run } = useRun(runId!)
  const { data: results } = useRunResults(runId!)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Test Run</h1>
          <p className="text-gray-500 font-mono">{runId}</p>
        </div>
        <StatusBadge status={run?.status || 'pending'} />
      </div>

      {/* Run Stats */}
      <div className="grid grid-cols-5 gap-4">
        <div className="bg-white p-4 rounded-lg shadow text-center">
          <p className="text-sm text-gray-500">Total</p>
          <p className="text-2xl font-bold">{run?.total_cases || 0}</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow text-center">
          <p className="text-sm text-gray-500">Passed</p>
          <p className="text-2xl font-bold text-green-600">{run?.passed || 0}</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow text-center">
          <p className="text-sm text-gray-500">Failed</p>
          <p className="text-2xl font-bold text-red-600">{run?.failed || 0}</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow text-center">
          <p className="text-sm text-gray-500">Skipped</p>
          <p className="text-2xl font-bold text-yellow-600">{run?.skipped || 0}</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow text-center">
          <p className="text-sm text-gray-500">Duration</p>
          <p className="text-2xl font-bold">{run?.duration_ms ? `${(run.duration_ms / 1000).toFixed(1)}s` : '-'}</p>
        </div>
      </div>

      {/* Results Table */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b border-gray-200">
          <h3 className="text-lg font-semibold">Results</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Case ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Duration</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Error</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {results?.map((result) => (
                <tr key={result.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <StatusBadge status={result.status} />
                  </td>
                  <td className="px-6 py-4 text-sm font-mono">{result.case_id.slice(0, 8)}</td>
                  <td className="px-6 py-4 text-sm text-gray-500">
                    {result.duration_ms ? `${result.duration_ms}ms` : '-'}
                  </td>
                  <td className="px-6 py-4 text-sm text-red-500 max-w-xs truncate">
                    {result.error_message || '-'}
                  </td>
                  <td className="px-6 py-4">
                    <Link
                      to={`/runs/${runId}/results/${result.id}`}
                      className="text-blue-600 hover:text-blue-800 text-sm"
                    >
                      View Details
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
