import { useParams, Link } from 'react-router-dom'
import { useRun, useRunResults } from '../hooks/useRuns'
import { useSuite } from '../hooks/useSuites'
import { StatusBadge } from '../components/StatusBadge'
import LoadMetricsDashboard from '../components/load/LoadMetricsDashboard'

export default function TestRunDetail() {
  const { runId } = useParams<{ runId: string }>()
  const { data: run } = useRun(runId!)
  const { data: results } = useRunResults(runId!)
  const { data: suite } = useSuite(run?.project_id || '', run?.suite_id || '')

  const isLoadTest = suite?.test_type === 'load'
  const isRunning = run?.status === 'running' || run?.status === 'pending'

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Execution Run</h1>
          <p className="text-slate-550 font-mono text-xs mt-0.5">{runId}</p>
        </div>
        <StatusBadge status={run?.status || 'pending'} />
      </div>

      {/* Run Stats */}
      <div className="grid grid-cols-5 gap-4">
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-850">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Total Targets</p>
          <p className="text-2xl font-bold text-slate-200 mt-1 font-mono">{run?.total_cases || 0}</p>
        </div>
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-850">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Successful</p>
          <p className="text-2xl font-bold text-emerald-400 mt-1 font-mono">{run?.passed || 0}</p>
        </div>
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-850">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Failed</p>
          <p className="text-2xl font-bold text-rose-500 mt-1 font-mono">{run?.failed || 0}</p>
        </div>
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-850">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Excluded</p>
          <p className="text-2xl font-bold text-amber-500 mt-1 font-mono">{run?.skipped || 0}</p>
        </div>
        <div className="glass-panel p-4 rounded-xl text-center border border-slate-850">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Runtime</p>
          <p className="text-2xl font-bold text-indigo-400 mt-1 font-mono">{run?.duration_ms ? `${(run.duration_ms / 1000).toFixed(1)}s` : '—'}</p>
        </div>
      </div>

      {/* Load metrics dashboard */}
      {isLoadTest && (
        <div className="glass-panel p-6 rounded-xl border border-slate-850">
          <LoadMetricsDashboard runId={runId!} isRunning={isRunning} />
        </div>
      )}

      {/* Results Table */}
      <div className="glass-panel rounded-xl border border-slate-850 overflow-hidden">
        <div className="p-5 border-b border-slate-850 bg-slate-900/30">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-400">Step Assertions / Metrics</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead className="bg-[#0B0F19]/50 text-slate-450 text-xs font-bold uppercase border-b border-slate-850">
              <tr>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4">Target ID</th>
                <th className="px-6 py-4">Runtime</th>
                <th className="px-6 py-4">Error / Issue Details</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-850 text-sm">
              {results?.map((result) => (
                <tr key={result.id} className="hover:bg-slate-800/10 transition-colors">
                  <td className="px-6 py-4">
                    <StatusBadge status={result.status} />
                  </td>
                  <td className="px-6 py-4 text-slate-350 font-mono font-semibold">{result.case_id.slice(0, 8)}</td>
                  <td className="px-6 py-4 text-slate-400 font-mono">
                    {result.duration_ms ? `${result.duration_ms}ms` : '—'}
                  </td>
                  <td className="px-6 py-4 text-rose-455 max-w-xs truncate font-medium">
                    {result.error_message || '—'}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <Link
                      to={`/runs/${runId}/results/${result.id}`}
                      className="inline-flex items-center text-xs font-bold text-indigo-400 hover:text-indigo-300 transition-colors"
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
