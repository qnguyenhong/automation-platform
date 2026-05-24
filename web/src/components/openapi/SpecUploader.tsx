import { useState, useCallback } from 'react'
import { Upload, FileText, Link, Loader2 } from 'lucide-react'

interface SpecUploaderProps {
  onParsed: (content: string) => void
  isLoading: boolean
}

export function SpecUploader({ onParsed, isLoading }: SpecUploaderProps) {
  const [mode, setMode] = useState<'file' | 'url'>('file')
  const [url, setUrl] = useState('')
  const [dragActive, setDragActive] = useState(false)

  const handleFile = useCallback((file: File) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const content = e.target?.result as string
      onParsed(content)
    }
    reader.readAsText(file)
  }, [onParsed])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDragActive(false)
    const file = e.dataTransfer.files[0]
    if (file) handleFile(file)
  }, [handleFile])

  const handleFileInput = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) handleFile(file)
  }, [handleFile])

  const handleUrlSubmit = useCallback(async () => {
    if (!url) return
    try {
      const resp = await fetch(url)
      const content = await resp.text()
      onParsed(content)
    } catch {
      alert('Failed to fetch spec from URL')
    }
  }, [url, onParsed])

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <button
          onClick={() => setMode('file')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            mode === 'file'
              ? 'bg-primary text-primary-foreground'
              : 'bg-muted text-muted-foreground hover:bg-muted/80'
          }`}
        >
          <FileText className="inline-block w-4 h-4 mr-2" />
          Upload File
        </button>
        <button
          onClick={() => setMode('url')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            mode === 'url'
              ? 'bg-primary text-primary-foreground'
              : 'bg-muted text-muted-foreground hover:bg-muted/80'
          }`}
        >
          <Link className="inline-block w-4 h-4 mr-2" />
          From URL
        </button>
      </div>

      {mode === 'file' ? (
        <div
          onDragOver={(e) => { e.preventDefault(); setDragActive(true) }}
          onDragLeave={() => setDragActive(false)}
          onDrop={handleDrop}
          className={`border-2 border-dashed rounded-lg p-12 text-center transition-colors ${
            dragActive
              ? 'border-primary bg-primary/5'
              : 'border-muted-foreground/25 hover:border-muted-foreground/50'
          }`}
        >
          <Upload className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
          <p className="text-lg font-medium mb-2">
            Drag and drop your OpenAPI spec
          </p>
          <p className="text-sm text-muted-foreground mb-4">
            Supports .yaml, .yml, and .json files
          </p>
          <label className="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground rounded-lg cursor-pointer hover:bg-primary/90 transition-colors">
            {isLoading ? (
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
            ) : (
              <Upload className="w-4 h-4 mr-2" />
            )}
            {isLoading ? 'Parsing...' : 'Browse Files'}
            <input
              type="file"
              accept=".yaml,.yml,.json"
              onChange={handleFileInput}
              className="hidden"
              disabled={isLoading}
            />
          </label>
        </div>
      ) : (
        <div className="space-y-3">
          <input
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="https://petstore.example.com/openapi.yaml"
            className="w-full px-4 py-3 border rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
          />
          <button
            onClick={handleUrlSubmit}
            disabled={!url || isLoading}
            className="w-full px-4 py-3 bg-primary text-primary-foreground rounded-lg font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
          >
            {isLoading ? (
              <Loader2 className="inline-block w-4 h-4 mr-2 animate-spin" />
            ) : null}
            {isLoading ? 'Parsing...' : 'Parse Spec'}
          </button>
        </div>
      )}
    </div>
  )
}
