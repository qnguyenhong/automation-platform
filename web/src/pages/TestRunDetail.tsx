import { useParams, Link } from 'react-router-dom'
import { useRun, useRunResults } from '../hooks/useRuns'
import StatusBadge from '../components/common/StatusBadge'
import LoadMetricsDashboard from '../components/load/LoadMetricsDashboard'

export default function TestRunDetail() {
  const { runId } = useParams<{ runId: string }>()
  const { data: run, isLoading: isRunLoading } = useRun(runId!)
  const { data: results } = useRunResults(runId!)

  const isRunning = run?.status === 'running' || run?.status === 'pending'
  // Determine if this run is a load test
  const isLoadTest = run?.total_cases === 1 && results && results.length > 0 && results[0].metrics !== null

  if (isRunLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="w-8 h-8 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-slate-800">Execution Run</h1>
          <p className="text-slate-500 font-mono text-xs mt-0.5">{runId}</p>
        </div>
        <StatusBadge status={run?.status || 'pending'} />
      </div>

      {/* Run Stats */}
      <div className="grid grid-cols-5 gap-4">
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-200 shadow-sm">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Total Targets</p>
          <p className="text-2xl font-bold text-slate-800 mt-1 font-mono">{run?.total_cases || 0}</p>
        </div>
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-200 shadow-sm">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Successful</p>
          <p className="text-2xl font-bold text-emerald-605 mt-1 font-mono">{run?.passed || 0}</p>
        </div>
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-200 shadow-sm">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Failed</p>
          <p className="text-2xl font-bold text-rose-600 mt-1 font-mono">{run?.failed || 0}</p>
        </div>
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-200 shadow-sm">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Excluded</p>
          <p className="text-2xl font-bold text-amber-600 mt-1 font-mono">{run?.skipped || 0}</p>
        </div>
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-200 shadow-sm">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Runtime</p>
          <p className="text-2xl font-bold text-indigo-600 mt-1 font-mono">{run?.duration_ms ? `${(run.duration_ms / 1000).toFixed(1)}s` : '—'}</p>
        </div>
      </div>

      {/* Load metrics dashboard */}
      {isLoadTest && (
        <div className="glass-panel p-6 rounded-xl border border-slate-200 shadow-sm">
          <LoadMetricsDashboard runId={runId!} isRunning={isRunning} />
        </div>
      )}

      {/* Results Table */}
      <div className="glass-panel rounded-xl border border-slate-200 overflow-hidden shadow-sm">
        <div className="p-5 border-b border-slate-200 bg-slate-50/50">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-700">Step Assertions / Metrics</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead className="bg-slate-50 text-slate-600 text-xs font-bold uppercase border-b border-slate-200">
              <tr>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4">Target ID</th>
                <th className="px-6 py-4">Runtime</th>
                <th className="px-6 py-4">Error / Issue Details</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 text-sm text-slate-700">
              {results?.map((result) => (
                <tr key={result.id} className="hover:bg-slate-50/80 transition-colors">
                  <td className="px-6 py-4">
                    <StatusBadge status={result.status} />
                  </td>
                  <td className="px-6 py-4 text-slate-700 font-mono font-semibold">{result.case_id.slice(0, 8)}</td>
                  <td className="px-6 py-4 text-slate-500 font-mono">
                    {result.duration_ms ? `${result.duration_ms}ms` : '—'}
                  </td>
                  <td className="px-6 py-4 text-rose-600 max-w-xs truncate font-medium">
                    {result.error_message || '—'}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <Link
                      to={`/runs/${runId}/results/${result.id}`}
                      className="inline-flex items-center text-xs font-bold text-indigo-600 hover:text-indigo-800 transition-colors"
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
