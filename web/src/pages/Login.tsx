import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/auth'
import client from '../api/client'
import { AlertCircle, Lock, Mail } from 'lucide-react'

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await client.post('/auth/login', { email, password })
      const { token, user } = response.data.data
      setAuth(token, user)
      navigate('/')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Authentication failed. Please verify credentials.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#070A13] relative overflow-hidden px-4">
      {/* Decorative Blur Orbs */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-indigo-500/10 rounded-full blur-[100px] animate-blob" />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-fuchsia-500/10 rounded-full blur-[100px] animate-blob animation-delay-2000" />

      {/* Auth Card */}
      <div className="glass-panel w-full max-w-md p-8 rounded-2xl shadow-2xl relative z-10 border border-slate-800/80 animate-in fade-in zoom-in-95 duration-300">
        <div className="text-center mb-8">
          <div className="w-12 h-12 bg-gradient-to-tr from-indigo-500 to-violet-500 rounded-xl flex items-center justify-center shadow-lg shadow-indigo-500/20 mx-auto mb-4">
            <span className="text-white text-lg font-black font-mono">A</span>
          </div>
          <h1 className="text-2xl font-extrabold text-white tracking-tight">
            Antigravity QA Hub
          </h1>
          <p className="text-xs text-slate-400 mt-1.5 leading-relaxed">
            Performance Load Simulation & API Automated Verification Engine
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          {error && (
            <div className="bg-rose-500/10 border border-rose-500/20 text-rose-400 p-4 rounded-xl text-xs font-semibold flex gap-2.5 items-start">
              <AlertCircle className="w-4 h-4 shrink-0 text-rose-400 mt-0.5" />
              <span>{error}</span>
            </div>
          )}

          <div className="space-y-1.5">
            <label className="block text-xs font-bold uppercase tracking-wider text-slate-400">
              Email Address
            </label>
            <div className="relative">
              <span className="absolute left-3.5 top-3.5 text-slate-500">
                <Mail className="w-4 h-4" />
              </span>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="developer@example.com"
                className="w-full pl-10 pr-4 py-3 bg-[#0B0F19]/50 border border-slate-850 rounded-xl text-slate-100 placeholder-slate-600 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm transition-all"
                required
                disabled={loading}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="block text-xs font-bold uppercase tracking-wider text-slate-400">
              Security Key / Password
            </label>
            <div className="relative">
              <span className="absolute left-3.5 top-3.5 text-slate-500">
                <Lock className="w-4 h-4" />
              </span>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className="w-full pl-10 pr-4 py-3 bg-[#0B0F19]/50 border border-slate-850 rounded-xl text-slate-100 placeholder-slate-600 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm transition-all"
                required
                disabled={loading}
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full mt-2 bg-gradient-to-r from-indigo-500 to-violet-600 hover:from-indigo-650 hover:to-violet-750 text-white font-bold py-3.5 rounded-xl shadow-lg shadow-indigo-500/10 active:scale-[0.98] transition-all disabled:opacity-50 text-sm"
          >
            {loading ? 'Authenticating Profile...' : 'Sign In to Workspace'}
          </button>
        </form>
      </div>
    </div>
  )
}
