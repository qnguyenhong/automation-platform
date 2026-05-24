import { useState } from 'react'
import { Eye, RefreshCw } from 'lucide-react'

interface VariablePreviewProps {
  value: string
  capturedVars?: Record<string, unknown>
}

export function VariablePreview({ value, capturedVars = {} }: VariablePreviewProps) {
  const [preview, setPreview] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const generatePreview = () => {
    setLoading(true)

    // Client-side preview of common variables
    let result = value
    const replacements: Record<string, string> = {
      '\\{\\{\\$uuid\\}\\}': crypto.randomUUID(),
      '\\{\\{\\$random_email\\}\\}': `user${Math.floor(Math.random() * 1000)}@example.com`,
      '\\{\\{\\$random_string\\(8\\)\\}\\}': Math.random().toString(36).slice(2, 10),
      '\\{\\{\\$random_int\\(1,100\\)\\}\\}': String(Math.floor(Math.random() * 100) + 1),
      '\\{\\{\\$random_bool\\}\\}': String(Math.random() > 0.5),
      '\\{\\{\\$random_name\\}\\}': 'John Smith',
      '\\{\\{\\$timestamp\\}\\}': String(Math.floor(Date.now() / 1000)),
      '\\{\\{\\$timestamp_ms\\}\\}': String(Date.now()),
      '\\{\\{\\$iso_date\\}\\}': new Date().toISOString(),
    }

    for (const [pattern, replacement] of Object.entries(replacements)) {
      result = result.replace(new RegExp(pattern, 'g'), replacement)
    }

    // Replace captured variables
    for (const [name, value] of Object.entries(capturedVars)) {
      result = result.replace(new RegExp(`\\{\\{${name}\\}\\}`, 'g'), String(value))
    }

    setPreview(result)
    setLoading(false)
  }

  if (!value.includes('{{')) return null

  return (
    <div className="mt-2">
      <button
        onClick={generatePreview}
        disabled={loading}
        className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
      >
        {loading ? (
          <RefreshCw className="h-3 w-3 animate-spin" />
        ) : (
          <Eye className="h-3 w-3" />
        )}
        Preview
      </button>
      {preview && (
        <div className="mt-1 p-2 bg-muted rounded text-xs font-mono break-all">
          {preview}
        </div>
      )}
    </div>
  )
}
