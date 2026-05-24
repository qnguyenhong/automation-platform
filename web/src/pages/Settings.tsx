import { useState } from 'react'
import { Bell, User, Key } from 'lucide-react'

export default function Settings() {
  const [activeTab, setActiveTab] = useState('notifications')

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-800">System Preferences</h1>
        <p className="text-slate-500 text-xs mt-1">Configure global notification pipelines, user roles, and access credentials.</p>
      </div>

      <div className="glass-panel rounded-xl border border-slate-200 overflow-hidden shadow-sm">
        <div className="border-b border-slate-200 bg-slate-50/50">
          <nav className="flex gap-8 px-6">
            <button
              onClick={() => setActiveTab('notifications')}
              className={`py-4 text-xs font-bold uppercase tracking-wider border-b-2 flex items-center gap-2 transition ${
                activeTab === 'notifications'
                  ? 'border-indigo-600 text-indigo-600'
                  : 'border-transparent text-slate-500 hover:text-slate-850'
              }`}
            >
              <Bell className="w-4 h-4" />
              Notifications
            </button>
            <button
              onClick={() => setActiveTab('users')}
              className={`py-4 text-xs font-bold uppercase tracking-wider border-b-2 flex items-center gap-2 transition ${
                activeTab === 'users'
                  ? 'border-indigo-600 text-indigo-600'
                  : 'border-transparent text-slate-500 hover:text-slate-850'
              }`}
            >
              <User className="w-4 h-4" />
              Users
            </button>
            <button
              onClick={() => setActiveTab('api-keys')}
              className={`py-4 text-xs font-bold uppercase tracking-wider border-b-2 flex items-center gap-2 transition ${
                activeTab === 'api-keys'
                  ? 'border-indigo-600 text-indigo-600'
                  : 'border-transparent text-slate-500 hover:text-slate-850'
              }`}
            >
              <Key className="w-4 h-4" />
              API Keys
            </button>
          </nav>
        </div>

        <div className="p-6 bg-white">
          {activeTab === 'notifications' && (
            <div className="space-y-3 animate-in fade-in duration-200">
              <h3 className="text-base font-bold text-slate-800">Notification Settings</h3>
              <p className="text-xs text-slate-500">Configure Slack hooks, SMTP email notifications, or pager alerts for execution failures.</p>
              <div className="p-12 bg-slate-50 border border-dashed border-slate-200 rounded-xl text-center text-slate-500 text-xs font-semibold shadow-inner">
                Slack and Webhook pipeline configs coming soon
              </div>
            </div>
          )}

          {activeTab === 'users' && (
            <div className="space-y-3 animate-in fade-in duration-200">
              <h3 className="text-base font-bold text-slate-800">User Management</h3>
              <p className="text-xs text-slate-500">Provision team developer profiles, access roles, and permission levels.</p>
              <div className="p-12 bg-slate-50 border border-dashed border-slate-200 rounded-xl text-center text-slate-500 text-xs font-semibold shadow-inner">
                Developer profile provisioning coming soon
              </div>
            </div>
          )}

          {activeTab === 'api-keys' && (
            <div className="space-y-3 animate-in fade-in duration-200">
              <h3 className="text-base font-bold text-slate-800">API Access Tokens</h3>
              <p className="text-xs text-slate-500">Generate secure API keys to integrate validation suite triggering directly with your CI/CD pipelines.</p>
              <div className="p-12 bg-slate-50 border border-dashed border-slate-200 rounded-xl text-center text-slate-500 text-xs font-semibold shadow-inner">
                Access token generation coming soon
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
