import { useState } from 'react'
import { Plus, Trash2 } from 'lucide-react'
import type { AssertionConfig } from '@/types/openapi'

interface AssertionBuilderProps {
  assertions: AssertionConfig[]
  onChange: (assertions: AssertionConfig[]) => void
  capturedVars: Record<string, unknown>
  onCapture?: (name: string, value: unknown) => void
}

const ASSERTION_TYPES = [
  { value: 'status', label: 'Status Code', needsPath: false },
  { value: 'response_time', label: 'Response Time (ms)', needsPath: false },
  { value: 'header', label: 'Header Value', needsPath: true },
  { value: 'json_path', label: 'JSON Path', needsPath: true },
]

export function AssertionBuilder({
  assertions,
  onChange,
  capturedVars,
  onCapture: _onCapture,
}: AssertionBuilderProps) {
  const [showCapture, setShowCapture] = useState<Record<number, boolean>>({})

  const addAssertion = () => {
    onChange([
      ...assertions,
      { type: 'status', expected: 200 },
    ])
  }

  const removeAssertion = (index: number) => {
    onChange(assertions.filter((_, i) => i !== index))
  }

  const updateAssertion = (index: number, updates: Partial<AssertionConfig>) => {
    onChange(
      assertions.map((a, i) => (i === index ? { ...a, ...updates } : a))
    )
  }

  const toggleCapture = (index: number) => {
    setShowCapture((prev) => ({ ...prev, [index]: !prev[index] }))
    if (showCapture[index]) {
      updateAssertion(index, { capture_as: undefined })
    }
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <label className="text-sm font-semibold text-slate-700">Assertions</label>
        <button
          onClick={addAssertion}
          className="flex items-center gap-1 text-sm text-indigo-600 hover:text-indigo-700 font-semibold transition-colors"
        >
          <Plus className="h-4 w-4" />
          Add
        </button>
      </div>

      {assertions.length === 0 && (
        <p className="text-sm text-slate-400">No assertions added yet</p>
      )}

      {assertions.map((assertion, index) => {
        const typeConfig = ASSERTION_TYPES.find((t) => t.value === assertion.type)
        return (
          <div key={index} className="p-3 border border-slate-200 rounded-lg space-y-2 bg-slate-50">
            <div className="flex items-center gap-2">
              <select
                value={assertion.type}
                onChange={(e) =>
                  updateAssertion(index, {
                    type: e.target.value as AssertionConfig['type'],
                    path: undefined,
                  })
                }
                className="px-2 py-1.5 border border-slate-200 rounded text-sm bg-white text-slate-700 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              >
                {ASSERTION_TYPES.map((t) => (
                  <option key={t.value} value={t.value}>
                    {t.label}
                  </option>
                ))}
              </select>

              {typeConfig?.needsPath && (
                <input
                  type="text"
                  value={assertion.path || ''}
                  onChange={(e) => updateAssertion(index, { path: e.target.value })}
                  placeholder={
                    assertion.type === 'header' ? 'Header name' : '$.data.id'
                  }
                  className="flex-1 px-2 py-1.5 border border-slate-200 rounded text-sm bg-white font-mono text-slate-700 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              )}

              <input
                type="text"
                value={String(assertion.expected ?? '')}
                onChange={(e) => {
                  let val: unknown = e.target.value
                  if (assertion.type === 'status' || assertion.type === 'response_time') {
                    val = Number(e.target.value)
                  }
                  updateAssertion(index, { expected: val })
                }}
                placeholder="Expected"
                className="w-32 px-2 py-1.5 border border-slate-200 rounded text-sm bg-white text-slate-700 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />

              <button
                onClick={() => toggleCapture(index)}
                className={`px-2 py-1.5 text-xs rounded font-semibold transition-colors ${
                  showCapture[index]
                    ? 'bg-indigo-600 text-white'
                    : 'bg-slate-100 text-slate-500 hover:bg-slate-200'
                }`}
              >
                Capture
              </button>

              <button
                onClick={() => removeAssertion(index)}
                className="text-slate-400 hover:text-rose-500 transition-colors"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>

            {showCapture[index] && (
              <div className="flex items-center gap-2 pl-2">
                <span className="text-xs text-slate-500">Capture as:</span>
                <input
                  type="text"
                  value={assertion.capture_as || ''}
                  onChange={(e) =>
                    updateAssertion(index, { capture_as: e.target.value || undefined })
                  }
                  placeholder="variable_name"
                  className="flex-1 px-2 py-1.5 border border-slate-200 rounded text-sm bg-white font-mono text-slate-700 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
            )}
          </div>
        )
      })}

      {Object.keys(capturedVars).length > 0 && (
        <div className="p-3 bg-indigo-50 border border-indigo-100 rounded-lg">
          <p className="text-xs font-semibold text-indigo-700 mb-2">Available Captured Variables:</p>
          <div className="flex flex-wrap gap-1">
            {Object.entries(capturedVars).map(([name, value]) => (
              <span
                key={name}
                className="px-2 py-0.5 bg-white border border-indigo-200 rounded text-xs font-mono text-indigo-600"
              >
                {`{{${name}}}`} = {String(value)}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
