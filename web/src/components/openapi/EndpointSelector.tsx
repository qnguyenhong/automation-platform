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
  GET: 'bg-green-100 text-green-800',
  POST: 'bg-blue-100 text-blue-800',
  PUT: 'bg-yellow-100 text-yellow-800',
  DELETE: 'bg-red-100 text-red-800',
  PATCH: 'bg-purple-100 text-purple-800',
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
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search endpoints..."
            className="w-full pl-10 pr-4 py-2 border rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>
        <button
          onClick={onSelectAll}
          className="px-3 py-2 text-sm border rounded-lg hover:bg-muted transition-colors"
        >
          Select All
        </button>
        <button
          onClick={onDeselectAll}
          className="px-3 py-2 text-sm border rounded-lg hover:bg-muted transition-colors"
        >
          Deselect All
        </button>
        <span className="text-sm text-muted-foreground">
          {selected.size} selected
        </span>
      </div>

      <div className="border rounded-lg divide-y">
        {Object.entries(grouped).map(([tag, eps]) => (
          <div key={tag}>
            <button
              onClick={() => toggleTag(tag)}
              className="w-full flex items-center gap-2 px-4 py-2 bg-muted/50 hover:bg-muted transition-colors"
            >
              {collapsedTags.has(tag) ? (
                <ChevronRight className="h-4 w-4" />
              ) : (
                <ChevronDown className="h-4 w-4" />
              )}
              <span className="font-medium capitalize">{tag}</span>
              <span className="text-sm text-muted-foreground">
                ({eps.length})
              </span>
            </button>
            {!collapsedTags.has(tag) && (
              <div className="divide-y">
                {eps.map((ep) => {
                  const key = endpointKey(ep)
                  return (
                    <label
                      key={key}
                      className="flex items-center gap-3 px-4 py-3 hover:bg-muted/30 cursor-pointer transition-colors"
                    >
                      <input
                        type="checkbox"
                        checked={selected.has(key)}
                        onChange={() => onToggle(key)}
                        className="h-4 w-4 rounded border-gray-300"
                      />
                      <span
                        className={`px-2 py-0.5 rounded text-xs font-mono font-bold ${
                          METHOD_COLORS[ep.method] || 'bg-gray-100 text-gray-800'
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
