import { useEffect, useRef, useCallback } from 'react'
import { useAuthStore } from '../store/auth'

type MessageHandler = (data: any) => void

export function useWebSocket(onMessage?: MessageHandler) {
  const wsRef = useRef<WebSocket | null>(null)
  const token = useAuthStore((state) => state.token)

  const connect = useCallback(() => {
    if (!token) return

    const wsUrl = import.meta.env.VITE_WS_URL || `ws://${window.location.host}/api/v1/ws`
    const ws = new WebSocket(`${wsUrl}?token=${token}`)

    ws.onopen = () => {
      console.log('WebSocket connected')
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        onMessage?.(data)
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    ws.onclose = () => {
      console.log('WebSocket disconnected, reconnecting in 3s...')
      setTimeout(connect, 3000)
    }

    wsRef.current = ws

    return () => {
      ws.close()
    }
  }, [token, onMessage])

  useEffect(() => {
    const cleanup = connect()
    return () => {
      cleanup?.()
    }
  }, [connect])

  return wsRef.current
}
