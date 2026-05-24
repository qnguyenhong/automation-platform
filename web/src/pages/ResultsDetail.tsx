import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { runsApi } from '../api/runs'
import StatusBadge from '../components/common/StatusBadge'

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
          <h1 className="text-2xl font-bold text-slate-800">Endpoint Assertion Details</h1>
          <p className="text-slate-500 font-mono text-xs mt-0.5">{resultId}</p>
        </div>
        <StatusBadge status={result?.status || 'pending'} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Request Data */}
        <div className="glass-panel rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-700 mb-4">Request Payload</h3>
          <pre className="bg-slate-50 border border-slate-200 p-4 rounded-xl overflow-auto text-xs text-indigo-700 max-h-96 shadow-inner font-mono">
            {JSON.stringify(result?.request_data, null, 2) || 'No request data'}
          </pre>
        </div>

        {/* Response Data */}
        <div className="glass-panel rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-700 mb-4">Response Payload</h3>
          <pre className="bg-slate-50 border border-slate-200 p-4 rounded-xl overflow-auto text-xs text-emerald-700 max-h-96 shadow-inner font-mono">
            {JSON.stringify(result?.response_data, null, 2) || 'No response data'}
          </pre>
        </div>
      </div>

      {/* Assertions */}
      <div className="glass-panel rounded-xl p-6 border border-slate-200 shadow-sm">
        <h3 className="text-sm font-bold uppercase tracking-wider text-slate-700 mb-4">Assertions</h3>
        <div className="space-y-3">
          {result?.assertions && result.assertions.length > 0 ? (
            result.assertions.map((assertion: any, index: number) => (
              <div
                key={index}
                className={`p-3.5 rounded-xl border flex items-start gap-3 transition ${
                  assertion.passed
                    ? 'bg-emerald-55/40 border-emerald-200 text-emerald-800'
                    : 'bg-rose-55/40 border-rose-200 text-rose-800'
                }`}
              >
                <span className={`text-base font-black select-none ${assertion.passed ? 'text-emerald-650' : 'text-rose-650'}`}>
                  {assertion.passed ? '✓' : '✗'}
                </span>
                <div>
                  <span className="font-semibold text-sm capitalize">{assertion.type.replace(/_/g, ' ')}</span>
                  {assertion.message && (
                    <p className="text-xs text-slate-500 mt-1 leading-relaxed">{assertion.message}</p>
                  )}
                </div>
              </div>
            ))
          ) : (
            <p className="text-xs text-slate-500 italic font-medium">No assertions evaluated for this target run.</p>
          )}
        </div>
      </div>

      {/* Stdout/Stderr */}
      {result?.stdout && (
        <div className="glass-panel rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-700 mb-4">Standard Output</h3>
          <pre className="bg-slate-55 border border-slate-200 p-4 rounded-xl overflow-auto text-xs text-slate-700 max-h-64 shadow-inner font-mono">
            {result.stdout}
          </pre>
        </div>
      )}

      {result?.stderr && (
        <div className="glass-panel rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-bold uppercase tracking-wider text-rose-600 mb-4">Standard Error</h3>
          <pre className="bg-rose-50 border border-rose-205 p-4 rounded-xl overflow-auto text-xs text-rose-700 max-h-64 shadow-inner font-mono">
            {result.stderr}
          </pre>
        </div>
      )}
    </div>
  )
}
