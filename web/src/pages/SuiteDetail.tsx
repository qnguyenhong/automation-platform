import { useParams, useNavigate } from 'react-router-dom'
import { useSuite, useCases } from '../hooks/useSuites'
import { useTriggerRun } from '../hooks/useRuns'
import { Play, Plus, Upload } from 'lucide-react'

export default function SuiteDetail() {
  const { projectId, suiteId } = useParams<{ projectId: string; suiteId: string }>()
  const navigate = useNavigate()
  const { data: suite } = useSuite(projectId!, suiteId!)
  const { data: cases } = useCases(projectId!, suiteId!)
  const triggerRun = useTriggerRun(projectId!, suiteId!)

  const handleRun = async () => {
    await triggerRun.mutateAsync({})
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{suite?.name}</h1>
          <p className="text-gray-500">{suite?.description}</p>
        </div>
        <button
          onClick={handleRun}
          disabled={triggerRun.isPending}
          className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50"
        >
          <Play className="w-4 h-4" />
          {triggerRun.isPending ? 'Running...' : 'Run Suite'}
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b border-gray-200 flex items-center justify-between">
          <h3 className="text-lg font-semibold">Test Cases</h3>
          <div className="flex items-center gap-2">
            <button
              onClick={() => navigate(`/projects/${projectId}/openapi-import`)}
              className="flex items-center gap-2 px-3 py-1.5 text-sm bg-blue-50 text-blue-700 hover:bg-blue-100 rounded-lg"
            >
              <Upload className="w-4 h-4" />
              Import from OpenAPI
            </button>
            <button className="flex items-center gap-2 px-3 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg">
              <Plus className="w-4 h-4" />
              Add Case
            </button>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tags</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {cases?.map((testCase) => (
                <tr key={testCase.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium">{testCase.name}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      testCase.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                    }`}>
                      {testCase.enabled ? 'Enabled' : 'Disabled'}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex flex-wrap gap-1">
                      {testCase.tags?.map((tag) => (
                        <span key={tag} className="px-2 py-0.5 text-xs bg-gray-100 rounded">
                          {tag}
                        </span>
                      ))}
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <button className="text-blue-600 hover:text-blue-800">Edit</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
