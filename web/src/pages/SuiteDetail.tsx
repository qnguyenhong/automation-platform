import { useParams, useNavigate, Link } from 'react-router-dom'
import { useSuite, useCases, useCreateCase, useUpdateCase, useDeleteCase } from '../hooks/useSuites'
import { useTriggerRun, useRuns } from '../hooks/useRuns'
import { Play, Plus, Upload, Trash2, Edit3, X, AlertCircle } from 'lucide-react'
import { useState, useEffect, useMemo } from 'react'
import { TestCase } from '../api/suites'
import K6Guide from '../components/load/K6Guide'

interface CaseFormData {
  name: string
  description: string
  enabled: boolean
  sort_order: number
  tags: string
  config: {
    method: string
    url: string
    headers: Record<string, string>
    body: string
    load: {
      vus: number
      duration: string
      ramp_up: string
      ramp_down: string
      rate_limit_rps: number
    }
    assertions: Array<{ type: string; expected: any }>
  }
}

const emptyForm: CaseFormData = {
  name: '',
  description: '',
  enabled: true,
  sort_order: 1,
  tags: '',
  config: {
    method: 'GET',
    url: '',
    headers: {},
    body: '',
    load: {
      vus: 2,
      duration: '15s',
      ramp_up: '2s',
      ramp_down: '2s',
      rate_limit_rps: 5,
    },
    assertions: [
      { type: 'p95_response_time', expected: 500 },
      { type: 'error_rate', expected: 0.05 },
    ],
  },
}

export default function SuiteDetail() {
  const { projectId, suiteId } = useParams<{ projectId: string; suiteId: string }>()
  const navigate = useNavigate()
  const { data: suite } = useSuite(projectId!, suiteId!)
  const { data: cases } = useCases(projectId!, suiteId!)
  const triggerRun = useTriggerRun(projectId!, suiteId!)
  const { data: runs } = useRuns(projectId!)

  const createCase = useCreateCase(projectId!, suiteId!)
  const updateCase = useUpdateCase(projectId!, suiteId!)
  const deleteCase = useDeleteCase(projectId!, suiteId!)

  const [showModal, setShowModal] = useState(false)
  const [editingCase, setEditingCase] = useState<TestCase | null>(null)
  const [formData, setFormData] = useState<CaseFormData>(emptyForm)
  const [selectedPreset, setSelectedPreset] = useState<string>('smoke')

  const [headerRows, setHeaderRows] = useState<Array<{ key: string; value: string }>>([])
  const [assertionRows, setAssertionRows] = useState<Array<{ type: string; expected: string }>>([])

  useEffect(() => {
    if (editingCase) {
      const tagsStr = editingCase.tags ? editingCase.tags.join(', ') : ''
      const cfg = editingCase.config || {}
      
      const formHeaders = cfg.headers || {}
      const headerRowsParsed = Object.keys(formHeaders).map((k) => ({
        key: k,
        value: String(formHeaders[k]),
      }))

      const assertionsList = cfg.assertions || []
      const assertionRowsParsed = assertionsList.map((a: any) => ({
        type: a.type || 'p95_response_time',
        expected: String(a.expected ?? ''),
      }))

      const loadVus = cfg.load?.vus ?? 2
      const loadDuration = cfg.load?.duration ?? '15s'
      const loadRampUp = cfg.load?.ramp_up ?? '2s'
      const loadRampDown = cfg.load?.ramp_down ?? '2s'
      const loadRps = cfg.load?.rate_limit_rps ?? 5

      setFormData({
        name: editingCase.name || '',
        description: editingCase.description || '',
        enabled: editingCase.enabled ?? true,
        sort_order: editingCase.sort_order || 1,
        tags: tagsStr,
        config: {
          method: cfg.method || 'GET',
          url: cfg.url || '',
          headers: formHeaders,
          body: cfg.body ? (typeof cfg.body === 'object' ? JSON.stringify(cfg.body, null, 2) : String(cfg.body)) : '',
          load: {
            vus: loadVus,
            duration: loadDuration,
            ramp_up: loadRampUp,
            ramp_down: loadRampDown,
            rate_limit_rps: loadRps,
          },
          assertions: assertionsList,
        },
      })
      setHeaderRows(headerRowsParsed)
      setAssertionRows(assertionRowsParsed)

      // Detect matching preset
      if (loadVus === 2 && loadDuration === '15s' && loadRampUp === '2s' && loadRampDown === '2s' && loadRps === 5) {
        setSelectedPreset('smoke')
      } else if (loadVus === 20 && loadDuration === '1m' && loadRampUp === '10s' && loadRampDown === '10s' && loadRps === 0) {
        setSelectedPreset('benchmark')
      } else if (loadVus === 100 && loadDuration === '2m' && loadRampUp === '15s' && loadRampDown === '15s' && loadRps === 0) {
        setSelectedPreset('stress')
      } else if (loadVus === 50 && loadDuration === '30s' && loadRampUp === '2s' && loadRampDown === '2s' && loadRps === 0) {
        setSelectedPreset('spike')
      } else {
        setSelectedPreset('custom')
      }
    } else {
      setFormData({
        ...emptyForm,
        config: {
          ...emptyForm.config,
          load: {
            vus: 2,
            duration: '15s',
            ramp_up: '2s',
            ramp_down: '2s',
            rate_limit_rps: 5,
          }
        }
      })
      setSelectedPreset('smoke')
      setHeaderRows([])
      setAssertionRows(
        suite?.test_type === 'load'
          ? [
              { type: 'p95_response_time', expected: '500' },
              { type: 'error_rate', expected: '0.05' },
            ]
          : [{ type: 'status', expected: '200' }]
      )
    }
  }, [editingCase, showModal, suite])

  const suiteRuns = useMemo(() => {
    if (!runs) return []
    return runs
      .filter((r) => r.suite_id === suiteId)
      .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
      .slice(0, 10)
  }, [runs, suiteId])

  const applyPreset = (preset: string) => {
    setSelectedPreset(preset)
    if (preset === 'smoke') {
      setFormData(prev => ({
        ...prev,
        config: {
          ...prev.config,
          load: { vus: 2, duration: '15s', ramp_up: '2s', ramp_down: '2s', rate_limit_rps: 5 }
        }
      }))
    } else if (preset === 'benchmark') {
      setFormData(prev => ({
        ...prev,
        config: {
          ...prev.config,
          load: { vus: 20, duration: '1m', ramp_up: '10s', ramp_down: '10s', rate_limit_rps: 0 }
        }
      }))
    } else if (preset === 'stress') {
      setFormData(prev => ({
        ...prev,
        config: {
          ...prev.config,
          load: { vus: 100, duration: '2m', ramp_up: '15s', ramp_down: '15s', rate_limit_rps: 0 }
        }
      }))
    } else if (preset === 'spike') {
      setFormData(prev => ({
        ...prev,
        config: {
          ...prev.config,
          load: { vus: 50, duration: '30s', ramp_up: '2s', ramp_down: '2s', rate_limit_rps: 0 }
        }
      }))
    }
  }

  const handleLoadChange = (field: string, val: any) => {
    setFormData(prev => {
      const nextLoad = {
        ...prev.config.load,
        [field]: val
      }
      
      // Determine if nextLoad matches any preset
      let nextPreset = 'custom'
      const { vus, duration, ramp_up, ramp_down, rate_limit_rps } = nextLoad
      if (vus === 2 && duration === '15s' && ramp_up === '2s' && ramp_down === '2s' && rate_limit_rps === 5) {
        nextPreset = 'smoke'
      } else if (vus === 20 && duration === '1m' && ramp_up === '10s' && ramp_down === '10s' && rate_limit_rps === 0) {
        nextPreset = 'benchmark'
      } else if (vus === 100 && duration === '2m' && ramp_up === '15s' && ramp_down === '15s' && rate_limit_rps === 0) {
        nextPreset = 'stress'
      } else if (vus === 50 && duration === '30s' && ramp_up === '2s' && ramp_down === '2s' && rate_limit_rps === 0) {
        nextPreset = 'spike'
      }
      
      setSelectedPreset(nextPreset)
      return {
        ...prev,
        config: {
          ...prev.config,
          load: nextLoad
        }
      }
    })
  }

  const handleRun = async () => {
    try {
      const run = await triggerRun.mutateAsync({})
      if (run && run.id) {
        navigate(`/runs/${run.id}`)
      }
    } catch (err) {
      console.error('Failed to trigger run:', err)
    }
  }

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()

    // Reconstruct headers
    const finalHeaders: Record<string, string> = {}
    headerRows.forEach((row) => {
      if (row.key.trim()) {
        finalHeaders[row.key.trim()] = row.value.trim()
      }
    })

    // Reconstruct assertions
    const finalAssertions = assertionRows
      .filter((row) => row.type)
      .map((row) => {
        let parsedExpected: any = row.expected
        if (!isNaN(Number(row.expected)) && row.expected.trim() !== '') {
          parsedExpected = Number(row.expected)
        }
        return {
          type: row.type,
          expected: parsedExpected,
        }
      })

    // Parse body if JSON
    let finalBody: any = formData.config.body
    if (formData.config.body.trim()) {
      try {
        finalBody = JSON.parse(formData.config.body)
      } catch {
        finalBody = formData.config.body
      }
    }

    const payload: Partial<TestCase> = {
      name: formData.name,
      description: formData.description,
      enabled: formData.enabled,
      sort_order: Number(formData.sort_order),
      tags: formData.tags
        .split(',')
        .map((t) => t.trim())
        .filter((t) => t.length > 0),
      config: {
        method: formData.config.method,
        url: formData.config.url,
        headers: finalHeaders,
        body: finalBody || null,
        load: suite?.test_type === 'load' ? {
          vus: Number(formData.config.load.vus),
          duration: formData.config.load.duration || '15s',
          ramp_up: formData.config.load.ramp_up || '0s',
          ramp_down: formData.config.load.ramp_down || '0s',
          rate_limit_rps: Number(formData.config.load.rate_limit_rps),
        } : undefined,
        assertions: finalAssertions,
      },
    }

    try {
      if (editingCase) {
        await updateCase.mutateAsync({ caseId: editingCase.id, data: payload })
      } else {
        await createCase.mutateAsync(payload)
      }
      setShowModal(false)
    } catch (err) {
      console.error('Failed to save test case:', err)
    }
  }

  const handleDelete = async (caseId: string) => {
    if (window.confirm('Are you sure you want to delete this target endpoint?')) {
      try {
        await deleteCase.mutateAsync(caseId)
        setShowModal(false)
      } catch (err) {
        console.error('Failed to delete case:', err)
      }
    }
  }

  const addHeaderRow = () => setHeaderRows([...headerRows, { key: '', value: '' }])
  const removeHeaderRow = (index: number) => setHeaderRows(headerRows.filter((_, i) => i !== index))
  const updateHeaderRow = (index: number, field: 'key' | 'value', val: string) => {
    const next = [...headerRows]
    next[index][field] = val
    setHeaderRows(next)
  }

  const addAssertionRow = () => setAssertionRows([...assertionRows, { type: 'p95_response_time', expected: '' }])
  const removeAssertionRow = (index: number) => setAssertionRows(assertionRows.filter((_, i) => i !== index))
  const updateAssertionRow = (index: number, field: 'type' | 'expected', val: string) => {
    const next = [...assertionRows]
    next[index][field] = val
    setAssertionRows(next)
  }

  const getSuiteLabelText = (type?: string) => {
    switch (type) {
      case 'load': return 'Performance Load Simulation'
      case 'api': return 'REST API Assertion'
      case 'e2e': return 'User Journey E2E Flow'
      case 'unit': return 'Isolated Unit Check'
      default: return 'Validation Suite'
    }
  }

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold text-white">{suite?.name}</h1>
            <span className="px-2.5 py-0.5 text-xs font-semibold bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 rounded-md">
              {getSuiteLabelText(suite?.test_type)}
            </span>
          </div>
          <p className="text-slate-400 mt-1 text-sm">{suite?.description || 'Configure parameters, endpoints, and load thresholds.'}</p>
        </div>
        <button
          onClick={handleRun}
          disabled={triggerRun.isPending}
          className="flex items-center gap-2 px-4 py-2.5 bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/20 text-emerald-400 font-bold rounded-lg shadow-sm transition disabled:opacity-50 text-sm animate-in fade-in duration-300"
        >
          <Play className="w-4 h-4 fill-current animate-pulse" />
          {triggerRun.isPending ? 'Running...' : 'Execute Suite'}
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
        {/* Left Column: Endpoints Table & K6 Guide */}
        <div className="lg:col-span-2 space-y-6">
          <div className="glass-panel rounded-xl border border-slate-850 overflow-hidden">
            <div className="p-6 border-b border-slate-850 bg-slate-900/30 flex items-center justify-between">
              <h3 className="text-sm font-bold uppercase tracking-wider text-slate-400">Target Endpoints</h3>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => navigate(`/projects/${projectId}/openapi-import`)}
                  className="flex items-center gap-2 px-3.5 py-2 text-xs font-semibold bg-indigo-500/10 text-indigo-400 hover:bg-indigo-500/20 border border-indigo-500/20 rounded-lg transition"
                >
                  <Upload className="w-4 h-4" />
                  Ingest from OpenAPI Spec
                </button>
                <button
                  onClick={() => {
                    setEditingCase(null)
                    setShowModal(true)
                  }}
                  className="flex items-center gap-2 px-3.5 py-2 text-xs font-bold bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg shadow-sm transition"
                >
                  <Plus className="w-4 h-4" />
                  Define Target Endpoint
                </button>
              </div>
            </div>
            
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead className="bg-[#0B0F19]/50 text-slate-400 text-xs font-bold uppercase border-b border-slate-850">
                  <tr>
                    <th className="px-6 py-4">Endpoint Label</th>
                    <th className="px-6 py-4">Target URL</th>
                    <th className="px-6 py-4">Status</th>
                    <th className="px-6 py-4">Tags</th>
                    <th className="px-6 py-4 text-right">Configure</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-850 text-sm">
                  {cases && cases.length > 0 ? (
                    cases.map((testCase) => {
                      const cfg = testCase.config || {}
                      const method = cfg.method || 'GET'
                      const url = cfg.url || ''

                      return (
                        <tr key={testCase.id} className="hover:bg-slate-800/10 transition-colors">
                          <td className="px-6 py-4 font-semibold text-slate-200">
                            <div>
                              {testCase.name}
                              {testCase.description && (
                                <p className="text-xs text-slate-500 font-normal mt-0.5">{testCase.description}</p>
                              )}
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center gap-2 font-mono text-xs">
                              <span className={`px-2 py-0.5 rounded font-bold border ${
                                method === 'GET' ? 'bg-sky-500/10 text-sky-400 border-sky-500/20' :
                                method === 'POST' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' :
                                method === 'PUT' ? 'bg-amber-500/10 text-amber-400 border-amber-500/20' : 
                                'bg-rose-500/10 text-rose-400 border-rose-500/20'
                              }`}>
                                {method}
                              </span>
                              <span className="text-slate-300 truncate max-w-xs">{url || '—'}</span>
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold border ${
                              testCase.enabled ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-slate-800 text-slate-400 border-slate-700'
                            }`}>
                              {testCase.enabled ? 'Enabled' : 'Disabled'}
                            </span>
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex flex-wrap gap-1">
                              {testCase.tags?.map((tag) => (
                                <span key={tag} className="px-2 py-0.5 text-xs bg-slate-900 border border-slate-800 text-slate-400 rounded font-mono">
                                  {tag}
                                </span>
                              ))}
                            </div>
                          </td>
                          <td className="px-6 py-4 text-right">
                            <button
                              onClick={() => {
                                setEditingCase(testCase)
                                setShowModal(true)
                              }}
                              className="p-1.5 text-slate-400 hover:text-indigo-400 rounded-md border border-slate-800 bg-slate-900 hover:border-slate-750 transition"
                            >
                              <Edit3 className="w-3.5 h-3.5" />
                            </button>
                          </td>
                        </tr>
                      )
                    })
                  ) : (
                    <tr>
                      <td colSpan={5} className="px-6 py-12 text-center text-slate-500">
                        <AlertCircle className="w-8 h-8 mx-auto mb-2 text-slate-650" />
                        No target endpoints defined for this execution suite.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>

          {suite?.test_type === 'load' && <K6Guide />}
        </div>

        {/* Right Column: Execution History */}
        <div className="space-y-6">
          <div className="glass-panel rounded-xl border border-slate-850 overflow-hidden">
            <div className="p-5 border-b border-slate-850 bg-slate-900/30 flex items-center justify-between">
              <h3 className="text-xs font-bold uppercase tracking-wider text-slate-400">Execution History</h3>
            </div>
            <div className="p-5 space-y-3">
              {suiteRuns.length > 0 ? (
                suiteRuns.map((run) => (
                  <Link
                    key={run.id}
                    to={`/runs/${run.id}`}
                    className="block p-3.5 bg-slate-950/40 border border-slate-850 hover:border-slate-800 rounded-xl hover:bg-slate-800/10 transition group"
                  >
                    <div className="flex items-center justify-between gap-2">
                      <span className="font-mono text-xs font-bold text-indigo-400 group-hover:text-indigo-300 transition-colors">
                        #{run.id.slice(0, 8)}
                      </span>
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded text-[10px] font-bold uppercase border ${
                        run.status === 'passed' ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-400' :
                        run.status === 'failed' ? 'bg-rose-500/10 border-rose-500/20 text-rose-400' :
                        run.status === 'running' ? 'bg-blue-500/10 border-blue-500/20 text-blue-400 animate-pulse' :
                        'bg-slate-800/40 border-slate-800 text-slate-500'
                      }`}>
                        {run.status === 'passed' ? 'Successful' : run.status}
                      </span>
                    </div>
                    
                    <div className="flex items-center justify-between text-[11px] text-slate-450 mt-2.5 font-medium">
                      <span>Trigger: <strong className="text-slate-350 capitalize font-semibold">{run.trigger}</strong></span>
                      <span className="font-mono">{run.duration_ms ? `${(run.duration_ms / 1000).toFixed(1)}s` : '—'}</span>
                    </div>
                    
                    <div className="text-[10px] text-slate-500 mt-1 font-mono">
                      {new Date(run.created_at).toLocaleString()}
                    </div>
                  </Link>
                ))
              ) : (
                <div className="p-8 text-center text-slate-500 text-xs italic">
                  No execution runs recorded yet for this suite.
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Case Configuration Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-950/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-[#0B0F19] rounded-2xl shadow-xl w-full max-w-3xl max-h-[90vh] flex flex-col border border-slate-800 animate-in fade-in zoom-in-95 duration-150 overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-850 flex items-center justify-between bg-slate-900/30">
              <h2 className="text-lg font-bold text-slate-100">
                {editingCase ? 'Configure Target Endpoint' : 'Define Target Endpoint'}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-1.5 hover:bg-slate-800 rounded-lg text-slate-400 hover:text-slate-200 transition"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleSave} className="flex-1 overflow-y-auto p-6 space-y-6">
              {/* General Fields */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="md:col-span-2">
                  <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Target Endpoint Name</label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-100 focus:ring-1 focus:ring-indigo-500 transition-all"
                    placeholder="e.g. GET User Profile API"
                    required
                  />
                </div>
                <div className="md:col-span-2">
                  <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Endpoint Description</label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm h-20 resize-none text-slate-100 focus:ring-1 focus:ring-indigo-500 transition-all"
                    placeholder="Describe what this target endpoint validates..."
                  />
                </div>
                <div>
                  <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Sort Order</label>
                  <input
                    type="number"
                    value={formData.sort_order}
                    onChange={(e) => setFormData({ ...formData, sort_order: Number(e.target.value) })}
                    className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-100 focus:ring-1 focus:ring-indigo-500 transition-all"
                    required
                  />
                </div>
                <div>
                  <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Tags (comma separated)</label>
                  <input
                    type="text"
                    value={formData.tags}
                    onChange={(e) => setFormData({ ...formData, tags: e.target.value })}
                    className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-100 focus:ring-1 focus:ring-indigo-500 transition-all"
                    placeholder="load, checkout, auth"
                  />
                </div>
                <div className="flex items-center gap-2.5 mt-4 ml-1">
                  <input
                    type="checkbox"
                    id="case-enabled"
                    checked={formData.enabled}
                    onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
                    className="w-4 h-4 text-indigo-650 bg-slate-950 border-slate-850 rounded focus:ring-indigo-500"
                  />
                  <label htmlFor="case-enabled" className="text-sm font-semibold text-slate-350">Enabled</label>
                </div>
              </div>

              {/* Endpoint configuration */}
              <div className="border-t border-slate-850 pt-6 space-y-4">
                <h3 className="font-bold text-slate-300 text-sm">Request Configuration</h3>
                <div className="flex gap-3">
                  <select
                    value={formData.config.method}
                    onChange={(e) => setFormData({
                      ...formData,
                      config: { ...formData.config, method: e.target.value }
                    })}
                    className="px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm font-bold text-slate-100 focus:ring-1 focus:ring-indigo-500 transition cursor-pointer"
                  >
                    <option value="GET">GET</option>
                    <option value="POST">POST</option>
                    <option value="PUT">PUT</option>
                    <option value="DELETE">DELETE</option>
                  </select>
                  <input
                    type="url"
                    value={formData.config.url}
                    onChange={(e) => setFormData({
                      ...formData,
                      config: { ...formData.config, url: e.target.value }
                    })}
                    className="flex-1 px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm font-mono text-slate-100 focus:ring-1 focus:ring-indigo-500 transition-all"
                    placeholder="https://api.example.com/endpoint"
                    required
                  />
                </div>

                {/* Headers Grid */}
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <label className="text-xs font-bold uppercase tracking-wider text-slate-400">Headers</label>
                    <button
                      type="button"
                      onClick={addHeaderRow}
                      className="text-xs text-indigo-400 hover:text-indigo-300 font-bold"
                    >
                      + Add Header
                    </button>
                  </div>
                  {headerRows.map((row, idx) => (
                    <div key={idx} className="flex gap-2 items-center">
                      <input
                        type="text"
                        value={row.key}
                        onChange={(e) => updateHeaderRow(idx, 'key', e.target.value)}
                        placeholder="Key (e.g. Authorization)"
                        className="flex-1 px-3 py-2 bg-slate-950/85 border border-slate-850 rounded-xl text-xs font-mono text-slate-100"
                      />
                      <input
                        type="text"
                        value={row.value}
                        onChange={(e) => updateHeaderRow(idx, 'value', e.target.value)}
                        placeholder="Value"
                        className="flex-1 px-3 py-2 bg-slate-950/85 border border-slate-850 rounded-xl text-xs font-mono text-slate-100"
                      />
                      <button
                        type="button"
                        onClick={() => removeHeaderRow(idx)}
                        className="p-1.5 hover:text-rose-500 text-slate-500 border border-slate-850 bg-slate-900 rounded-md transition"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  ))}
                </div>

                {/* Body Area */}
                {formData.config.method !== 'GET' && (
                  <div>
                    <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Request Body (JSON)</label>
                    <textarea
                      value={formData.config.body}
                      onChange={(e) => setFormData({
                        ...formData,
                        config: { ...formData.config, body: e.target.value }
                      })}
                      className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-xs font-mono h-24 text-slate-100 focus:ring-1 focus:ring-indigo-500 transition-all resize-none"
                      placeholder='{ "key": "value" }'
                    />
                  </div>
                )}
              </div>

              {/* Load Configuration */}
              {suite?.test_type === 'load' && (
                <div className="border-t border-slate-850 pt-6 space-y-4">
                  <h3 className="font-bold text-slate-300 text-sm">Concurrency Load Settings</h3>
                  
                  {/* Presets Grid */}
                  <div className="space-y-2">
                    <label className="block text-xs font-bold uppercase tracking-wider text-slate-400">
                      Select Load Profile Preset
                    </label>
                    <div className="grid grid-cols-2 md:grid-cols-5 gap-2">
                      {(['smoke', 'benchmark', 'stress', 'spike', 'custom'] as const).map((preset) => (
                        <button
                          key={preset}
                          type="button"
                          onClick={() => {
                            if (preset !== 'custom') {
                              applyPreset(preset)
                            } else {
                              setSelectedPreset('custom')
                            }
                          }}
                          className={`px-3 py-2 text-xs font-semibold rounded-lg border transition capitalize text-center ${
                            selectedPreset === preset
                              ? 'bg-indigo-500/10 border-indigo-500/30 text-indigo-400 shadow-sm font-bold'
                              : 'bg-slate-950 border-slate-850 hover:bg-slate-900 text-slate-400 hover:text-slate-200'
                          }`}
                        >
                          {preset === 'smoke' ? 'Smoke Check' :
                           preset === 'benchmark' ? 'Standard Load' :
                           preset === 'stress' ? 'Stress Load' :
                           preset === 'spike' ? 'Spike Surge' : 'Custom Profile'}
                        </button>
                      ))}
                    </div>
                    
                    <div className="p-3 bg-slate-900/30 rounded-xl border border-slate-850 text-xs text-slate-450 leading-relaxed font-semibold">
                      {selectedPreset === 'smoke' && (
                        <p><strong>Smoke Check:</strong> Validates endpoint health and basic responsiveness with minimal traffic (2 Simulated VUs, 5 RPS limit, 15s duration).</p>
                      )}
                      {selectedPreset === 'benchmark' && (
                        <p><strong>Standard Load:</strong> Simulates standard client load to establish performance benchmarks under realistic usage (20 Simulated VUs, uncapped RPS, 1m duration).</p>
                      )}
                      {selectedPreset === 'stress' && (
                        <p><strong>Stress Load:</strong> Pushes the endpoint to its upper limits to discover bottlenecks, memory leaks, or scaling limits (100 Simulated VUs, uncapped RPS, 2m duration).</p>
                      )}
                      {selectedPreset === 'spike' && (
                        <p><strong>Spike Surge:</strong> Evaluates system resilience and recovery against sudden surges in visitor volume (50 Simulated VUs, uncapped RPS, 30s duration with 2s rapid ramp).</p>
                      )}
                      {selectedPreset === 'custom' && (
                        <p><strong>Custom Profile:</strong> Configure Simulated Virtual Users, Simulation Duration, and Ramp-up/down criteria manually below.</p>
                      )}
                    </div>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                      <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Simulated Virtual Users (VUs)</label>
                      <input
                        type="number"
                        value={formData.config.load.vus}
                        onChange={(e) => handleLoadChange('vus', Number(e.target.value))}
                        className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-100"
                        min="1"
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Simulation Duration</label>
                      <input
                        type="text"
                        value={formData.config.load.duration}
                        onChange={(e) => handleLoadChange('duration', e.target.value)}
                        className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-100"
                        placeholder="e.g. 30s, 2m, 1h"
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Rate Limiter (RPS, 0 = Uncapped)</label>
                      <input
                        type="number"
                        value={formData.config.load.rate_limit_rps}
                        onChange={(e) => handleLoadChange('rate_limit_rps', Number(e.target.value))}
                        className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-100"
                        min="0"
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Ramp Up Phase</label>
                      <input
                        type="text"
                        value={formData.config.load.ramp_up}
                        onChange={(e) => handleLoadChange('ramp_up', e.target.value)}
                        className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-100"
                        placeholder="e.g. 5s"
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Ramp Down Phase</label>
                      <input
                        type="text"
                        value={formData.config.load.ramp_down}
                        onChange={(e) => handleLoadChange('ramp_down', e.target.value)}
                        className="w-full px-3.5 py-2.5 bg-slate-950/80 border border-slate-850 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-100"
                        placeholder="e.g. 5s"
                        required
                      />
                    </div>
                  </div>
                </div>
              )}

              {/* Assertions */}
              <div className="border-t border-slate-850 pt-6 space-y-4">
                <div className="flex items-center justify-between">
                  <h3 className="font-bold text-slate-350 text-sm">Service Level Objectives (SLOs)</h3>
                  <button
                    type="button"
                    onClick={addAssertionRow}
                    className="text-xs text-indigo-400 hover:text-indigo-300 font-bold"
                  >
                    + Add SLO/Assertion
                  </button>
                </div>
                <div className="space-y-2.5">
                  {assertionRows.map((row, idx) => (
                    <div key={idx} className="flex gap-2 items-center">
                      <select
                        value={row.type}
                        onChange={(e) => updateAssertionRow(idx, 'type', e.target.value)}
                        className="px-3.5 py-2 bg-slate-950/90 border border-slate-850 rounded-xl text-xs text-slate-200 focus:outline-none focus:border-indigo-500 font-semibold cursor-pointer"
                      >
                        {suite?.test_type === 'load' ? (
                          <>
                            <option value="p95_response_time">P95 Response Time (ms) &lt;=</option>
                            <option value="p99_response_time">P99 Response Time (ms) &lt;=</option>
                            <option value="error_rate">Max Error Rate (fraction 0-1) &lt;=</option>
                            <option value="throughput_min">Min Throughput (RPS) &gt;=</option>
                            <option value="avg_response_time">Avg Response Time (ms) &lt;=</option>
                          </>
                        ) : (
                          <>
                            <option value="status">Status Equals</option>
                            <option value="response_time">Response Time (ms) &lt;=</option>
                            <option value="header">Header Value</option>
                            <option value="json_path">JSON Body Property</option>
                          </>
                        )}
                      </select>
                      <input
                        type="text"
                        value={row.expected}
                        onChange={(e) => updateAssertionRow(idx, 'expected', e.target.value)}
                        placeholder="Expected value (e.g. 500, 0.02)"
                        className="flex-1 px-3.5 py-2 bg-slate-950/85 border border-slate-850 rounded-xl text-xs font-mono text-slate-100"
                        required
                      />
                      <button
                        type="button"
                        onClick={() => removeAssertionRow(idx)}
                        className="p-1.5 hover:text-rose-500 text-slate-500 border border-slate-850 bg-slate-900 rounded-md transition"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            </form>

            <div className="px-6 py-4 bg-slate-900/30 border-t border-slate-850 flex items-center justify-between rounded-b-2xl">
              {editingCase ? (
                <button
                  type="button"
                  onClick={() => handleDelete(editingCase.id)}
                  className="flex items-center gap-1.5 text-sm font-semibold text-rose-550 hover:text-rose-400 transition"
                >
                  <Trash2 className="w-4 h-4" />
                  Delete Endpoint
                </button>
              ) : (
                <div />
              )}
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-4 py-2 border border-slate-800 text-slate-400 text-sm font-semibold rounded-xl hover:bg-slate-850 hover:text-slate-200 transition"
                >
                  Cancel
                </button>
                <button
                  onClick={handleSave}
                  className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-bold rounded-xl shadow-md shadow-indigo-500/10 transition"
                >
                  Save Configuration
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
