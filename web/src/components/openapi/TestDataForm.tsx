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

    // Build default config from endpoint
    const defaultConfig: ImportEndpointConfig = {
      method: ep.method,
      path: ep.path,
      headers: { 'Content-Type': 'application/json' },
      params: {},
      body: ep.request_body?.content?.['application/json']?.example || {},
      assertions: [{ type: 'status', expected: 200 }],
    }

    // Add path parameters
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
          <div key={key} className="border rounded-lg overflow-hidden">
            <button
              onClick={() => toggleCollapse(key)}
              className="w-full flex items-center gap-3 px-4 py-3 bg-muted/50 hover:bg-muted transition-colors"
            >
              {isCollapsed ? (
                <ChevronRight className="h-4 w-4" />
              ) : (
                <ChevronDown className="h-4 w-4" />
              )}
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
              <span className="font-mono text-sm">{ep.path}</span>
              {ep.summary && (
                <span className="text-sm text-muted-foreground ml-auto">
                  {ep.summary}
                </span>
              )}
            </button>

            {!isCollapsed && (
              <div className="p-4 space-y-4">
                <div className="flex gap-1 border-b">
                  {(['params', 'headers', 'body', 'assertions'] as const).map((t) => (
                    <button
                      key={t}
                      onClick={() =>
                        setActiveTab((prev) => ({ ...prev, [key]: t }))
                      }
                      className={`px-3 py-2 text-sm font-medium border-b-2 transition-colors ${
                        tab === t
                          ? 'border-primary text-primary'
                          : 'border-transparent text-muted-foreground hover:text-foreground'
                      }`}
                    >
                      {t.charAt(0).toUpperCase() + t.slice(1)}
                    </button>
                  ))}
                </div>

                {tab === 'params' && (
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <label className="text-sm font-medium">Query Parameters</label>
                      <button
                        onClick={() =>
                          onConfigChange(key, {
                            ...config,
                            params: { ...config.params, '': '' },
                          })
                        }
                        className="text-xs text-primary hover:text-primary/80"
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
                          className="w-40 px-2 py-1 border rounded text-sm bg-background"
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
                      <label className="text-sm font-medium">Headers</label>
                      <button
                        onClick={() =>
                          onConfigChange(key, {
                            ...config,
                            headers: { ...config.headers, '': '' },
                          })
                        }
                        className="text-xs text-primary hover:text-primary/80"
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
                          className="w-40 px-2 py-1 border rounded text-sm bg-background"
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
                    <label className="text-sm font-medium">Request Body (JSON)</label>
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
                      className="w-full px-3 py-2 border rounded-lg bg-background text-foreground font-mono text-sm focus:outline-none focus:ring-2 focus:ring-primary"
                    />
                    <p className="text-xs text-muted-foreground">
                      Use {'{{variable_name}}'} for dynamic values. Type {'{{'} to see available variables.
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
