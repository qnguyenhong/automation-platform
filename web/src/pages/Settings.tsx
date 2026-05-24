import { useState } from 'react'
import { Bell, User, Key } from 'lucide-react'

export default function Settings() {
  const [activeTab, setActiveTab] = useState('notifications')

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">System Preferences</h1>
        <p className="text-slate-400 text-xs mt-1">Configure global notification pipelines, user roles, and access credentials.</p>
      </div>

      <div className="glass-panel rounded-xl border border-slate-850 overflow-hidden">
        <div className="border-b border-slate-850 bg-slate-900/10">
          <nav className="flex gap-8 px-6">
            <button
              onClick={() => setActiveTab('notifications')}
              className={`py-4 text-xs font-bold uppercase tracking-wider border-b-2 flex items-center gap-2 transition ${
                activeTab === 'notifications'
                  ? 'border-indigo-500 text-indigo-400'
                  : 'border-transparent text-slate-500 hover:text-slate-300'
              }`}
            >
              <Bell className="w-4 h-4" />
              Notifications
            </button>
            <button
              onClick={() => setActiveTab('users')}
              className={`py-4 text-xs font-bold uppercase tracking-wider border-b-2 flex items-center gap-2 transition ${
                activeTab === 'users'
                  ? 'border-indigo-500 text-indigo-400'
                  : 'border-transparent text-slate-500 hover:text-slate-300'
              }`}
            >
              <User className="w-4 h-4" />
              Users
            </button>
            <button
              onClick={() => setActiveTab('api-keys')}
              className={`py-4 text-xs font-bold uppercase tracking-wider border-b-2 flex items-center gap-2 transition ${
                activeTab === 'api-keys'
                  ? 'border-indigo-500 text-indigo-400'
                  : 'border-transparent text-slate-500 hover:text-slate-300'
              }`}
            >
              <Key className="w-4 h-4" />
              API Keys
            </button>
          </nav>
        </div>

        <div className="p-6 bg-slate-900/10">
          {activeTab === 'notifications' && (
            <div className="space-y-3 animate-in fade-in duration-200">
              <h3 className="text-base font-bold text-slate-200">Notification Settings</h3>
              <p className="text-xs text-slate-450">Configure Slack hooks, SMTP email notifications, or pager alerts for execution failures.</p>
              <div className="p-12 bg-slate-950/60 border border-dashed border-slate-800 rounded-xl text-center text-slate-550 text-xs font-semibold">
                Slack and Webhook pipeline configs coming soon
              </div>
            </div>
          )}

          {activeTab === 'users' && (
            <div className="space-y-3 animate-in fade-in duration-200">
              <h3 className="text-base font-bold text-slate-200">User Management</h3>
              <p className="text-xs text-slate-450">Provision team developer profiles, access roles, and permission levels.</p>
              <div className="p-12 bg-slate-950/60 border border-dashed border-slate-800 rounded-xl text-center text-slate-550 text-xs font-semibold">
                Developer profile provisioning coming soon
              </div>
            </div>
          )}

          {activeTab === 'api-keys' && (
            <div className="space-y-3 animate-in fade-in duration-200">
              <h3 className="text-base font-bold text-slate-200">API Access Tokens</h3>
              <p className="text-xs text-slate-450">Generate secure API keys to integrate validation suite triggering directly with your CI/CD pipelines.</p>
              <div className="p-12 bg-slate-950/60 border border-dashed border-slate-800 rounded-xl text-center text-slate-550 text-xs font-semibold">
                Access token generation coming soon
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
