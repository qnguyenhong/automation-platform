import { useState, useMemo } from 'react'
import { Search, ChevronDown, ChevronRight } from 'lucide-react'
import type { ParsedEndpoint } from '@/types/openapi'

interface EndpointSelectorProps {
  endpoints: ParsedEndpoint[]
  selected: Set<string>
  onToggle: (key: string) => void
  onSelectAll: () => void
  onDeselectAll: () => void
}

const METHOD_COLORS: Record<string, string> = {
  GET: 'bg-sky-50 text-sky-600 border border-sky-200',
  POST: 'bg-emerald-50 text-emerald-600 border border-emerald-200',
  PUT: 'bg-amber-50 text-amber-600 border border-amber-200',
  DELETE: 'bg-rose-50 text-rose-600 border border-rose-200',
  PATCH: 'bg-violet-50 text-violet-600 border border-violet-200',
}

export function EndpointSelector({
  endpoints,
  selected,
  onToggle,
  onSelectAll,
  onDeselectAll,
}: EndpointSelectorProps) {
  const [search, setSearch] = useState('')
  const [collapsedTags, setCollapsedTags] = useState<Set<string>>(new Set())

  const filtered = useMemo(() => {
    if (!search) return endpoints
    const q = search.toLowerCase()
    return endpoints.filter(
      (ep) =>
        ep.path.toLowerCase().includes(q) ||
        ep.summary.toLowerCase().includes(q) ||
        ep.method.toLowerCase().includes(q)
    )
  }, [endpoints, search])

  const grouped = useMemo(() => {
    const groups: Record<string, ParsedEndpoint[]> = {}
    for (const ep of filtered) {
      const tag = (ep.tags && ep.tags[0]) || 'default'
      if (!groups[tag]) groups[tag] = []
      groups[tag].push(ep)
    }
    return groups
  }, [filtered])

  const toggleTag = (tag: string) => {
    setCollapsedTags((prev) => {
      const next = new Set(prev)
      if (next.has(tag)) next.delete(tag)
      else next.add(tag)
      return next
    })
  }

  const endpointKey = (ep: ParsedEndpoint) => `${ep.method} ${ep.path}`

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search endpoints..."
            className="w-full pl-10 pr-4 py-2 border border-slate-200 rounded-lg bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500 placeholder-slate-400 text-sm"
          />
        </div>
        <button
          onClick={onSelectAll}
          className="px-3 py-2 text-xs font-semibold border border-slate-200 rounded-lg hover:bg-slate-50 text-slate-600 transition-colors"
        >
          Select All
        </button>
        <button
          onClick={onDeselectAll}
          className="px-3 py-2 text-xs font-semibold border border-slate-200 rounded-lg hover:bg-slate-50 text-slate-600 transition-colors"
        >
          Deselect All
        </button>
        <span className="text-xs font-semibold text-indigo-600 shrink-0">
          {selected.size} selected
        </span>
      </div>

      <div className="border border-slate-200 rounded-lg divide-y divide-slate-100 overflow-hidden">
        {Object.entries(grouped).map(([tag, eps]) => (
          <div key={tag}>
            <button
              onClick={() => toggleTag(tag)}
              className="w-full flex items-center gap-2 px-4 py-2.5 bg-slate-50 hover:bg-slate-100 transition-colors text-left"
            >
              {collapsedTags.has(tag) ? (
                <ChevronRight className="h-4 w-4 text-slate-400" />
              ) : (
                <ChevronDown className="h-4 w-4 text-slate-400" />
              )}
              <span className="font-semibold text-sm text-slate-700 capitalize">{tag}</span>
              <span className="text-xs text-slate-400">
                ({eps.length})
              </span>
            </button>
            {!collapsedTags.has(tag) && (
              <div className="divide-y divide-slate-50">
                {eps.map((ep) => {
                  const key = endpointKey(ep)
                  return (
                    <label
                      key={key}
                      className="flex items-center gap-3 px-4 py-3 hover:bg-slate-50 cursor-pointer transition-colors"
                    >
                      <input
                        type="checkbox"
                        checked={selected.has(key)}
                        onChange={() => onToggle(key)}
                        className="h-4 w-4 rounded border-slate-300 accent-indigo-600"
                      />
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
                    </label>
                  )
                })}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
