import { useAuthStore } from '../../store/auth'
import { useWebSocket } from '../../hooks/useWebSocket'
import { Bell, User } from 'lucide-react'

export default function Header() {
  const { user, logout } = useAuthStore()

  useWebSocket((data) => {
    console.log('WebSocket message:', data)
  })

  return (
    <header className="bg-white border-b border-gray-200 px-6 py-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <h2 className="text-lg font-semibold text-gray-900">Automation Platform</h2>
        </div>

        <div className="flex items-center gap-4">
          <button className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg">
            <Bell className="w-5 h-5" />
          </button>

          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center">
              <User className="w-4 h-4 text-white" />
            </div>
            <div className="text-sm">
              <p className="font-medium text-gray-900">{user?.name || 'User'}</p>
              <p className="text-gray-500">{user?.email || ''}</p>
            </div>
            <button
              onClick={logout}
              className="text-sm text-gray-500 hover:text-gray-700"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </header>
  )
}
