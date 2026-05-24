import { useState, useRef, useEffect, useCallback } from 'react'
import { BUILT_IN_VARIABLES, type VariableDefinition } from '@/types/openapi'

interface VariableInputProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  className?: string
  capturedVars?: Record<string, unknown>
}

export function VariableInput({
  value,
  onChange,
  placeholder,
  className = '',
  capturedVars = {},
}: VariableInputProps) {
  const [showAutocomplete, setShowAutocomplete] = useState(false)
  const [filter, setFilter] = useState('')
  const [selectedIndex, setSelectedIndex] = useState(0)
  const inputRef = useRef<HTMLInputElement>(null)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const allVariables = [
    ...BUILT_IN_VARIABLES,
    ...Object.keys(capturedVars).map((name) => ({
      name,
      description: `Captured from prior test`,
      example: String(capturedVars[name]),
      category: 'captured' as const,
    })),
  ]

  const filtered = allVariables.filter(
    (v) =>
      v.name.toLowerCase().includes(filter.toLowerCase()) ||
      v.description.toLowerCase().includes(filter.toLowerCase())
  )

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value
      onChange(newValue)

      const cursorPos = e.target.selectionStart || newValue.length
      const beforeCursor = newValue.slice(0, cursorPos)
      const openBrace = beforeCursor.lastIndexOf('{{')

      if (openBrace !== -1 && !beforeCursor.slice(openBrace).includes('}}')) {
        const searchTerm = beforeCursor.slice(openBrace + 2)
        setFilter(searchTerm)
        setShowAutocomplete(true)
        setSelectedIndex(0)
      } else {
        setShowAutocomplete(false)
      }
    },
    [onChange]
  )

  const insertVariable = useCallback(
    (variable: VariableDefinition) => {
      const cursorPos = inputRef.current?.selectionStart || value.length
      const beforeCursor = value.slice(0, cursorPos)
      const openBrace = beforeCursor.lastIndexOf('{{')

      if (openBrace !== -1) {
        const before = value.slice(0, openBrace)
        const afterCursor = value.slice(cursorPos)
        const newValue = `${before}{{${variable.name}}}${afterCursor}`
        onChange(newValue)
      } else {
        onChange(`${value}{{${variable.name}}}`)
      }

      setShowAutocomplete(false)
      inputRef.current?.focus()
    },
    [value, onChange]
  )

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(e.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(e.target as Node)
      ) {
        setShowAutocomplete(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showAutocomplete) return

    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setSelectedIndex((i) => Math.min(i + 1, filtered.length - 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setSelectedIndex((i) => Math.max(i - 1, 0))
    } else if (e.key === 'Enter' && filtered[selectedIndex]) {
      e.preventDefault()
      insertVariable(filtered[selectedIndex])
    } else if (e.key === 'Escape') {
      setShowAutocomplete(false)
    }
  }

  const highlightVariables = (text: string) => {
    const parts = text.split(/(\{\{[^}]+\}\})/g)
    return parts.map((part, i) => {
      if (part.startsWith('{{') && part.endsWith('}}')) {
        return (
          <span key={i} className="text-indigo-600 font-medium bg-indigo-50 px-0.5 rounded">
            {part}
          </span>
        )
      }
      return part
    })
  }

  return (
    <div className="relative">
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          value={value}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          onFocus={() => {
            const cursorPos = inputRef.current?.selectionStart || value.length
            const beforeCursor = value.slice(0, cursorPos)
            if (beforeCursor.includes('{{') && !beforeCursor.slice(beforeCursor.lastIndexOf('{{')).includes('}}')) {
              setShowAutocomplete(true)
            }
          }}
          placeholder={placeholder || 'Type {{ to see variables...'}
          className={`w-full px-3 py-1.5 border border-slate-200 rounded-lg bg-white text-slate-800 focus:outline-none focus:ring-1 focus:ring-indigo-500 font-mono text-sm placeholder-slate-400 ${className}`}
        />
        {value.includes('{{') && (
          <div className="absolute inset-0 px-3 py-1.5 pointer-events-none font-mono text-sm whitespace-pre">
            {highlightVariables(value)}
          </div>
        )}
      </div>

      {showAutocomplete && filtered.length > 0 && (
        <div
          ref={dropdownRef}
          className="absolute z-50 w-full mt-1 bg-white border border-slate-200 rounded-lg shadow-lg max-h-64 overflow-y-auto"
        >
          {filtered.map((variable, index) => (
            <button
              key={variable.name}
              onClick={() => insertVariable(variable)}
              className={`w-full text-left px-3 py-2 flex items-center justify-between transition-colors ${
                index === selectedIndex ? 'bg-indigo-50' : 'hover:bg-slate-50'
              }`}
            >
              <div>
                <span className="font-mono text-sm font-medium text-slate-800">
                  {`{{${variable.name}}}`}
                </span>
                <span className="text-xs text-slate-400 ml-2">
                  {variable.description}
                </span>
              </div>
              <span className="text-xs text-slate-400 font-mono">
                {variable.example}
              </span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
