import { useParams } from 'react-router-dom'
import { useSuites, useCreateSuite } from '../hooks/useSuites'
import { useProject } from '../hooks/useProjects'
import { useState } from 'react'
import { Plus, Play, Settings } from 'lucide-react'
import { Link } from 'react-router-dom'

export default function TestSuites() {
  const { projectId } = useParams<{ projectId: string }>()
  const { data: project } = useProject(projectId!)
  const { data: suites, isLoading } = useSuites(projectId!)
  const createSuite = useCreateSuite(projectId!)
  const [showCreate, setShowCreate] = useState(false)
  const [newSuite, setNewSuite] = useState({
    name: '',
    description: '',
    test_type: 'api' as const,
  })

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    await createSuite.mutateAsync(newSuite)
    setShowCreate(false)
    setNewSuite({ name: '', description: '', test_type: 'api' })
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{project?.name || 'Project'}</h1>
          <p className="text-gray-500">Test Suites</p>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          <Plus className="w-4 h-4" />
          New Suite
        </button>
      </div>

      {/* Create Suite Modal */}
      {showCreate && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-96">
            <h2 className="text-lg font-semibold mb-4">Create Test Suite</h2>
            <form onSubmit={handleCreate} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
                <input
                  type="text"
                  value={newSuite.name}
                  onChange={(e) => setNewSuite({ ...newSuite, name: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  value={newSuite.description}
                  onChange={(e) => setNewSuite({ ...newSuite, description: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Test Type</label>
                <select
                  value={newSuite.test_type}
                  onChange={(e) => setNewSuite({ ...newSuite, test_type: e.target.value as any })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                >
                  <option value="api">API</option>
                  <option value="e2e">E2E</option>
                  <option value="load">Load</option>
                  <option value="unit">Unit</option>
                </select>
              </div>
              <div className="flex justify-end gap-2">
                <button
                  type="button"
                  onClick={() => setShowCreate(false)}
                  className="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                >
                  Create
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Suites List */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {suites?.map((suite) => (
          <div key={suite.id} className="bg-white rounded-lg shadow p-6">
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="font-semibold text-lg">{suite.name}</h3>
                <p className="text-sm text-gray-500">{suite.description}</p>
              </div>
              <span className="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-800 rounded">
                {suite.test_type}
              </span>
            </div>

            <div className="flex flex-wrap gap-2 mb-4">
              {suite.tags?.map((tag) => (
                <span key={tag} className="px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded">
                  {tag}
                </span>
              ))}
            </div>

            <div className="flex items-center gap-2">
              <Link
                to={`/projects/${projectId}/suites/${suite.id}`}
                className="flex-1 flex items-center justify-center gap-2 px-3 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg text-sm"
              >
                <Settings className="w-4 h-4" />
                Configure
              </Link>
              <button className="flex items-center justify-center gap-2 px-3 py-2 bg-green-100 hover:bg-green-200 text-green-700 rounded-lg text-sm">
                <Play className="w-4 h-4" />
                Run
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
