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
          <h1 className="text-2xl font-bold text-white">Endpoint Assertion Details</h1>
          <p className="text-slate-500 font-mono text-xs mt-0.5">{resultId}</p>
        </div>
        <StatusBadge status={result?.status || 'pending'} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Request Data */}
        <div className="glass-panel rounded-xl p-6 border border-slate-850">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-400 mb-4">Request Payload</h3>
          <pre className="bg-slate-950/80 border border-slate-850 p-4 rounded-xl overflow-auto text-xs text-indigo-300 max-h-96">
            {JSON.stringify(result?.request_data, null, 2) || 'No request data'}
          </pre>
        </div>

        {/* Response Data */}
        <div className="glass-panel rounded-xl p-6 border border-slate-850">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-400 mb-4">Response Payload</h3>
          <pre className="bg-slate-950/80 border border-slate-850 p-4 rounded-xl overflow-auto text-xs text-emerald-350 max-h-96">
            {JSON.stringify(result?.response_data, null, 2) || 'No response data'}
          </pre>
        </div>
      </div>

      {/* Assertions */}
      <div className="glass-panel rounded-xl p-6 border border-slate-850">
        <h3 className="text-sm font-bold uppercase tracking-wider text-slate-400 mb-4">Assertions</h3>
        <div className="space-y-3">
          {result?.assertions && result.assertions.length > 0 ? (
            result.assertions.map((assertion: any, index: number) => (
              <div
                key={index}
                className={`p-3.5 rounded-xl border flex items-start gap-3 transition ${
                  assertion.passed
                    ? 'bg-emerald-500/5 border-emerald-500/20 text-emerald-400'
                    : 'bg-rose-500/5 border-rose-500/20 text-rose-400'
                }`}
              >
                <span className={`text-base font-bold select-none ${assertion.passed ? 'text-emerald-450' : 'text-rose-455'}`}>
                  {assertion.passed ? '✓' : '✗'}
                </span>
                <div>
                  <span className="font-semibold text-sm capitalize">{assertion.type.replace(/_/g, ' ')}</span>
                  {assertion.message && (
                    <p className="text-xs text-slate-400 mt-1 leading-relaxed">{assertion.message}</p>
                  )}
                </div>
              </div>
            ))
          ) : (
            <p className="text-xs text-slate-500 italic">No assertions evaluated for this target run.</p>
          )}
        </div>
      </div>

      {/* Stdout/Stderr */}
      {result?.stdout && (
        <div className="glass-panel rounded-xl p-6 border border-slate-850">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-400 mb-4">Standard Output</h3>
          <pre className="bg-slate-950/80 border border-slate-850 p-4 rounded-xl overflow-auto text-xs text-slate-300 max-h-64">
            {result.stdout}
          </pre>
        </div>
      )}

      {result?.stderr && (
        <div className="glass-panel rounded-xl p-6 border border-slate-850">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-400 mb-4 text-rose-400">Standard Error</h3>
          <pre className="bg-rose-950/20 border border-rose-500/20 p-4 rounded-xl overflow-auto text-xs text-rose-400 max-h-64">
            {result.stderr}
          </pre>
        </div>
      )}
    </div>
  )
}
