import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/auth'
import Layout from './components/layout/Layout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import TestSuites from './pages/TestSuites'
import SuiteDetail from './pages/SuiteDetail'
import TestRunDetail from './pages/TestRunDetail'
import ResultsDetail from './pages/ResultsDetail'
import Workers from './pages/Workers'
import WorkerDetail from './pages/WorkerDetail'
import Settings from './pages/Settings'
import NotFound from './pages/NotFound'
import OpenAPIImport from './pages/OpenAPIImport'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((state) => state.token)
  if (!token) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/"
        element={
          <PrivateRoute>
            <Layout />
          </PrivateRoute>
        }
      >
        <Route index element={<Dashboard />} />
        <Route path="projects/:projectId/suites" element={<TestSuites />} />
        <Route path="projects/:projectId/suites/:suiteId" element={<SuiteDetail />} />
        <Route path="projects/:projectId/openapi-import" element={<OpenAPIImport />} />
        <Route path="runs/:runId" element={<TestRunDetail />} />
        <Route path="runs/:runId/results/:resultId" element={<ResultsDetail />} />
        <Route path="workers" element={<Workers />} />
        <Route path="workers/:workerId" element={<WorkerDetail />} />
        <Route path="settings" element={<Settings />} />
      </Route>
      <Route path="*" element={<NotFound />} />
    </Routes>
  )
}
