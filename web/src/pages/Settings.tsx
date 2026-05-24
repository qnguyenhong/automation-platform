import { useState } from 'react'
import { Bell, User, Key } from 'lucide-react'

export default function Settings() {
  const [activeTab, setActiveTab] = useState('notifications')

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Settings</h1>

      <div className="bg-white rounded-lg shadow">
        <div className="border-b border-gray-200">
          <nav className="flex gap-8 px-6">
            <button
              onClick={() => setActiveTab('notifications')}
              className={`py-4 text-sm font-medium border-b-2 ${
                activeTab === 'notifications'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <div className="flex items-center gap-2">
                <Bell className="w-4 h-4" />
                Notifications
              </div>
            </button>
            <button
              onClick={() => setActiveTab('users')}
              className={`py-4 text-sm font-medium border-b-2 ${
                activeTab === 'users'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <div className="flex items-center gap-2">
                <User className="w-4 h-4" />
                Users
              </div>
            </button>
            <button
              onClick={() => setActiveTab('api-keys')}
              className={`py-4 text-sm font-medium border-b-2 ${
                activeTab === 'api-keys'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <div className="flex items-center gap-2">
                <Key className="w-4 h-4" />
                API Keys
              </div>
            </button>
          </nav>
        </div>

        <div className="p-6">
          {activeTab === 'notifications' && (
            <div>
              <h3 className="text-lg font-semibold mb-4">Notification Settings</h3>
              <p className="text-gray-500">Configure Slack, email, or webhook notifications for test failures.</p>
              <div className="mt-4 p-8 bg-gray-50 rounded-lg text-center text-gray-400">
                Notification configuration coming soon
              </div>
            </div>
          )}

          {activeTab === 'users' && (
            <div>
              <h3 className="text-lg font-semibold mb-4">User Management</h3>
              <p className="text-gray-500">Manage user accounts and roles.</p>
              <div className="mt-4 p-8 bg-gray-50 rounded-lg text-center text-gray-400">
                User management coming soon
              </div>
            </div>
          )}

          {activeTab === 'api-keys' && (
            <div>
              <h3 className="text-lg font-semibold mb-4">API Keys</h3>
              <p className="text-gray-500">Manage API keys for CI/CD integration.</p>
              <div className="mt-4 p-8 bg-gray-50 rounded-lg text-center text-gray-400">
                API key management coming soon
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
