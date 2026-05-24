import { useState, useMemo } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeft, ArrowRight, Check, Loader2 } from 'lucide-react'
import { useParseSpec, useImportEndpoints } from '../hooks/useOpenAPI'
import { useVariableStore } from '../store/variables'
import { SpecUploader } from '../components/openapi/SpecUploader'
import { EndpointSelector } from '../components/openapi/EndpointSelector'
import { TestDataForm } from '../components/openapi/TestDataForm'
import type { ParsedEndpoint, ImportEndpointConfig } from '../types/openapi'

const STEPS = [
  { id: 1, title: 'Upload Spec', description: 'Upload or paste your OpenAPI specification' },
  { id: 2, title: 'Select Endpoints', description: 'Choose which APIs to test' },
  { id: 3, title: 'Configure Tests', description: 'Set up test data and assertions' },
  { id: 4, title: 'Review & Import', description: 'Review and save test cases' },
]

export default function OpenAPIImport() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const [step, setStep] = useState(1)

  const [endpoints, setEndpoints] = useState<ParsedEndpoint[]>([])
  const [specInfo, setSpecInfo] = useState<{ title: string; version: string } | null>(null)

  const [selectedKeys, setSelectedKeys] = useState<Set<string>>(new Set())
  const [configs, setConfigs] = useState<Record<string, ImportEndpointConfig>>({})
  const [suiteName, setSuiteName] = useState('')

  const { captured, setCaptured } = useVariableStore()

  const parseMutation = useParseSpec()
  const importMutation = useImportEndpoints()

  const handleParsed = async (content: string) => {
    try {
      const result = await parseMutation.mutateAsync({ content })
      setEndpoints(result.endpoints)
      setSpecInfo(result.spec)
      setStep(2)
    } catch {
      alert('Failed to parse OpenAPI spec')
    }
  }

  const handleToggleEndpoint = (key: string) => {
    setSelectedKeys((prev) => {
      const next = new Set(prev)
      if (next.has(key)) next.delete(key)
      else next.add(key)
      return next
    })
  }

  const handleSelectAll = () => {
    setSelectedKeys(new Set(endpoints.map((ep) => `${ep.method} ${ep.path}`)))
  }

  const handleDeselectAll = () => {
    setSelectedKeys(new Set())
  }

  const selectedEndpoints = useMemo(
    () =>
      endpoints.filter((ep) => selectedKeys.has(`${ep.method} ${ep.path}`)),
    [endpoints, selectedKeys]
  )

  const handleImport = async () => {
    if (!projectId || !suiteName) return

    try {
      const result = await importMutation.mutateAsync({
        project_id: projectId,
        suite_name: suiteName,
        endpoints: selectedEndpoints.map((ep) => {
          const key = `${ep.method} ${ep.path}`
          return configs[key] || {
            method: ep.method,
            path: ep.path,
            headers: { 'Content-Type': 'application/json' },
            params: {},
            body: {},
            assertions: [{ type: 'status' as const, expected: 200 }],
          }
        }),
      })

      navigate(`/projects/${projectId}/suites/${result.suite_id}`)
    } catch {
      alert('Failed to import endpoints')
    }
  }

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4 border-b border-slate-200 pb-5">
        <button
          onClick={() => navigate(-1)}
          className="p-2 bg-white border border-slate-200 text-slate-500 hover:text-slate-700 hover:border-slate-300 rounded-lg transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </button>
        <div>
          <h1 className="text-2xl font-bold text-slate-800">Import from OpenAPI</h1>
          <p className="text-slate-500 text-xs mt-0.5">
            Import API endpoints and create target validation structures automatically.
          </p>
        </div>
      </div>

      {/* Progress Steps Indicators */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 p-4 glass-panel rounded-2xl border border-slate-200">
        {STEPS.map((s) => (
          <div key={s.id} className="flex items-center gap-3">
            <div
              className={`w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold border transition ${
                step > s.id
                  ? 'bg-emerald-50 border-emerald-300 text-emerald-600'
                  : step === s.id
                  ? 'bg-indigo-50 border-indigo-300 text-indigo-600 animate-pulse font-black'
                  : 'bg-slate-50 border-slate-200 text-slate-400'
              }`}
            >
              {step > s.id ? <Check className="h-4 w-4" /> : s.id}
            </div>
            <div className="text-left">
              <span
                className={`block text-xs font-semibold uppercase tracking-wider ${
                  step >= s.id ? 'text-slate-700 font-bold' : 'text-slate-400'
                }`}
              >
                {s.title}
              </span>
              <span className="text-[10px] text-slate-400 leading-none hidden md:block mt-0.5">{s.description}</span>
            </div>
          </div>
        ))}
      </div>

      <div className="glass-panel rounded-2xl p-6 border border-slate-200">
        {step === 1 && (
          <div className="space-y-4">
            <h2 className="text-base font-bold text-slate-800">{STEPS[0].title}</h2>
            <p className="text-xs text-slate-500">{STEPS[0].description}</p>
            <SpecUploader onParsed={handleParsed} isLoading={parseMutation.isPending} />
          </div>
        )}

        {step === 2 && (
          <div className="space-y-4">
            <h2 className="text-base font-bold text-slate-800">{STEPS[1].title}</h2>
            <p className="text-xs text-indigo-600 font-mono font-semibold">
              {specInfo?.title} v{specInfo?.version} — {endpoints.length} Endpoints Discovered
            </p>
            <EndpointSelector
              endpoints={endpoints}
              selected={selectedKeys}
              onToggle={handleToggleEndpoint}
              onSelectAll={handleSelectAll}
              onDeselectAll={handleDeselectAll}
            />
          </div>
        )}

        {step === 3 && (
          <div className="space-y-4">
            <h2 className="text-base font-bold text-slate-800">{STEPS[2].title}</h2>
            <p className="text-xs text-slate-500">
              Configure parameters and default values for the {selectedEndpoints.length} selected target endpoints.
            </p>
            <TestDataForm
              endpoints={selectedEndpoints}
              configs={configs}
              onConfigChange={(key, config) =>
                setConfigs((prev) => ({ ...prev, [key]: config }))
              }
              capturedVars={captured}
              onCapture={setCaptured}
            />
          </div>
        )}

        {step === 4 && (
          <div className="space-y-4">
            <h2 className="text-base font-bold text-slate-800">{STEPS[3].title}</h2>
            <p className="text-xs text-slate-500">
              Review and establish validation suite configs before final ingest.
            </p>

            <div className="space-y-4 pt-2">
              <div>
                <label className="block text-xs font-bold uppercase tracking-wider text-slate-500 mb-1.5">
                  Validation Suite Name
                </label>
                <input
                  type="text"
                  value={suiteName}
                  onChange={(e) => setSuiteName(e.target.value)}
                  placeholder={`API Tests - ${specInfo?.title || 'My API'}`}
                  className="w-full px-3.5 py-2.5 bg-white border border-slate-200 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-800 focus:ring-1 focus:ring-indigo-500 transition-all placeholder-slate-400"
                />
              </div>

              <div className="p-4 bg-slate-50 border border-slate-200 rounded-xl">
                <h3 className="text-xs font-bold uppercase tracking-wider text-slate-500 mb-2">Ingest Summary</h3>
                <ul className="space-y-1.5 text-xs text-slate-600 leading-relaxed font-semibold">
                  <li>Spec Reference: <span className="text-indigo-600 font-mono">{specInfo?.title} v{specInfo?.version}</span></li>
                  <li>Target Count: <span className="text-slate-800">{selectedEndpoints.length} target endpoints</span></li>
                </ul>
              </div>

              <div className="space-y-2">
                <h3 className="text-xs font-bold uppercase tracking-wider text-slate-500">Endpoints to Import</h3>
                <div className="max-h-64 overflow-y-auto space-y-1.5 pr-2">
                  {selectedEndpoints.map((ep) => (
                    <div
                      key={`${ep.method} ${ep.path}`}
                      className="flex items-center gap-2 text-xs p-2 bg-white border border-slate-200 rounded-lg"
                    >
                      <span
                        className={`px-2 py-0.5 rounded text-[10px] font-mono font-bold border ${
                          ep.method === 'GET' ? 'bg-sky-50 text-sky-600 border-sky-200' :
                          ep.method === 'POST' ? 'bg-emerald-50 text-emerald-600 border-emerald-200' :
                          ep.method === 'PUT' ? 'bg-amber-50 text-amber-600 border-amber-200' :
                          'bg-rose-50 text-rose-600 border-rose-200'
                        }`}
                      >
                        {ep.method}
                      </span>
                      <span className="font-mono text-slate-700">{ep.path}</span>
                      {ep.summary && (
                        <span className="text-slate-400 font-medium">— {ep.summary}</span>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      <div className="flex justify-between items-center">
        <button
          onClick={() => setStep((s) => Math.max(1, s - 1))}
          disabled={step === 1}
          className="px-4 py-2 border border-slate-200 text-slate-600 text-sm font-semibold rounded-xl hover:bg-slate-50 hover:text-slate-800 transition disabled:opacity-50"
        >
          Back
        </button>

        {step < 4 ? (
          <button
            onClick={() => setStep((s) => Math.min(4, s + 1))}
            disabled={step === 2 && selectedKeys.size === 0}
            className="px-4 py-2 bg-indigo-600 text-white rounded-xl text-sm font-semibold hover:bg-indigo-700 shadow-sm transition disabled:opacity-50 flex items-center gap-2"
          >
            Next Step
            <ArrowRight className="h-4 w-4" />
          </button>
        ) : (
          <button
            onClick={handleImport}
            disabled={!suiteName || importMutation.isPending}
            className="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-sm font-bold shadow-sm transition disabled:opacity-50 flex items-center gap-2"
          >
            {importMutation.isPending ? (
              <Loader2 className="mr-1 h-4 w-4 animate-spin" />
            ) : (
              <Check className="mr-1 h-4 w-4" />
            )}
            Import {selectedEndpoints.length} Targets
          </button>
        )}
      </div>
    </div>
  )
}
