import { useState, useMemo } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeft, ArrowRight, Check, Loader2 } from 'lucide-react'
import { useParseSpec, useImportEndpoints } from '@/hooks/useOpenAPI'
import { useVariableStore } from '@/store/variables'
import { SpecUploader } from '@/components/openapi/SpecUploader'
import { EndpointSelector } from '@/components/openapi/EndpointSelector'
import { TestDataForm } from '@/components/openapi/TestDataForm'
import type { ParsedEndpoint, ImportEndpointConfig } from '@/types/openapi'

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

  const [specContent, setSpecContent] = useState<string | null>(null)
  const [endpoints, setEndpoints] = useState<ParsedEndpoint[]>([])
  const [specInfo, setSpecInfo] = useState<{ title: string; version: string } | null>(null)

  const [selectedKeys, setSelectedKeys] = useState<Set<string>>(new Set())
  const [configs, setConfigs] = useState<Record<string, ImportEndpointConfig>>({})
  const [suiteName, setSuiteName] = useState('')

  const { captured, setCaptured } = useVariableStore()

  const parseMutation = useParseSpec()
  const importMutation = useImportEndpoints()

  const handleParsed = async (content: string) => {
    setSpecContent(content)
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
      <div className="flex items-center gap-4">
        <button
          onClick={() => navigate(-1)}
          className="p-2 hover:bg-muted rounded-lg transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </button>
        <div>
          <h1 className="text-2xl font-bold">Import from OpenAPI</h1>
          <p className="text-muted-foreground">
            Import API endpoints and create test cases automatically
          </p>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {STEPS.map((s, i) => (
          <div key={s.id} className="flex items-center gap-2">
            <div
              className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium ${
                step > s.id
                  ? 'bg-primary text-primary-foreground'
                  : step === s.id
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-muted text-muted-foreground'
              }`}
            >
              {step > s.id ? <Check className="h-4 w-4" /> : s.id}
            </div>
            <span
              className={`text-sm font-medium ${
                step >= s.id ? 'text-foreground' : 'text-muted-foreground'
              }`}
            >
              {s.title}
            </span>
            {i < STEPS.length - 1 && (
              <div className="w-8 h-px bg-muted-foreground/25" />
            )}
          </div>
        ))}
      </div>

      <div className="border rounded-lg p-6">
        {step === 1 && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">{STEPS[0].title}</h2>
            <p className="text-muted-foreground">{STEPS[0].description}</p>
            <SpecUploader onParsed={handleParsed} isLoading={parseMutation.isPending} />
          </div>
        )}

        {step === 2 && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">{STEPS[1].title}</h2>
            <p className="text-muted-foreground">
              {specInfo?.title} v{specInfo?.version} — {endpoints.length} endpoints found
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
            <h2 className="text-lg font-semibold">{STEPS[2].title}</h2>
            <p className="text-muted-foreground">
              Configure test data for {selectedEndpoints.length} selected endpoints
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
            <h2 className="text-lg font-semibold">{STEPS[3].title}</h2>
            <p className="text-muted-foreground">
              Review your configuration before importing
            </p>

            <div className="space-y-4">
              <div>
                <label className="text-sm font-medium">Suite Name</label>
                <input
                  type="text"
                  value={suiteName}
                  onChange={(e) => setSuiteName(e.target.value)}
                  placeholder={`API Tests - ${specInfo?.title || 'My API'}`}
                  className="w-full mt-1 px-3 py-2 border rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                />
              </div>

              <div className="p-4 bg-muted rounded-lg">
                <h3 className="font-medium mb-2">Summary</h3>
                <ul className="space-y-1 text-sm">
                  <li>Spec: {specInfo?.title} v{specInfo?.version}</li>
                  <li>Endpoints: {selectedEndpoints.length}</li>
                  <li>Test cases to create: {selectedEndpoints.length}</li>
                </ul>
              </div>

              <div className="space-y-2">
                <h3 className="font-medium text-sm">Selected Endpoints</h3>
                {selectedEndpoints.map((ep) => (
                  <div
                    key={`${ep.method} ${ep.path}`}
                    className="flex items-center gap-2 text-sm"
                  >
                    <span
                      className={`px-2 py-0.5 rounded text-xs font-mono font-bold ${
                        ep.method === 'GET'
                          ? 'bg-green-100 text-green-800'
                          : ep.method === 'POST'
                          ? 'bg-blue-100 text-blue-800'
                          : ep.method === 'PUT'
                          ? 'bg-yellow-100 text-yellow-800'
                          : ep.method === 'DELETE'
                          ? 'bg-red-100 text-red-800'
                          : 'bg-gray-100 text-gray-800'
                      }`}
                    >
                      {ep.method}
                    </span>
                    <span className="font-mono">{ep.path}</span>
                    {ep.summary && (
                      <span className="text-muted-foreground">— {ep.summary}</span>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>

      <div className="flex justify-between">
        <button
          onClick={() => setStep((s) => Math.max(1, s - 1))}
          disabled={step === 1}
          className="px-4 py-2 border rounded-lg text-sm font-medium hover:bg-muted transition-colors disabled:opacity-50"
        >
          Back
        </button>

        {step < 4 ? (
          <button
            onClick={() => setStep((s) => Math.min(4, s + 1))}
            disabled={step === 2 && selectedKeys.size === 0}
            className="px-4 py-2 bg-primary text-primary-foreground rounded-lg text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
          >
            Next
            <ArrowRight className="inline-block ml-2 h-4 w-4" />
          </button>
        ) : (
          <button
            onClick={handleImport}
            disabled={!suiteName || importMutation.isPending}
            className="px-6 py-2 bg-primary text-primary-foreground rounded-lg text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
          >
            {importMutation.isPending ? (
              <Loader2 className="inline-block mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Check className="inline-block mr-2 h-4 w-4" />
            )}
            Import {selectedEndpoints.length} Endpoints
          </button>
        )}
      </div>
    </div>
  )
}
