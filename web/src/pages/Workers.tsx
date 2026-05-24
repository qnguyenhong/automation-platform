import { useQuery } from '@tanstack/react-query'
import { workersApi } from '../api/workers'
import { Server, Wifi, WifiOff, Activity } from 'lucide-react'
import { Link } from 'react-router-dom'

export default function Workers() {
  const { data: workers, isLoading } = useQuery({
    queryKey: ['workers'],
    queryFn: workersApi.list,
  })

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Workers</h1>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-green-500 rounded-full"></div>
            <span className="text-sm text-gray-500">
              {workers?.filter((w) => w.status === 'online').length || 0} Online
            </span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-gray-400 rounded-full"></div>
            <span className="text-sm text-gray-500">
              {workers?.filter((w) => w.status === 'offline').length || 0} Offline
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {workers?.map((worker) => (
          <div key={worker.id} className="bg-white rounded-lg shadow p-6">
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg ${
                  worker.status === 'online' ? 'bg-green-100' :
                  worker.status === 'busy' ? 'bg-yellow-100' : 'bg-gray-100'
                }`}>
                  <Server className={`w-5 h-5 ${
                    worker.status === 'online' ? 'text-green-600' :
                    worker.status === 'busy' ? 'text-yellow-600' : 'text-gray-600'
                  }`} />
                </div>
                <div>
                  <h3 className="font-semibold">{worker.name}</h3>
                  <p className="text-sm text-gray-500">{worker.hostname}</p>
                </div>
              </div>
              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                worker.status === 'online' ? 'bg-green-100 text-green-800' :
                worker.status === 'busy' ? 'bg-yellow-100 text-yellow-800' :
                'bg-gray-100 text-gray-800'
              }`}>
                {worker.status}
              </span>
            </div>

            <div className="space-y-2 mb-4">
              <div className="flex items-center gap-2 text-sm text-gray-500">
                <Wifi className="w-4 h-4" />
                <span>IP: {worker.ip_address || 'N/A'}</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-500">
                <Activity className="w-4 h-4" />
                <span>Executors: {worker.executor_types?.join(', ') || 'None'}</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-500">
                <span>Max Concurrent: {worker.max_concurrent}</span>
              </div>
            </div>

            <div className="flex flex-wrap gap-1 mb-4">
              {Object.entries(worker.labels || {}).map(([key, value]) => (
                <span key={key} className="px-2 py-0.5 text-xs bg-gray-100 rounded">
                  {key}: {String(value)}
                </span>
              ))}
            </div>

            <Link
              to={`/workers/${worker.id}`}
              className="block text-center text-sm text-blue-600 hover:text-blue-800"
            >
              View Details
            </Link>
          </div>
        ))}
      </div>
    </div>
  )
}
