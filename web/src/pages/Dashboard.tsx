import { useQuery } from '@tanstack/react-query'
import { dashboardApi } from '../api/dashboard'
import { runsApi } from '../api/runs'
import { useProjects } from '../hooks/useProjects'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts'
import {
  Activity,
  CheckCircle,
  XCircle,
  Server,
  ArrowUpRight,
} from 'lucide-react'

const COLORS = ['#10B981', '#EF4444', '#F59E0B', '#64748B']

export default function Dashboard() {
  const { data: projects } = useProjects()
  const [selectedProject, setSelectedProject] = useState<string>('')

  const { data: summary } = useQuery({
    queryKey: ['dashboard', 'summary', selectedProject],
    queryFn: () => dashboardApi.summary(selectedProject || undefined),
  })

  const { data: trends } = useQuery({
    queryKey: ['dashboard', 'trends', selectedProject],
    queryFn: () => dashboardApi.trends(selectedProject || undefined),
  })

  const { data: recentRuns } = useQuery({
    queryKey: ['runs', selectedProject],
    queryFn: () => runsApi.list(selectedProject, 1, 10),
  })

  const responseData = recentRuns as any
  const runsList = responseData?.data || []

  const pieData = summary
    ? [
        { name: 'Successful', value: summary.total_passed },
        { name: 'Failed', value: summary.total_failed },
        { name: 'Other', value: Math.max(0, summary.total_runs - summary.total_passed - summary.total_failed) },
      ]
    : []

  return (
    <div className="space-y-6">
      {/* Header and Project Select */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-800">Analytics Dashboard</h1>
          <p className="text-xs text-slate-500 mt-1">Global performance metrics and automated verification summary.</p>
        </div>

        <select
          value={selectedProject}
          onChange={(e) => setSelectedProject(e.target.value)}
          className="px-3.5 py-2 bg-white border border-slate-200 hover:border-slate-350 text-slate-700 text-xs font-semibold rounded-lg focus:outline-none focus:ring-1 focus:ring-indigo-500 transition cursor-pointer shadow-sm"
        >
          <option value="">All Workspaces</option>
          {projects?.map((project) => (
            <option key={project.id} value={project.id}>
              {project.name}
            </option>
          ))}
        </select>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex items-center gap-4 shadow-sm">
          <div className="p-3 bg-blue-50 text-blue-600 border border-blue-200 rounded-xl">
            <Activity className="w-5 h-5" />
          </div>
          <div>
            <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Total Executions</p>
            <p className="text-2xl font-bold text-slate-800 font-mono mt-0.5">{summary?.total_runs || 0}</p>
          </div>
        </div>

        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex items-center gap-4 shadow-sm">
          <div className="p-3 bg-emerald-50 text-emerald-600 border border-emerald-200 rounded-xl">
            <CheckCircle className="w-5 h-5" />
          </div>
          <div>
            <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Successful Runs</p>
            <p className="text-2xl font-bold text-emerald-600 font-mono mt-0.5">{summary?.total_passed || 0}</p>
          </div>
        </div>

        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex items-center gap-4 shadow-sm">
          <div className="p-3 bg-rose-50 text-rose-600 border border-rose-200 rounded-xl">
            <XCircle className="w-5 h-5" />
          </div>
          <div>
            <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Failed Runs</p>
            <p className="text-2xl font-bold text-rose-650 font-mono mt-0.5">{summary?.total_failed || 0}</p>
          </div>
        </div>

        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex items-center gap-4 shadow-sm">
          <div className="p-3 bg-purple-50 text-purple-600 border border-purple-200 rounded-xl">
            <Server className="w-5 h-5" />
          </div>
          <div>
            <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Active Runner Nodes</p>
            <p className="text-2xl font-bold text-purple-650 font-mono mt-0.5">{summary?.active_workers || 0}</p>
          </div>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="glass-panel p-6 rounded-xl border border-slate-200/60 shadow-sm">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-650 mb-4">Execution Status Trends</h3>
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={trends || []}>
              <CartesianGrid strokeDasharray="3 3" stroke="#E2E8F0" vertical={false} />
              <XAxis dataKey="date" stroke="#64748B" fontSize={11} tickLine={false} />
              <YAxis stroke="#64748B" fontSize={11} tickLine={false} />
              <Tooltip
                contentStyle={{
                  backgroundColor: '#FFFFFF',
                  borderColor: '#E2E8F0',
                  borderRadius: '10px',
                  color: '#1E293B',
                  fontSize: '12px',
                }}
              />
              <Bar dataKey="passed" name="Successful" fill="#10B981" radius={[4, 4, 0, 0]} />
              <Bar dataKey="failed" name="Failed" fill="#EF4444" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="glass-panel p-6 rounded-xl border border-slate-200/60 shadow-sm">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-650 mb-4">Global Distribution</h3>
          <div className="flex items-center justify-center">
            <ResponsiveContainer width="100%" height={260}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {pieData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#FFFFFF',
                    borderColor: '#E2E8F0',
                    borderRadius: '10px',
                    color: '#1E293B',
                    fontSize: '12px',
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
            <div className="flex flex-col gap-2 shrink-0 pr-8 text-xs font-semibold text-slate-500">
              <div className="flex items-center gap-2">
                <span className="w-2.5 h-2.5 bg-emerald-500 rounded-full" />
                <span>Successful ({summary?.total_passed || 0})</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-2.5 h-2.5 bg-rose-500 rounded-full" />
                <span>Failed ({summary?.total_failed || 0})</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-2.5 h-2.5 bg-slate-450 rounded-full" />
                <span>Other</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Runs */}
      <div className="glass-panel rounded-xl border border-slate-200/60 overflow-hidden shadow-sm">
        <div className="p-5 border-b border-slate-200 bg-slate-50/50 flex items-center justify-between">
          <h3 className="text-sm font-bold uppercase tracking-wider text-slate-650">Recent Execution History</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead className="bg-slate-50 text-slate-600 text-xs font-bold uppercase border-b border-slate-200">
              <tr>
                <th className="px-6 py-4">Execution Run ID</th>
                <th className="px-6 py-4">Execution Status</th>
                <th className="px-6 py-4">Initiator</th>
                <th className="px-6 py-4 font-mono">Runtime</th>
                <th className="px-6 py-4">Launched At</th>
                <th className="px-6 py-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 text-sm text-slate-700">
              {runsList?.map((run: any) => (
                <tr key={run.id} className="hover:bg-slate-50/80 transition-colors">
                  <td className="px-6 py-4 font-mono font-bold text-indigo-600">{run.id.slice(0, 8)}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                      run.status === 'passed' ? 'bg-emerald-50 border-emerald-200 text-emerald-700' :
                      run.status === 'failed' ? 'bg-rose-50 border-rose-200 text-rose-700' :
                      run.status === 'running' ? 'bg-blue-50 border-blue-200 text-blue-700 animate-pulse' :
                      'bg-slate-50 border-slate-200 text-slate-600'
                    }`}>
                      {run.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-700 font-medium">{run.trigger}</td>
                  <td className="px-6 py-4 text-slate-600 font-mono">
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
      </div>
    </div>
  )
}
