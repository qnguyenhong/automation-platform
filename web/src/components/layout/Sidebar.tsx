import { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useProjects, useCreateProject } from '../../hooks/useProjects'
import {
  LayoutDashboard,
  FolderKanban,
  PlayCircle,
  Server,
  Settings,
  Plus,
  X,
  Check,
} from 'lucide-react'

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Workers', href: '/workers', icon: Server },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export default function Sidebar() {
  const location = useLocation()
  const { data: projects } = useProjects()
  const createProjectMutation = useCreateProject()

  const [isCreating, setIsCreating] = useState(false)
  const [projectName, setProjectName] = useState('')

  return (
    <aside className="w-64 bg-white border-r border-gray-200 flex flex-col">
      <div className="p-4 border-b border-gray-200">
        <h1 className="text-xl font-bold text-gray-900">Automation Platform</h1>
      </div>

      <nav className="flex-1 overflow-y-auto p-4 space-y-1">
        {navigation.map((item) => {
          const isActive = location.pathname === item.href
          return (
            <Link
              key={item.name}
              to={item.href}
              className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                isActive
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              <item.icon className="w-5 h-5" />
              {item.name}
            </Link>
          )
        })}

        <div className="pt-4">
          <div className="flex items-center justify-between px-3 py-2 text-sm font-medium text-gray-500">
            <div className="flex items-center gap-2">
              <FolderKanban className="w-4 h-4" />
              Projects
            </div>
            {!isCreating && (
              <button
                onClick={() => setIsCreating(true)}
                className="p-1 hover:bg-gray-100 rounded text-gray-500 hover:text-gray-900 transition-colors"
                title="Create Project"
              >
                <Plus className="w-4 h-4" />
              </button>
            )}
          </div>

          {isCreating && (
            <form
              onSubmit={async (e) => {
                e.preventDefault()
                if (!projectName.trim()) return
                try {
                  await createProjectMutation.mutateAsync({ name: projectName.trim() })
                  setProjectName('')
                  setIsCreating(false)
                } catch {
                  alert('Failed to create project')
                }
              }}
              className="px-3 py-2 flex items-center gap-1"
            >
              <input
                type="text"
                value={projectName}
                onChange={(e) => setProjectName(e.target.value)}
                placeholder="Project name..."
                className="flex-1 min-w-0 px-2 py-1 text-xs border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                autoFocus
                disabled={createProjectMutation.isPending}
              />
              <button
                type="submit"
                disabled={createProjectMutation.isPending || !projectName.trim()}
                className="p-1 hover:bg-green-50 text-green-600 rounded hover:text-green-800 disabled:opacity-50"
              >
                <Check className="w-4 h-4" />
              </button>
              <button
                type="button"
                onClick={() => {
                  setIsCreating(false)
                  setProjectName('')
                }}
                disabled={createProjectMutation.isPending}
                className="p-1 hover:bg-red-50 text-red-600 rounded hover:text-red-800 disabled:opacity-50"
              >
                <X className="w-4 h-4" />
              </button>
            </form>
          )}

          {projects?.map((project) => (
            <Link
              key={project.id}
              to={`/projects/${project.id}/suites`}
              className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                location.pathname.startsWith(`/projects/${project.id}`)
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              <PlayCircle className="w-4 h-4" />
              {project.name}
            </Link>
          ))}
        </div>
      </nav>
    </aside>
  )
}
