import { useState } from 'react'
import { ChevronDown, ChevronRight } from 'lucide-react'
import { VariableInput } from './VariableInput'
import { VariablePreview } from './VariablePreview'
import { AssertionBuilder } from './AssertionBuilder'
import type { ParsedEndpoint, ImportEndpointConfig, AssertionConfig } from '@/types/openapi'

interface TestDataFormProps {
  endpoints: ParsedEndpoint[]
  configs: Record<string, ImportEndpointConfig>
  onConfigChange: (key: string, config: ImportEndpointConfig) => void
  capturedVars: Record<string, unknown>
  onCapture: (name: string, value: unknown) => void
}

const METHOD_COLORS: Record<string, string> = {
  GET: 'bg-sky-50 text-sky-600 border border-sky-200',
  POST: 'bg-emerald-50 text-emerald-600 border border-emerald-200',
  PUT: 'bg-amber-50 text-amber-600 border border-amber-200',
  DELETE: 'bg-rose-50 text-rose-600 border border-rose-200',
}

export function TestDataForm({
  endpoints,
  configs,
  onConfigChange,
  capturedVars,
  onCapture,
}: TestDataFormProps) {
  const [collapsed, setCollapsed] = useState<Set<string>>(new Set())
  const [activeTab, setActiveTab] = useState<Record<string, 'params' | 'headers' | 'body' | 'assertions'>>({})

  const endpointKey = (ep: ParsedEndpoint) => `${ep.method} ${ep.path}`

  const toggleCollapse = (key: string) => {
    setCollapsed((prev) => {
      const next = new Set(prev)
      if (next.has(key)) next.delete(key)
      else next.add(key)
      return next
    })
  }

  const getConfig = (key: string, ep: ParsedEndpoint): ImportEndpointConfig => {
    if (configs[key]) return configs[key]

    const defaultConfig: ImportEndpointConfig = {
      method: ep.method,
      path: ep.path,
      headers: { 'Content-Type': 'application/json' },
      params: {},
      body: ep.request_body?.content?.['application/json']?.example || {},
      assertions: [{ type: 'status', expected: 200 }],
    }

    for (const param of ep.parameters || []) {
      if (param.in === 'path') {
        defaultConfig.path = defaultConfig.path.replace(
          `{${param.name}}`,
          `{{${param.name}}}`
        )
      } else if (param.in === 'query') {
        defaultConfig.params[param.name] = param.example
          ? String(param.example)
          : `{{${param.name}}}`
      } else if (param.in === 'header') {
        defaultConfig.headers[param.name] = param.example
          ? String(param.example)
          : `{{${param.name}}}`
      }
    }

    return defaultConfig
  }

  const getTab = (key: string) => activeTab[key] || 'body'

  return (
    <div className="space-y-3">
      {endpoints.map((ep) => {
        const key = endpointKey(ep)
        const config = getConfig(key, ep)
        const isCollapsed = collapsed.has(key)
        const tab = getTab(key)

        return (
          <div key={key} className="border border-slate-200 rounded-lg overflow-hidden">
            <button
              onClick={() => toggleCollapse(key)}
              className="w-full flex items-center gap-3 px-4 py-3 bg-slate-50 hover:bg-slate-100 transition-colors text-left"
            >
              {isCollapsed ? (
                <ChevronRight className="h-4 w-4 text-slate-400" />
              ) : (
                <ChevronDown className="h-4 w-4 text-slate-400" />
              )}
              <span
                className={`px-2 py-0.5 rounded text-[10px] font-mono font-bold ${
                  METHOD_COLORS[ep.method] || 'bg-slate-100 text-slate-600 border border-slate-200'
                }`}
              >
                {ep.method}
              </span>
              <span className="font-mono text-sm text-slate-700">{ep.path}</span>
              {ep.summary && (
                <span className="text-xs text-slate-400 ml-auto">
                  {ep.summary}
                </span>
              )}
            </button>

            {!isCollapsed && (
              <div className="p-4 space-y-4">
                <div className="flex gap-1 border-b border-slate-200">
                  {(['params', 'headers', 'body', 'assertions'] as const).map((t) => (
                    <button
                      key={t}
                      onClick={() =>
                        setActiveTab((prev) => ({ ...prev, [key]: t }))
                      }
                      className={`px-3 py-2 text-sm font-semibold border-b-2 transition-colors ${
                        tab === t
                          ? 'border-indigo-500 text-indigo-600'
                          : 'border-transparent text-slate-500 hover:text-slate-700'
                      }`}
                    >
                      {t.charAt(0).toUpperCase() + t.slice(1)}
                    </button>
                  ))}
                </div>

                {tab === 'params' && (
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <label className="text-sm font-semibold text-slate-700">Query Parameters</label>
                      <button
                        onClick={() =>
                          onConfigChange(key, {
                            ...config,
                            params: { ...config.params, '': '' },
                          })
                        }
                        className="text-xs text-indigo-600 hover:text-indigo-700 font-semibold"
                      >
                        + Add
                      </button>
                    </div>
                    {Object.entries(config.params).map(([name, val], i) => (
                      <div key={i} className="flex items-center gap-2">
                        <input
                          type="text"
                          value={name}
                          onChange={(e) => {
                            const newParams = { ...config.params }
                            delete newParams[name]
                            newParams[e.target.value] = val
                            onConfigChange(key, { ...config, params: newParams })
                          }}
                          placeholder="Name"
                          className="w-40 px-2 py-1.5 border border-slate-200 rounded text-sm bg-white text-slate-700 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                        />
                        <div className="flex-1">
                          <VariableInput
                            value={val}
                            onChange={(v) =>
                              onConfigChange(key, {
                                ...config,
                                params: { ...config.params, [name]: v },
                              })
                            }
                            capturedVars={capturedVars}
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                {tab === 'headers' && (
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <label className="text-sm font-semibold text-slate-700">Headers</label>
                      <button
                        onClick={() =>
                          onConfigChange(key, {
                            ...config,
                            headers: { ...config.headers, '': '' },
                          })
                        }
                        className="text-xs text-indigo-600 hover:text-indigo-700 font-semibold"
                      >
                        + Add
                      </button>
                    </div>
                    {Object.entries(config.headers).map(([name, val], i) => (
                      <div key={i} className="flex items-center gap-2">
                        <input
                          type="text"
                          value={name}
                          onChange={(e) => {
                            const newHeaders = { ...config.headers }
                            delete newHeaders[name]
                            newHeaders[e.target.value] = val
                            onConfigChange(key, { ...config, headers: newHeaders })
                          }}
                          placeholder="Name"
                          className="w-40 px-2 py-1.5 border border-slate-200 rounded text-sm bg-white text-slate-700 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                        />
                        <div className="flex-1">
                          <VariableInput
                            value={val}
                            onChange={(v) =>
                              onConfigChange(key, {
                                ...config,
                                headers: { ...config.headers, [name]: v },
                              })
                            }
                            capturedVars={capturedVars}
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                {tab === 'body' && (
                  <div className="space-y-2">
                    <label className="text-sm font-semibold text-slate-700">Request Body (JSON)</label>
                    <textarea
                      value={
                        typeof config.body === 'string'
                          ? config.body
                          : JSON.stringify(config.body, null, 2)
                      }
                      onChange={(e) => {
                        try {
                          const parsed = JSON.parse(e.target.value)
                          onConfigChange(key, { ...config, body: parsed })
                        } catch {
                          onConfigChange(key, { ...config, body: e.target.value })
                        }
                      }}
                      rows={8}
                      className="w-full px-3 py-2 border border-slate-200 rounded-lg bg-white text-slate-800 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                    <p className="text-xs text-slate-400">
                      Use {'{{variable_name}}'} for dynamic values.
                    </p>
                    <VariablePreview
                      value={
                        typeof config.body === 'string'
                          ? config.body
                          : JSON.stringify(config.body)
                      }
                      capturedVars={capturedVars}
                    />
                  </div>
                )}

                {tab === 'assertions' && (
                  <AssertionBuilder
                    assertions={config.assertions}
                    onChange={(assertions: AssertionConfig[]) =>
                      onConfigChange(key, { ...config, assertions })
                    }
                    capturedVars={capturedVars}
                    onCapture={onCapture}
                  />
                )}
              </div>
            )}
          </div>
        )
      })}
    </div>
  )
}
