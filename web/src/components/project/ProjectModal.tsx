import { useState, useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { X, Trash2, AlertTriangle, FolderKanban } from 'lucide-react'
import { Project } from '../../api/projects'
import { useCreateProject, useUpdateProject, useDeleteProject } from '../../hooks/useProjects'

interface ProjectModalProps {
  isOpen: boolean
  onClose: () => void
  projectToEdit?: Project | null
}

export default function ProjectModal({ isOpen, onClose, projectToEdit }: ProjectModalProps) {
  const navigate = useNavigate()
  const location = useLocation()
  
  const createProjectMutation = useCreateProject()
  const updateProjectMutation = useUpdateProject()
  const deleteProjectMutation = useDeleteProject()

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [deleteConfirmText, setDeleteConfirmText] = useState('')
  const [error, setError] = useState('')

  const isEditMode = !!projectToEdit

  useEffect(() => {
    if (isOpen) {
      if (projectToEdit) {
        setName(projectToEdit.name || '')
        setDescription(projectToEdit.description || '')
      } else {
        setName('')
        setDescription('')
      }
      setShowDeleteConfirm(false)
      setDeleteConfirmText('')
      setError('')
    }
  }, [isOpen, projectToEdit])

  if (!isOpen) return null

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) {
      setError('Project name is required')
      return
    }

    try {
      if (isEditMode && projectToEdit) {
        await updateProjectMutation.mutateAsync({
          id: projectToEdit.id,
          data: { name: name.trim(), description: description.trim() },
        })
      } else {
        const newProj = await createProjectMutation.mutateAsync({
          name: name.trim(),
          description: description.trim(),
        })
        if (newProj && newProj.id) {
          navigate(`/projects/${newProj.id}/suites`)
        }
      }
      onClose()
    } catch (err) {
      setError('An error occurred while saving the project.')
      console.error(err)
    }
  }

  const handleDelete = async () => {
    if (!projectToEdit) return
    if (deleteConfirmText !== projectToEdit.name) {
      setError('Please type the exact project name to confirm deletion.')
      return
    }

    try {
      await deleteProjectMutation.mutateAsync(projectToEdit.id)
      
      // If user is currently viewing the deleted project, redirect to dashboard
      if (location.pathname.includes(`/projects/${projectToEdit.id}`)) {
        navigate('/')
      }
      
      onClose()
    } catch (err) {
      setError('Failed to delete project. Please try again.')
      console.error(err)
    }
  }

  return (
    <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-md border border-slate-100 animate-in fade-in zoom-in-95 duration-150 overflow-hidden flex flex-col">
        {/* Header */}
        <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between bg-slate-50">
          <div className="flex items-center gap-2 text-slate-800">
            <FolderKanban className="w-5 h-5 text-indigo-600" />
            <h2 className="text-lg font-bold">
              {isEditMode ? 'Configure Project Settings' : 'Initialize Workspace Project'}
            </h2>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 hover:bg-slate-200 rounded-lg text-slate-400 hover:text-slate-700 transition"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Error notification */}
        {error && (
          <div className="bg-rose-50 border-b border-rose-100 px-6 py-3 text-xs font-semibold text-rose-700 flex items-center gap-2">
            <AlertTriangle className="w-4 h-4 shrink-0" />
            <span>{error}</span>
          </div>
        )}

        {/* Content */}
        {!showDeleteConfirm ? (
          <form onSubmit={handleSubmit} className="p-6 space-y-4 flex-1">
            <div>
              <label className="block text-xs font-bold uppercase tracking-wider text-slate-500 mb-1">
                Project Name
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-3.5 py-2 border border-slate-200 rounded-lg focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm placeholder-slate-400"
                placeholder="e.g. Core Billing Engine"
                required
                maxLength={100}
                disabled={createProjectMutation.isPending || updateProjectMutation.isPending}
              />
            </div>

            <div>
              <label className="block text-xs font-bold uppercase tracking-wider text-slate-500 mb-1">
                Description
              </label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="w-full px-3.5 py-2 border border-slate-200 rounded-lg focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm h-28 placeholder-slate-400 resize-none"
                placeholder="Detail the target environments, systems, or purposes of this workspace project..."
                maxLength={500}
                disabled={createProjectMutation.isPending || updateProjectMutation.isPending}
              />
            </div>

            <div className="flex justify-between items-center pt-4 border-t border-slate-100">
              {isEditMode ? (
                <button
                  type="button"
                  onClick={() => setShowDeleteConfirm(true)}
                  className="flex items-center gap-1.5 px-3 py-2 text-xs font-bold text-rose-600 hover:bg-rose-50 rounded-lg transition"
                >
                  <Trash2 className="w-4 h-4" />
                  Delete Workspace
                </button>
              ) : (
                <div />
              )}

              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={onClose}
                  className="px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 border border-slate-200 rounded-lg transition"
                  disabled={createProjectMutation.isPending || updateProjectMutation.isPending}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-semibold rounded-lg shadow-sm transition disabled:opacity-50"
                  disabled={
                    createProjectMutation.isPending ||
                    updateProjectMutation.isPending ||
                    !name.trim()
                  }
                >
                  {isEditMode ? 'Apply Updates' : 'Launch Project'}
                </button>
              </div>
            </div>
          </form>
        ) : (
          <div className="p-6 space-y-4">
            <div className="p-4 bg-amber-50 rounded-lg border border-amber-200 flex gap-3 text-amber-800">
              <AlertTriangle className="w-5 h-5 shrink-0 text-amber-600" />
              <div className="text-xs space-y-1">
                <p className="font-bold">Dangerous Action Ahead</p>
                <p>
                  Deleting this project will permanently remove all associated validation suites,
                  endpoints, and historic execution runs. This action cannot be reversed.
                </p>
              </div>
            </div>

            <div className="space-y-2">
              <p className="text-xs text-slate-500 font-semibold">
                To confirm, type <span className="font-bold text-slate-800 font-mono">"{projectToEdit?.name}"</span> below:
              </p>
              <input
                type="text"
                value={deleteConfirmText}
                onChange={(e) => setDeleteConfirmText(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:border-rose-500 text-sm font-mono"
                placeholder="Type project name..."
                required
              />
            </div>

            <div className="flex justify-end gap-2 pt-4 border-t border-slate-100">
              <button
                type="button"
                onClick={() => {
                  setShowDeleteConfirm(false)
                  setDeleteConfirmText('')
                  setError('')
                }}
                className="px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 border border-slate-200 rounded-lg transition"
                disabled={deleteProjectMutation.isPending}
              >
                Go Back
              </button>
              <button
                type="button"
                onClick={handleDelete}
                className="px-4 py-2 bg-rose-600 hover:bg-rose-700 text-white text-sm font-semibold rounded-lg shadow-sm transition disabled:opacity-50"
                disabled={
                  deleteProjectMutation.isPending ||
                  deleteConfirmText !== projectToEdit?.name
                }
              >
                {deleteProjectMutation.isPending ? 'Deleting...' : 'Permanently Delete'}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
