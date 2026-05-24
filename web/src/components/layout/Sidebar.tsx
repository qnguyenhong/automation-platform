import { Link, useLocation } from 'react-router-dom'
import { useProjects } from '../../hooks/useProjects'
import {
  LayoutDashboard,
  FolderKanban,
  PlayCircle,
  Server,
  Settings,
  ChevronDown,
} from 'lucide-react'

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Workers', href: '/workers', icon: Server },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export default function Sidebar() {
  const location = useLocation()
  const { data: projects } = useProjects()

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
          <div className="flex items-center gap-2 px-3 py-2 text-sm font-medium text-gray-500">
            <FolderKanban className="w-4 h-4" />
            Projects
          </div>

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
