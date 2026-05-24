import { useAuthStore } from '../../store/auth'
import { useWebSocket } from '../../hooks/useWebSocket'
import { Bell, User, LogOut } from 'lucide-react'

export default function Header() {
  const { user, logout } = useAuthStore()

  useWebSocket((data) => {
    console.log('WebSocket message:', data)
  })

  return (
    <header className="bg-white/80 backdrop-blur-md border-b border-slate-200/80 px-6 py-3.5 z-10">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <h2 className="text-sm font-semibold tracking-wider uppercase text-slate-500 font-mono">
            Control Plane Active
          </h2>
          <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
        </div>

        <div className="flex items-center gap-4">
          <button className="p-2 text-slate-500 hover:text-slate-800 hover:bg-slate-100 border border-transparent hover:border-slate-200 rounded-lg transition-all">
            <Bell className="w-4 h-4" />
          </button>

          <div className="h-6 w-[1px] bg-slate-200" />

          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-violet-500 rounded-full flex items-center justify-center shadow-md shadow-indigo-500/10">
              <User className="w-4 h-4 text-white" />
            </div>
            <div className="text-sm text-left">
              <p className="font-semibold text-slate-800">{user?.name || 'Developer'}</p>
              <p className="text-[10px] text-slate-500 font-mono">{user?.email || 'admin@example.com'}</p>
            </div>
            <button
              onClick={logout}
              className="flex items-center gap-1.5 ml-2 text-xs font-semibold text-slate-600 hover:text-rose-600 px-2.5 py-1.5 rounded-lg border border-slate-200 bg-slate-50 hover:bg-rose-50 hover:border-rose-200 transition-all active:scale-95"
            >
              <LogOut className="w-3.5 h-3.5" />
              Sign Out
            </button>
          </div>
        </div>
      </div>
    </header>
  )
}
