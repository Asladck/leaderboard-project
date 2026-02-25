import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import { Layout } from './components/Layout';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Games } from './pages/Games';
import { GameLeaderboard } from './pages/GameLeaderboard';
import { GlobalLeaderboard } from './pages/GlobalLeaderboard';
import { CreateGame } from './pages/CreateGame';
import './App.css';

function AppRoutes() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <div>Loading...</div>
      </div>
    );
  }

  return (
    <Routes>
      <Route path="/login" element={isAuthenticated ? <Navigate to="/games" replace /> : <Login />} />
      <Route path="/register" element={isAuthenticated ? <Navigate to="/games" replace /> : <Register />} />
      <Route
        path="/games"
        element={
          <ProtectedRoute>
            <Layout>
              <Games />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/games/:gameId"
        element={
          <ProtectedRoute>
            <Layout>
              <GameLeaderboard />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/leaderboard/global"
        element={
          <ProtectedRoute>
            <Layout>
              <GlobalLeaderboard />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/leaderboard/my"
        element={<Navigate to="/leaderboard/global" replace />}
      />
      <Route
        path="/admin/create"
        element={
          <ProtectedRoute>
            <Layout>
              <CreateGame />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route path="/" element={<Navigate to="/games" replace />} />
      <Route path="*" element={<Navigate to="/games" replace />} />
    </Routes>
  );
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
