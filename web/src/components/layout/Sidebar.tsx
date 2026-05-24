import { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useProjects } from '../../hooks/useProjects'
import { Project } from '../../api/projects'
import ProjectModal from '../project/ProjectModal'
import {
  LayoutDashboard,
  FolderKanban,
  PlayCircle,
  Server,
  Settings,
  Plus,
  Settings2,
} from 'lucide-react'

const navigation = [
  { name: 'Performance Analytics', href: '/', icon: LayoutDashboard },
  { name: 'Worker Nodes', href: '/workers', icon: Server },
  { name: 'System Preferences', href: '/settings', icon: Settings },
]

export default function Sidebar() {
  const location = useLocation()
  const { data: projects } = useProjects()

  const [isModalOpen, setIsModalOpen] = useState(false)
  const [projectToEdit, setProjectToEdit] = useState<Project | null>(null)

  const handleCreateClick = () => {
    setProjectToEdit(null)
    setIsModalOpen(true)
  }

  const handleEditClick = (e: React.MouseEvent, project: Project) => {
    e.preventDefault()
    e.stopPropagation()
    setProjectToEdit(project)
    setIsModalOpen(true)
  }

  return (
    <aside className="w-64 bg-[#0B0F19] border-r border-slate-800/80 flex flex-col z-20">
      <div className="p-5 border-b border-slate-800/60">
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 bg-gradient-to-tr from-indigo-500 to-violet-500 rounded-md flex items-center justify-center shadow-lg shadow-indigo-500/20">
            <span className="text-white text-xs font-black font-mono">A</span>
          </div>
          <h1 className="text-lg font-bold text-gradient-indigo tracking-tight">
            Antigravity QA
          </h1>
        </div>
      </div>

      <nav className="flex-1 overflow-y-auto p-4 space-y-5">
        <div className="space-y-1">
          {navigation.map((item) => {
            const isActive = location.pathname === item.href
            return (
              <Link
                key={item.name}
                to={item.href}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-semibold transition-all border ${
                  isActive
                    ? 'bg-indigo-500/10 border-indigo-500/30 text-indigo-400'
                    : 'border-transparent text-slate-400 hover:text-slate-200 hover:bg-slate-800/40'
                }`}
              >
                <item.icon className="w-4 h-4 shrink-0" />
                {item.name}
              </Link>
            )
          })}
        </div>

        <div className="pt-2">
          <div className="flex items-center justify-between px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-slate-500 border-b border-slate-800/40 pb-2 mb-2">
            <div className="flex items-center gap-2">
              <FolderKanban className="w-3.5 h-3.5 text-slate-500" />
              Workspace Projects
            </div>
            <button
              onClick={handleCreateClick}
              className="p-1 hover:bg-slate-800/60 rounded text-slate-500 hover:text-slate-200 transition-colors"
              title="Initialize Project"
            >
              <Plus className="w-3.5 h-3.5" />
            </button>
          </div>

          <div className="space-y-1 mt-2 max-h-72 overflow-y-auto pr-1">
            {projects?.map((project) => {
              const isActive = location.pathname.startsWith(`/projects/${project.id}`)
              return (
                <div key={project.id} className="relative group flex items-center">
                  <Link
                    to={`/projects/${project.id}/suites`}
                    className={`w-full flex items-center gap-3 pl-3 pr-9 py-2 rounded-lg text-sm font-medium transition-all border ${
                      isActive
                        ? 'bg-indigo-500/5 border-indigo-500/20 text-indigo-400'
                        : 'border-transparent text-slate-400 hover:text-slate-200 hover:bg-slate-800/20'
                    }`}
                  >
                    <PlayCircle className={`w-4 h-4 shrink-0 ${isActive ? 'text-indigo-400' : 'text-slate-500'}`} />
                    <span className="truncate">{project.name}</span>
                  </Link>
                  <button
                    onClick={(e) => handleEditClick(e, project)}
                    className="absolute right-2 hidden group-hover:flex items-center p-1 bg-slate-900 border border-slate-700 text-slate-400 hover:text-indigo-400 rounded-md transition shadow-md"
                    title="Configure Project Settings"
                  >
                    <Settings2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              )
            })}
          </div>
        </div>
      </nav>

      <ProjectModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        projectToEdit={projectToEdit}
      />
    </aside>
  )
}
