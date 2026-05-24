interface StatusBadgeProps {
  status: string
}

const statusStyles: Record<string, string> = {
  pending: 'bg-amber-50 border border-amber-200 text-amber-700',
  running: 'bg-blue-50 border border-blue-200 text-blue-700 animate-pulse',
  passed: 'bg-emerald-50 border border-emerald-200 text-emerald-700',
  failed: 'bg-rose-50 border border-rose-200 text-rose-700',
  cancelled: 'bg-slate-100 border border-slate-200 text-slate-600',
  error: 'bg-rose-50 border border-rose-200 text-rose-700',
  skipped: 'bg-slate-100 border border-slate-200 text-slate-500',
  online: 'bg-emerald-50 border border-emerald-200 text-emerald-700',
  offline: 'bg-slate-100 border border-slate-200 text-slate-500',
  busy: 'bg-amber-50 border border-amber-200 text-amber-700',
  draining: 'bg-orange-50 border border-orange-200 text-orange-700',
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const style = statusStyles[status] || 'bg-slate-100 border border-slate-200 text-slate-600'
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold capitalize ${style}`}>
      {status}
    </span>
  )
}
