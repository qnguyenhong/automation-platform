import { useParams, Link, useNavigate } from 'react-router-dom'
import { useSuites, useCreateSuite } from '../hooks/useSuites'
import { useProject } from '../hooks/useProjects'
import { useState, useMemo } from 'react'
import { Plus, Play, Settings, Filter, HelpCircle, Layers, Cpu, Zap, Activity, Settings2, Folder } from 'lucide-react'
import { useMutation } from '@tanstack/react-query'
import { runsApi } from '../api/runs'
import ProjectModal from '../components/project/ProjectModal'

type SuiteType = 'all' | 'api' | 'e2e' | 'load' | 'unit';

export default function TestSuites() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const { data: project } = useProject(projectId!)
  const { data: suites, isLoading } = useSuites(projectId!)
  const createSuite = useCreateSuite(projectId!)
  
  const [showCreate, setShowCreate] = useState(false)
  const [isProjectModalOpen, setIsProjectModalOpen] = useState(false)
  const [filterType, setFilterType] = useState<SuiteType>('all')
  const [newSuite, setNewSuite] = useState({
    name: '',
    description: '',
    test_type: 'load' as const, // Default to load testing since it's the core focus
  })

  const triggerRun = useMutation({
    mutationFn: (suiteId: string) => runsApi.trigger(projectId!, suiteId),
    onSuccess: (run) => {
      if (run && run.id) {
        navigate(`/runs/${run.id}`)
      }
    },
  })

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    await createSuite.mutateAsync(newSuite)
    setShowCreate(false)
    setNewSuite({ name: '', description: '', test_type: 'load' })
  }

  // Filtered suites
  const filteredSuites = useMemo(() => {
    if (!suites) return []
    if (filterType === 'all') return suites
    return suites.filter((s) => s.test_type === filterType)
  }, [suites, filterType])

  const getSuiteTypeBadge = (type: string) => {
    switch (type) {
      case 'load':
        return (
          <span className="flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-semibold bg-rose-50 text-rose-700 border border-rose-200">
            <Activity className="w-3.5 h-3.5 animate-pulse" />
            Performance Load Simulation
          </span>
        )
      case 'api':
        return (
          <span className="flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-semibold bg-sky-50 text-sky-700 border border-sky-200">
            <Zap className="w-3.5 h-3.5" />
            REST API Assertion
          </span>
        )
      case 'e2e':
        return (
          <span className="flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-semibold bg-indigo-50 text-indigo-700 border border-indigo-200">
            <Layers className="w-3.5 h-3.5" />
            User Journey E2E Flow
          </span>
        )
      default:
        return (
          <span className="flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-semibold bg-slate-100 text-slate-700 border border-slate-200">
            <Cpu className="w-3.5 h-3.5" />
            Isolated Unit Check
          </span>
        )
    }
  }

  const getFilterLabel = (type: SuiteType) => {
    switch (type) {
      case 'all': return 'All Categories'
      case 'load': return 'Performance Load'
      case 'api': return 'API Assertion'
      case 'e2e': return 'E2E Flow'
      case 'unit': return 'Unit Check'
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="w-8 h-8 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin"></div>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Header and Create Button */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-200 pb-5">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-3xl font-extrabold text-slate-800 tracking-tight">{project?.name || 'Project'}</h1>
            <button
              onClick={() => setIsProjectModalOpen(true)}
              className="p-1.5 hover:bg-slate-100 text-slate-500 hover:text-indigo-600 rounded-lg transition border border-transparent hover:border-slate-200"
              title="Configure Project Settings"
            >
              <Settings2 className="w-5 h-5" />
            </button>
          </div>
          <p className="text-slate-550 mt-1.5 text-sm max-w-2xl leading-relaxed font-medium">
            {project?.description || 'Configure parameters, execution criteria, and simulated workloads for validation runs.'}
          </p>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          className="flex items-center justify-center gap-2 px-4 py-2.5 bg-indigo-600 text-white font-semibold rounded-lg hover:bg-indigo-700 shadow-md shadow-indigo-500/10 transition-all active:scale-[0.98] text-sm shrink-0"
        >
          <Plus className="w-4 h-4" />
          Define Validation Suite
        </button>
      </div>

      {/* Filter Tabs */}
      <div className="flex items-center justify-between flex-wrap gap-3">
        <div className="flex items-center gap-1 bg-white p-1 border border-slate-200 rounded-xl shadow-sm">
          {(['all', 'load', 'api', 'e2e', 'unit'] as const).map((type) => (
            <button
              key={type}
              onClick={() => setFilterType(type)}
              className={`px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-all border ${
                filterType === type
                  ? 'bg-slate-100 border-slate-200/80 text-slate-800 font-bold shadow-sm'
                  : 'border-transparent text-slate-500 hover:text-slate-800'
              }`}
            >
              {getFilterLabel(type)}
            </button>
          ))}
        </div>
        <div className="flex items-center gap-1.5 text-xs text-slate-500 font-bold font-mono">
          <Filter className="w-3.5 h-3.5" />
          ACTIVE SUITES: {filteredSuites.length}
        </div>
      </div>

      {/* Create Suite Modal */}
      {showCreate && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-md border border-slate-200 animate-in fade-in zoom-in-95 duration-150 overflow-hidden flex flex-col">
            <div className="px-6 py-4 border-b border-slate-200 flex items-center justify-between bg-slate-50/50">
              <h2 className="text-lg font-bold text-slate-800 font-sans">Define Validation Suite</h2>
              <button
                onClick={() => setShowCreate(false)}
                className="p-1.5 hover:bg-slate-100 rounded-lg text-slate-500 hover:text-slate-800 transition"
              >
                <Plus className="w-5 h-5 rotate-45" />
              </button>
            </div>
            
            <form onSubmit={handleCreate} className="p-6 space-y-5">
              <div>
                <label className="block text-xs font-bold uppercase tracking-wider text-slate-550 mb-1">Validation Suite Name</label>
                <input
                  type="text"
                  value={newSuite.name}
                  onChange={(e) => setNewSuite({ ...newSuite, name: e.target.value })}
                  className="w-full px-3.5 py-2.5 bg-white border border-slate-200 rounded-xl focus:outline-none focus:border-indigo-500 text-sm placeholder-slate-400 text-slate-800 focus:ring-1 focus:ring-indigo-500 transition-all"
                  placeholder="e.g. Core Checkout API Load Suite"
                  required
                />
              </div>
              
              <div>
                <label className="block text-xs font-bold uppercase tracking-wider text-slate-550 mb-1">Suite Purpose & Description</label>
                <textarea
                  value={newSuite.description}
                  onChange={(e) => setNewSuite({ ...newSuite, description: e.target.value })}
                  className="w-full px-3.5 py-2.5 bg-white border border-slate-200 rounded-xl focus:outline-none focus:border-indigo-500 text-sm h-24 placeholder-slate-400 text-slate-800 focus:ring-1 focus:ring-indigo-500 transition-all resize-none"
                  placeholder="What endpoints, SLOs, or user scenarios does this suite validate?"
                />
              </div>
              
              <div>
                <label className="block text-xs font-bold uppercase tracking-wider text-slate-550 mb-1.5 flex items-center gap-1">
                  Suite Execution Mode
                  <HelpCircle className="w-3.5 h-3.5 text-slate-450" />
                </label>
                <select
                  value={newSuite.test_type}
                  onChange={(e) => setNewSuite({ ...newSuite, test_type: e.target.value as any })}
                  className="w-full px-3.5 py-2.5 bg-white border border-slate-200 rounded-xl focus:outline-none focus:border-indigo-500 text-sm text-slate-700 focus:ring-1 focus:ring-indigo-500 transition-all cursor-pointer"
                >
                  <option value="load">Performance Load Simulation (High-concurrency virtual user loading)</option>
                  <option value="api">REST API Assertion Flow (Functional endpoint validation sequence)</option>
                  <option value="e2e">User Journey E2E Flow (Simulated browser session testing)</option>
                  <option value="unit">Isolated Unit Check (Micro-service assert checkups)</option>
                </select>
              </div>
 
              <div className="flex justify-end gap-2 pt-4 border-t border-slate-205">
                <button
                  type="button"
                  onClick={() => setShowCreate(false)}
                  className="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 border border-slate-200 rounded-lg transition"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-indigo-650 hover:bg-indigo-700 text-white text-sm font-semibold rounded-lg shadow-sm transition"
                >
                  Confirm Suite Creation
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
 
      {/* Suites List */}
      {filteredSuites.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredSuites.map((suite) => (
            <div
              key={suite.id}
              className="glass-panel glass-panel-hover p-6 rounded-xl border border-slate-200 flex flex-col justify-between hover:shadow-md transition group"
            >
              <div>
                <div className="flex items-start justify-between gap-3 mb-3">
                  <h3 className="font-bold text-lg text-slate-800 group-hover:text-indigo-600 transition-colors">
                    {suite.name}
                  </h3>
                </div>
                
                <div className="mb-4">
                  {getSuiteTypeBadge(suite.test_type)}
                </div>
 
                <p className="text-sm text-slate-600 line-clamp-3 mb-6 leading-relaxed font-medium">
                  {suite.description || 'No description provided. Click Configure Suite to define endpoints and criteria.'}
                </p>
 
                {suite.tags && suite.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1.5 mb-6">
                    {suite.tags.map((tag) => (
                      <span
                        key={tag}
                        className="px-2 py-0.5 text-xs bg-slate-50 border border-slate-200 text-slate-500 rounded font-semibold font-mono"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                )}
              </div>
 
              <div className="flex items-center gap-2 border-t border-slate-200 pt-4 mt-auto">
                <Link
                  to={`/projects/${projectId}/suites/${suite.id}`}
                  className="flex-1 flex items-center justify-center gap-2 px-3 py-2 bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 font-semibold rounded-lg text-xs transition"
                >
                  <Settings className="w-3.5 h-3.5" />
                  Configure Suite
                </Link>
                <button
                  onClick={() => triggerRun.mutate(suite.id)}
                  disabled={triggerRun.isPending}
                  className="flex items-center justify-center gap-2 px-4 py-2 bg-emerald-50 hover:bg-emerald-100 border border-emerald-200 text-emerald-700 font-bold rounded-lg text-xs transition disabled:opacity-50"
                >
                  <Play className="w-3.5 h-3.5 fill-current" />
                  {triggerRun.isPending ? 'Executing...' : 'Execute Suite'}
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="glass-panel rounded-xl border border-slate-200 p-12 text-center text-slate-500 shadow-sm">
          <Folder className="w-12 h-12 mx-auto mb-3 text-slate-400" />
          <p className="font-semibold text-slate-700">No validation suites found</p>
          <p className="text-sm text-slate-500 mt-1">Create a new suite or select a different filter category above.</p>
        </div>
      )}

      {project && (
        <ProjectModal
          isOpen={isProjectModalOpen}
          onClose={() => setIsProjectModalOpen(false)}
          projectToEdit={project}
        />
      )}
    </div>
  )
}
