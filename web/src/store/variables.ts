import { create } from 'zustand'

interface VariableState {
  captured: Record<string, unknown>
  setCaptured: (name: string, value: unknown) => void
  setCapturedBatch: (vars: Record<string, unknown>) => void
  clear: () => void
}

export const useVariableStore = create<VariableState>((set) => ({
  captured: {},
  setCaptured: (name, value) =>
    set((state) => ({ captured: { ...state.captured, [name]: value } })),
  setCapturedBatch: (vars) =>
    set((state) => ({ captured: { ...state.captured, ...vars } })),
  clear: () => set({ captured: {} }),
}))
