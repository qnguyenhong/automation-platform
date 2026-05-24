import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { workersApi } from '../api/workers'
import { Activity } from 'lucide-react'

export default function WorkerDetail() {
  const { workerId } = useParams<{ workerId: string }>()

  const { data: worker } = useQuery({
    queryKey: ['workers', workerId],
    queryFn: () => workersApi.get(workerId!),
    enabled: !!workerId,
  })

  if (!worker) {
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
          <h1 className="text-2xl font-bold text-slate-800">{worker.name}</h1>
          <p className="text-slate-500 font-mono text-xs mt-0.5">{worker.hostname}</p>
        </div>
        <span className={`inline-flex items-center px-3 py-1 rounded-md text-xs font-bold uppercase border ${
          worker.status === 'online' ? 'bg-emerald-50 border border-emerald-250 text-emerald-700' :
          worker.status === 'busy' ? 'bg-amber-50 border border-amber-255 text-amber-700' :
          'bg-slate-50 border border-slate-200 text-slate-600'
        }`}>
          {worker.status}
        </span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Worker Info */}
        <div className="glass-panel rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-700 mb-4">Worker Information</h3>
          <dl className="space-y-4">
            <div>
              <dt className="text-xs text-slate-550 font-bold uppercase tracking-wider">Node ID</dt>
              <dd className="font-mono text-sm text-indigo-650 font-semibold mt-1">{worker.id}</dd>
            </div>
            <div>
              <dt className="text-xs text-slate-550 font-bold uppercase tracking-wider">IP Address</dt>
              <dd className="text-sm font-mono mt-1 text-slate-800">{worker.ip_address || '—'}</dd>
            </div>
            <div>
              <dt className="text-xs text-slate-550 font-bold uppercase tracking-wider">Agent Version</dt>
              <dd className="text-sm font-mono mt-1 text-slate-800">{worker.version || '—'}</dd>
            </div>
            <div>
              <dt className="text-xs text-slate-550 font-bold uppercase tracking-wider">Max Concurrency Limit</dt>
              <dd className="text-sm mt-1 text-slate-800 font-semibold">{worker.max_concurrent} active concurrent threads</dd>
            </div>
            <div>
              <dt className="text-xs text-slate-550 font-bold uppercase tracking-wider">Last Heartbeat Ping</dt>
              <dd className="text-sm mt-1 text-slate-700">{worker.last_heartbeat ? new Date(worker.last_heartbeat).toLocaleString() : 'Never'}</dd>
            </div>
            <div>
              <dt className="text-xs text-slate-550 font-bold uppercase tracking-wider">Node Registration Time</dt>
              <dd className="text-sm mt-1 text-slate-700">{new Date(worker.registered_at).toLocaleString()}</dd>
            </div>
          </dl>
        </div>

        {/* Executors */}
        <div className="glass-panel rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-700 mb-4">Supported Execution Engines</h3>
          <div className="space-y-2">
            {worker.executor_types?.map((type) => (
              <div key={type} className="flex items-center gap-2.5 p-3 bg-slate-50 border border-slate-200 rounded-xl shadow-sm">
                <Activity className="w-4 h-4 text-indigo-600" />
                <span className="font-semibold text-sm text-slate-800 capitalize">{type}</span>
              </div>
            ))}
          </div>

          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-700 mt-6 mb-4">Node Labels</h3>
          <div className="flex flex-wrap gap-2">
            {worker.labels && Object.keys(worker.labels).length > 0 ? (
              Object.entries(worker.labels).map(([key, value]) => (
                <span key={key} className="px-3 py-1 text-xs font-mono bg-slate-55 border border-slate-200 text-slate-700 rounded-md font-semibold">
                  {key}: {String(value)}
                </span>
              ))
            ) : (
              <p className="text-xs text-slate-500 italic font-medium">No custom labels configured on this node.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
