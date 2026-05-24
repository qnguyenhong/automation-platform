import { useQuery } from '@tanstack/react-query'
import { workersApi } from '../api/workers'
import { Server, Wifi, Activity } from 'lucide-react'
import { Link } from 'react-router-dom'

export default function Workers() {
  const { data: workers, isLoading } = useQuery({
    queryKey: ['workers'],
    queryFn: workersApi.list,
  })

  if (isLoading) {
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
          <h1 className="text-2xl font-bold text-slate-800 font-sans">Worker Nodes</h1>
          <p className="text-slate-500 text-xs mt-1">Status and allocation profile of distributed execution agents.</p>
        </div>
        <div className="flex items-center gap-4 text-xs font-semibold font-mono bg-white p-2.5 rounded-lg border border-slate-205 shadow-sm">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
            <span className="text-slate-600">
              ONLINE: {workers?.filter((w) => w.status === 'online').length || 0}
            </span>
          </div>
          <div className="w-[1px] h-3 bg-slate-200" />
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 bg-slate-400 rounded-full"></div>
            <span className="text-slate-500">
              OFFLINE: {workers?.filter((w) => w.status === 'offline').length || 0}
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {workers?.map((worker) => (
          <div key={worker.id} className="glass-panel glass-panel-hover p-6 rounded-xl border border-slate-200 flex flex-col justify-between shadow-sm">
            <div>
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className={`p-2.5 rounded-lg border ${
                    worker.status === 'online' ? 'bg-emerald-50 border-emerald-200 text-emerald-700' :
                    worker.status === 'busy' ? 'bg-amber-50 border-amber-200 text-amber-700' :
                    'bg-slate-50 border-slate-200 text-slate-500'
                  }`}>
                    <Server className="w-5 h-5" />
                  </div>
                  <div>
                    <h3 className="font-bold text-slate-800">{worker.name}</h3>
                    <p className="text-xs text-slate-500 font-mono mt-0.5">{worker.hostname}</p>
                  </div>
                </div>
                <span className={`inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold uppercase border ${
                  worker.status === 'online' ? 'bg-emerald-50 border border-emerald-250 text-emerald-700' :
                  worker.status === 'busy' ? 'bg-amber-50 border border-amber-250 text-amber-700' :
                  'bg-slate-50 border border-slate-200 text-slate-600'
                }`}>
                  {worker.status}
                </span>
              </div>

              <div className="space-y-2 mb-5 pt-2 border-t border-slate-200">
                <div className="flex items-center gap-2 text-xs text-slate-600">
                  <Wifi className="w-3.5 h-3.5 text-slate-450 shrink-0" />
                  <span className="font-mono">IP: {worker.ip_address || '—'}</span>
                </div>
                <div className="flex items-center gap-2 text-xs text-slate-600">
                  <Activity className="w-3.5 h-3.5 text-slate-450 shrink-0" />
                  <span>Engine Nodes: <strong className="font-semibold text-slate-800 capitalize">{worker.executor_types?.join(', ') || 'None'}</strong></span>
                </div>
                <div className="flex items-center gap-2 text-xs text-slate-600">
                  <span className="w-3.5 h-3.5 border border-slate-200 rounded flex items-center justify-center text-[8px] font-bold text-slate-450 font-mono shrink-0">C</span>
                  <span>Concurrent Limit: <strong className="font-semibold text-slate-800">{worker.max_concurrent} Threads</strong></span>
                </div>
              </div>

              {worker.labels && Object.keys(worker.labels).length > 0 && (
                <div className="flex flex-wrap gap-1.5 mb-6">
                  {Object.entries(worker.labels).map(([key, value]) => (
                    <span key={key} className="px-2 py-0.5 text-[10px] font-mono bg-slate-50 border border-slate-200 text-slate-600 rounded font-semibold">
                      {key}: {String(value)}
                    </span>
                  ))}
                </div>
              )}
            </div>

            <Link
              to={`/workers/${worker.id}`}
              className="block text-center text-xs font-bold text-indigo-600 hover:text-indigo-850 py-2 border border-slate-200 hover:border-slate-350 bg-white rounded-lg transition shadow-sm"
            >
              Analyze Worker Node
            </Link>
          </div>
        ))}
      </div>
    </div>
  )
}
