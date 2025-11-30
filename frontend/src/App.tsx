import './App.css'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import AccountPage from './pages/Account'
import Dashboard from './pages/Dashboard'
import FilesPage from './pages/Files'
import UploadPage from './pages/Upload'
import Login from './pages/Login'
import Register from './pages/Register'
import HomePage from './pages/Home'

function App() {
  return (
        <Router>
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/account" element={<AccountPage />} />
                <Route path="/dashboard" element={<Dashboard />} />
                <Route path="/files" element={<FilesPage />} />
                <Route path="/upload" element={<UploadPage />} />
                <Route path="/files" element={<FilesPage />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
            </Routes>
        </Router>
  )
}

export default App
