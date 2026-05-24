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
    return <div>Loading...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{worker.name}</h1>
          <p className="text-gray-500">{worker.hostname}</p>
        </div>
        <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
          worker.status === 'online' ? 'bg-green-100 text-green-800' :
          worker.status === 'busy' ? 'bg-yellow-100 text-yellow-800' :
          'bg-gray-100 text-gray-800'
        }`}>
          {worker.status}
        </span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Worker Info */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">Worker Information</h3>
          <dl className="space-y-3">
            <div>
              <dt className="text-sm text-gray-500">ID</dt>
              <dd className="font-mono text-sm">{worker.id}</dd>
            </div>
            <div>
              <dt className="text-sm text-gray-500">IP Address</dt>
              <dd>{worker.ip_address || 'N/A'}</dd>
            </div>
            <div>
              <dt className="text-sm text-gray-500">Version</dt>
              <dd>{worker.version || 'N/A'}</dd>
            </div>
            <div>
              <dt className="text-sm text-gray-500">Max Concurrent</dt>
              <dd>{worker.max_concurrent}</dd>
            </div>
            <div>
              <dt className="text-sm text-gray-500">Last Heartbeat</dt>
              <dd>{worker.last_heartbeat ? new Date(worker.last_heartbeat).toLocaleString() : 'Never'}</dd>
            </div>
            <div>
              <dt className="text-sm text-gray-500">Registered At</dt>
              <dd>{new Date(worker.registered_at).toLocaleString()}</dd>
            </div>
          </dl>
        </div>

        {/* Executors */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">Supported Executors</h3>
          <div className="space-y-2">
            {worker.executor_types?.map((type) => (
              <div key={type} className="flex items-center gap-2 p-2 bg-gray-50 rounded-lg">
                <Activity className="w-4 h-4 text-blue-600" />
                <span className="font-medium">{type}</span>
              </div>
            ))}
          </div>

          <h3 className="text-lg font-semibold mt-6 mb-4">Labels</h3>
          <div className="flex flex-wrap gap-2">
            {Object.entries(worker.labels || {}).map(([key, value]) => (
              <span key={key} className="px-3 py-1 text-sm bg-gray-100 rounded">
                {key}: {String(value)}
              </span>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
