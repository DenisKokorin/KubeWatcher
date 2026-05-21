import { lazy, Suspense } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { Box, CircularProgress } from '@mui/material'
import Login from './pages/Login'
const Dashboard = lazy(() => import('./pages/Dashboard'))
const Nodes = lazy(() => import('./pages/Nodes'))
const Pods = lazy(() => import('./pages/Pods'))
const Employees = lazy(() => import('./pages/Employees'))
const Reviews = lazy(() => import('./pages/Reviews'))
const ExternalData = lazy(() => import('./pages/ExternalData'))
const AdminUsers = lazy(() => import('./pages/AdminUsers'))
import Sidebar from './components/Sidebar'
import { ProtectedRoute } from './components/ProtectedRoute'
import { AuthProvider, useAuth } from './hooks/useAuth'

const AppContent = () => {
  const { isAuthenticated, loading } = useAuth()

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
        <CircularProgress />
      </Box>
    )
  }

  if (!isAuthenticated) {
    return <Login />
  }

  return (
    <Box sx={{ display: 'flex' }}>
      <Sidebar />
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Suspense fallback={<Box display="flex" justifyContent="center" alignItems="center" minHeight="400px"><CircularProgress /></Box>}>
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/nodes" element={<Nodes />} />
            <Route path="/pods" element={<Pods />} />
            <Route path="/employees" element={<Employees />} />
            <Route path="/reviews" element={<Reviews />} />
            <Route path="/external" element={<ExternalData />} />
            <Route
              path="/admin/users"
              element={
                <ProtectedRoute requiredRole="admin">
                  <AdminUsers />
                </ProtectedRoute>
              }
            />
          </Routes>
        </Suspense>
      </Box>
    </Box>
  )
}

const App = () => (
  <AuthProvider>
    <AppContent />
  </AuthProvider>
)

export default App