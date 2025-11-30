import './App.css'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import AccountPage from './pages/Account'
import Dashboard from './pages/Dashboard'
import FilesPage from './pages/Files'
import UploadPage from './pages/Upload'
import Login from './pages/Login'
import Register from './pages/Register'
import HomePage from './pages/Home'
import { AuthProvider } from './hooks/AuthContext'
import ProtectedRoute from './hooks/ProtectedRoute'

function App() {
  return (
    <AuthProvider>
      <Router>
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />

                <Route path="/account" element={<ProtectedRoute> <AccountPage /> </ProtectedRoute>} />
                <Route path="/dashboard" element={<ProtectedRoute>  <Dashboard /> </ProtectedRoute>} />
                <Route path="/files" element={<ProtectedRoute> <FilesPage /> </ProtectedRoute>} />
                <Route path="/upload" element={<ProtectedRoute> <UploadPage /> </ProtectedRoute>} />
                <Route path="/files" element={<ProtectedRoute> <FilesPage /> </ProtectedRoute>} />

            </Routes>
        </Router>
    </AuthProvider>
  )
}

export default App
