interface StatusBadgeProps {
  status: string
}

const statusColors: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  running: 'bg-blue-100 text-blue-800',
  passed: 'bg-green-100 text-green-800',
  failed: 'bg-red-100 text-red-800',
  cancelled: 'bg-gray-100 text-gray-800',
  error: 'bg-red-100 text-red-800',
  skipped: 'bg-gray-100 text-gray-800',
  online: 'bg-green-100 text-green-800',
  offline: 'bg-gray-100 text-gray-800',
  busy: 'bg-yellow-100 text-yellow-800',
  draining: 'bg-orange-100 text-orange-800',
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const colorClass = statusColors[status] || 'bg-gray-100 text-gray-800'

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colorClass}`}>
      {status}
    </span>
  )
}
