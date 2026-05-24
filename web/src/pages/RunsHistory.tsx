import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useAllRuns } from '../hooks/useRuns'
import {
  PlayCircle,
  ArrowUpRight,
  Search,
  ChevronLeft,
  ChevronRight,
  Filter,
} from 'lucide-react'

export default function RunsHistory() {
  const [page, setPage] = useState(1)
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const perPage = 15

  const { data, isLoading } = useAllRuns(page, perPage)

  // Extract pagination details from JSONWithMeta structure
  // Note: JSONWithMeta returns { data: T[], meta: { page, per_page, total, total_pages } }
  // Let's typecast/safely handle the structure
  const responseData = data as any
  const runs = responseData?.data || []
  const pagination = responseData?.meta || { page: 1, per_page: 15, total: 0, total_pages: 1 }

  // Client side search and filter for quick search experience
  const filteredRuns = runs.filter((run: any) => {
    const matchesSearch =
      run.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
      run.suite_id.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesStatus = statusFilter === '' || run.status === statusFilter
    return matchesSearch && matchesStatus
  })

  const totalPages = pagination.total_pages || 1

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-800">Global Execution History</h1>
        <p className="text-xs text-slate-500 mt-1">
          Historical overview of all test executions and load profile runs across the workspace.
        </p>
      </div>

      {/* Filters Toolbar */}
      <div className="glass-panel p-4 rounded-xl flex flex-col md:flex-row items-center justify-between gap-4">
        <div className="relative w-full md:w-80">
          <span className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
            <Search className="w-4 h-4" />
          </span>
          <input
            type="text"
            placeholder="Search by Run ID or Suite ID..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2 bg-white border border-slate-200 rounded-lg text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
          />
        </div>

        <div className="flex items-center gap-3 w-full md:w-auto">
          <div className="flex items-center gap-2 text-slate-500 text-xs font-semibold shrink-0">
            <Filter className="w-3.5 h-3.5" />
            Filter Status:
          </div>
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="px-3 py-2 bg-white border border-slate-200 text-slate-700 text-xs font-semibold rounded-lg focus:outline-none focus:ring-1 focus:ring-indigo-500 transition cursor-pointer"
          >
            <option value="">All Statuses</option>
            <option value="passed">Passed</option>
            <option value="failed">Failed</option>
            <option value="running">Running</option>
            <option value="pending">Pending</option>
            <option value="cancelled">Cancelled</option>
            <option value="error">Error</option>
          </select>
        </div>
      </div>

      {/* Runs Table */}
      <div className="glass-panel rounded-xl border border-slate-200 overflow-hidden shadow-sm">
        {isLoading ? (
          <div className="flex flex-col items-center justify-center p-12 space-y-4">
            <div className="w-10 h-10 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin"></div>
            <p className="text-slate-500 text-sm font-medium animate-pulse">Loading history logs...</p>
          </div>
        ) : filteredRuns.length === 0 ? (
          <div className="text-center p-12 text-slate-500">
            <PlayCircle className="w-12 h-12 mx-auto mb-3 text-slate-300" />
            <p className="font-semibold text-slate-700 text-base">No Execution Logs Found</p>
            <p className="text-xs text-slate-500 mt-1">Try resetting the status filter or search query.</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead className="bg-slate-50 text-slate-600 text-xs font-bold uppercase border-b border-slate-200">
                  <tr>
                    <th className="px-6 py-4">Run ID</th>
                    <th className="px-6 py-4">Suite ID</th>
                    <th className="px-6 py-4">Status</th>
                    <th className="px-6 py-4">Trigger</th>
                    <th className="px-6 py-4 font-mono">Duration</th>
                    <th className="px-6 py-4">Timestamp</th>
                    <th className="px-6 py-4 text-right">Link</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100 text-sm text-slate-700">
                  {filteredRuns.map((run: any) => (
                    <tr key={run.id} className="hover:bg-slate-50/80 transition-colors">
                      <td className="px-6 py-4 font-mono font-bold text-indigo-600">
                        {run.id.slice(0, 8)}...
                      </td>
                      <td className="px-6 py-4 font-mono text-xs text-slate-500">
                        {run.suite_id.slice(0, 8)}...
                      </td>
                      <td className="px-6 py-4">
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                            run.status === 'passed'
                              ? 'bg-emerald-50 border-emerald-200 text-emerald-700'
                              : run.status === 'failed'
                              ? 'bg-rose-50 border-rose-200 text-rose-700'
                              : run.status === 'running'
                              ? 'bg-blue-50 border-blue-200 text-blue-700 animate-pulse'
                              : 'bg-slate-50 border-slate-200 text-slate-650'
                          }`}
                        >
                          {run.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 font-medium text-slate-700 capitalize">{run.trigger}</td>
                      <td className="px-6 py-4 font-mono text-slate-600">
                        {run.duration_ms ? `${(run.duration_ms / 1000).toFixed(1)}s` : '—'}
                      </td>
                      <td className="px-6 py-4 text-slate-500">
                        {new Date(run.created_at).toLocaleString()}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <Link
                          to={`/runs/${run.id}`}
                          className="inline-flex items-center gap-1 text-xs font-bold text-indigo-600 hover:text-indigo-800 transition-colors"
                        >
                          Analyze
                          <ArrowUpRight className="w-3.5 h-3.5" />
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination controls */}
            {totalPages > 1 && (
              <div className="p-4 border-t border-slate-200 bg-slate-50/50 flex items-center justify-between">
                <span className="text-xs text-slate-555 font-medium">
                  Showing Page {page} of {totalPages}
                </span>
                <div className="flex items-center gap-2">
                  <button
                    disabled={page <= 1}
                    onClick={() => setPage(page - 1)}
                    className="p-1.5 border border-slate-200 rounded-lg hover:bg-slate-100 disabled:opacity-50 transition cursor-pointer text-slate-600"
                  >
                    <ChevronLeft className="w-4 h-4" />
                  </button>
                  <button
                    disabled={page >= totalPages}
                    onClick={() => setPage(page + 1)}
                    className="p-1.5 border border-slate-200 rounded-lg hover:bg-slate-100 disabled:opacity-50 transition cursor-pointer text-slate-600"
                  >
                    <ChevronRight className="w-4 h-4" />
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
