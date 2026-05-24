import { useEffect, useState, useMemo } from 'react'
import { loadMetricsApi, LoadMetrics } from '../../api/loadMetrics'
import { useWebSocket } from '../../hooks/useWebSocket'
import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  BarChart, Bar, Cell
} from 'recharts'
import { Activity, Clock, AlertTriangle, Database, ArrowUpRight } from 'lucide-react'

interface DashboardProps {
  runId: string
  isRunning: boolean
}

interface LivePoint {
  second: number
  rps: number
  p95: number
}

export default function LoadMetricsDashboard({ runId, isRunning }: DashboardProps) {
  const [metrics, setMetrics] = useState<LoadMetrics | null>(null)
  const [loading, setLoading] = useState(!isRunning)
  const [liveData, setLiveData] = useState<LivePoint[]>([])

  // Fetch final metrics when test completes
  const fetchFinalMetrics = async () => {
    try {
      setLoading(true)
      const data = await loadMetricsApi.getByRun(runId)
      if (data) {
        setMetrics(data)
      }
    } catch (err) {
      console.error('Failed to fetch final load metrics:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (!isRunning) {
      fetchFinalMetrics()
    }
  }, [runId, isRunning])

  // Listen for real-time WebSocket metrics stream
  useWebSocket((data) => {
    if (isRunning && data.type === 'load_metric') {
      try {
        const payload = typeof data.metrics === 'string' ? JSON.parse(data.metrics) : data.metrics
        setLiveData((prev) => {
          // Avoid duplicate seconds
          if (prev.some((p) => p.second === payload.second)) {
            return prev
          }
          return [...prev, {
            second: payload.second,
            rps: Number(payload.rps || 0),
            p95: Number(payload.p95 || 0)
          }].sort((a, b) => a.second - b.second)
        })
      } catch (err) {
        console.error('Error parsing live socket metrics:', err)
      }
    }
  })

  // Format bytes helper
  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  // Derive charts data depending on whether we are live or viewing historic results
  const throughputData = useMemo(() => {
    if (isRunning) {
      return liveData.map((pt) => ({
        second: `${pt.second}s`,
        RPS: pt.rps
      }))
    }

    if (!metrics || !metrics.time_series) return []
    const rpsPoints = metrics.time_series.rps || []
    return rpsPoints.map((pt) => ({
      second: `${pt.second}s`,
      RPS: pt.value
    }))
  }, [isRunning, liveData, metrics])

  const latencyData = useMemo(() => {
    if (isRunning) {
      return liveData.map((pt) => ({
        second: `${pt.second}s`,
        P95: pt.p95
      }))
    }

    if (!metrics || !metrics.time_series) return []
    const p95Points = metrics.time_series.p95 || []
    return p95Points.map((pt) => ({
      second: `${pt.second}s`,
      P95: pt.value
    }))
  }, [isRunning, liveData, metrics])

  // Aggregate stats dynamically for live dashboard
  const stats = useMemo(() => {
    if (isRunning) {
      const totalReq = liveData.reduce((sum, p) => sum + p.rps, 0)
      const avgRps = liveData.length > 0 ? (totalReq / liveData.length) : 0
      const maxP95 = liveData.reduce((max, p) => p.p95 > max ? p.p95 : max, 0)
      return {
        totalRequests: totalReq,
        avgRps: parseFloat(avgRps.toFixed(1)),
        errorRate: '0.0%',
        p95: maxP95,
        avgLatency: '—',
        totalBytes: '—'
      }
    }

    if (!metrics) return null

    return {
      totalRequests: metrics.total_requests,
      avgRps: parseFloat(metrics.throughput_rps.toFixed(1)),
      errorRate: `${(metrics.error_rate * 100).toFixed(2)}%`,
      p95: metrics.p95_latency_ms,
      avgLatency: `${metrics.avg_latency_ms.toFixed(1)} ms`,
      totalBytes: formatBytes(metrics.total_bytes)
    }
  }, [isRunning, liveData, metrics])

  // Latency profile chart
  const latencyProfileData = useMemo(() => {
    if (!metrics) return []
    return [
      { name: 'Min', value: metrics.min_latency_ms },
      { name: 'P50', value: metrics.p50_latency_ms },
      { name: 'P90', value: metrics.p90_latency_ms },
      { name: 'P95', value: metrics.p95_latency_ms },
      { name: 'P99', value: metrics.p99_latency_ms },
      { name: 'Max', value: metrics.max_latency_ms }
    ]
  }, [metrics])

  // Status code distribution pie data
  const statusCodeData = useMemo(() => {
    if (!metrics || !metrics.status_codes) return []
    return Object.entries(metrics.status_codes).map(([code, count]) => ({
      name: `HTTP ${code}`,
      value: count
    }))
  }, [metrics])

  const COLORS = ['#6366F1', '#10B981', '#F59E0B', '#EF4444', '#A855F7', '#64748B']

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-64 space-y-4">
        <div className="w-12 h-12 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin"></div>
        <p className="text-slate-500 font-medium animate-pulse">Retrieving execution metrics...</p>
      </div>
    )
  }

  if (!stats) {
    return (
      <div className="glass-panel rounded-xl border border-slate-200 p-8 text-center text-slate-500 shadow-sm">
        <Activity className="w-12 h-12 mx-auto mb-3 text-slate-400 animate-pulse" />
        <p className="font-semibold text-slate-700">Awaiting execution metrics stream...</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Title / Realtime Badge */}
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-bold uppercase tracking-wider text-slate-550">Load Execution Analytics</h3>
        {isRunning && (
          <span className="flex items-center gap-1.5 px-3 py-1 text-xs font-semibold bg-emerald-50 border border-emerald-200 text-emerald-700 rounded-full">
            <span className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></span>
            Live Streaming Metrics
          </span>
        )}
      </div>

      {/* Summary KPI Cards Grid */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
        {/* Total Requests */}
        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex flex-col justify-between shadow-sm">
          <div className="flex items-center justify-between text-slate-500">
            <span className="text-[10px] font-bold uppercase tracking-wider">Total Requests</span>
            <Database className="w-4 h-4 text-indigo-500" />
          </div>
          <div className="mt-2.5">
            <h4 className="text-2xl font-bold text-slate-800 font-mono">{stats.totalRequests.toLocaleString()}</h4>
            <span className="text-[10px] text-slate-500 font-medium">requests sent</span>
          </div>
        </div>

        {/* Avg Throughput */}
        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex flex-col justify-between shadow-sm">
          <div className="flex items-center justify-between text-slate-500">
            <span className="text-[10px] font-bold uppercase tracking-wider">Avg Throughput</span>
            <Activity className="w-4 h-4 text-emerald-500" />
          </div>
          <div className="mt-2.5">
            <h4 className="text-2xl font-bold text-emerald-600 font-mono">{stats.avgRps}</h4>
            <span className="text-[10px] text-slate-500 font-medium">req / sec (RPS)</span>
          </div>
        </div>

        {/* P95 Latency */}
        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex flex-col justify-between shadow-sm">
          <div className="flex items-center justify-between text-slate-500">
            <span className="text-[10px] font-bold uppercase tracking-wider">P95 Latency</span>
            <Clock className="w-4 h-4 text-amber-500" />
          </div>
          <div className="mt-2.5">
            <h4 className="text-2xl font-bold text-amber-600 font-mono">{stats.p95} <span className="text-xs font-normal text-slate-500">ms</span></h4>
            <span className="text-[10px] text-slate-500 font-medium">95% response limit</span>
          </div>
        </div>

        {/* Avg Latency */}
        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex flex-col justify-between shadow-sm">
          <div className="flex items-center justify-between text-slate-500">
            <span className="text-[10px] font-bold uppercase tracking-wider">Avg Latency</span>
            <Clock className="w-4 h-4 text-indigo-500" />
          </div>
          <div className="mt-2.5">
            <h4 className="text-2xl font-bold text-slate-800 font-mono">{stats.avgLatency}</h4>
            <span className="text-[10px] text-slate-500 font-medium">average response</span>
          </div>
        </div>

        {/* Error Rate */}
        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex flex-col justify-between shadow-sm">
          <div className="flex items-center justify-between text-slate-500">
            <span className="text-[10px] font-bold uppercase tracking-wider">Error Rate</span>
            <AlertTriangle className="w-4 h-4 text-rose-500" />
          </div>
          <div className="mt-2.5">
            <h4 className={`text-2xl font-bold font-mono ${stats.errorRate !== '0.00%' && stats.errorRate !== '0.0%' ? 'text-rose-600' : 'text-slate-800'}`}>{stats.errorRate}</h4>
            <span className="text-[10px] text-slate-500 font-medium">failed requests</span>
          </div>
        </div>

        {/* Bytes Transferred */}
        <div className="glass-panel p-5 rounded-xl border border-slate-200/60 flex flex-col justify-between shadow-sm">
          <div className="flex items-center justify-between text-slate-500">
            <span className="text-[10px] font-bold uppercase tracking-wider">Transferred</span>
            <ArrowUpRight className="w-4 h-4 text-sky-500" />
          </div>
          <div className="mt-2.5">
            <h4 className="text-2xl font-bold text-slate-800 font-mono">{stats.totalBytes}</h4>
            <span className="text-[10px] text-slate-500 font-medium">total payload size</span>
          </div>
        </div>
      </div>

      {/* Main Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Throughput chart */}
        <div className="glass-panel p-6 rounded-xl border border-slate-200/60 shadow-sm">
          <h4 className="text-xs font-bold text-slate-550 uppercase tracking-wider mb-4">Throughput Over Time</h4>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={throughputData}>
                <defs>
                  <linearGradient id="colorRps" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.12}/>
                    <stop offset="95%" stopColor="#10b981" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#E2E8F0" vertical={false} />
                <XAxis dataKey="second" stroke="#64748B" fontSize={11} tickLine={false} />
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
                <Area type="monotone" dataKey="RPS" stroke="#10b981" strokeWidth={2} fillOpacity={1} fill="url(#colorRps)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Latency over time */}
        <div className="glass-panel p-6 rounded-xl border border-slate-200/60 shadow-sm">
          <h4 className="text-xs font-bold text-slate-550 uppercase tracking-wider mb-4">P95 Latency Over Time</h4>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={latencyData}>
                <defs>
                  <linearGradient id="colorLatency" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#f59e0b" stopOpacity={0.12}/>
                    <stop offset="95%" stopColor="#f59e0b" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#E2E8F0" vertical={false} />
                <XAxis dataKey="second" stroke="#64748B" fontSize={11} tickLine={false} />
                <YAxis stroke="#64748B" fontSize={11} tickLine={false} unit="ms" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#FFFFFF',
                    borderColor: '#E2E8F0',
                    borderRadius: '10px',
                    color: '#1E293B',
                    fontSize: '12px',
                  }}
                />
                <Area type="monotone" dataKey="P95" stroke="#f59e0b" strokeWidth={2} fillOpacity={1} fill="url(#colorLatency)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Final test profiles charts (only shown if not running) */}
        {!isRunning && metrics && (
          <>
            {/* Latency Percentiles Histogram */}
            <div className="glass-panel p-6 rounded-xl border border-slate-200/60 shadow-sm">
              <h4 className="text-xs font-bold text-slate-550 uppercase tracking-wider mb-4">Latency Profile Distribution</h4>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={latencyProfileData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#E2E8F0" vertical={false} />
                    <XAxis dataKey="name" stroke="#64748B" fontSize={11} tickLine={false} />
                    <YAxis stroke="#64748B" fontSize={11} tickLine={false} unit="ms" />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: '#FFFFFF',
                        borderColor: '#E2E8F0',
                        borderRadius: '10px',
                        color: '#1E293B',
                        fontSize: '12px',
                      }}
                      formatter={(value) => [`${value} ms`, 'Latency']}
                    />
                    <Bar dataKey="value" fill="#6366f1" radius={[4, 4, 0, 0]}>
                      {latencyProfileData.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>

            {/* Status Codes Distribution */}
            <div className="glass-panel p-6 rounded-xl border border-slate-200/60 shadow-sm">
              <h4 className="text-xs font-bold text-slate-550 uppercase tracking-wider mb-4">HTTP Status Distribution</h4>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={statusCodeData} layout="vertical">
                    <CartesianGrid strokeDasharray="3 3" stroke="#E2E8F0" horizontal={false} />
                    <XAxis type="number" stroke="#64748B" fontSize={11} tickLine={false} />
                    <YAxis type="category" dataKey="name" stroke="#64748B" fontSize={11} tickLine={false} />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: '#FFFFFF',
                        borderColor: '#E2E8F0',
                        borderRadius: '10px',
                        color: '#1E293B',
                        fontSize: '12px',
                      }}
                    />
                    <Bar dataKey="value" fill="#10b981" radius={[0, 4, 4, 0]}>
                      {statusCodeData.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[(index + 1) % COLORS.length]} />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
