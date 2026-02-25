import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiService } from '../services/api';
import type { LoginRequest } from '../types/api';

interface AuthContextType {
  isAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    setIsAuthenticated(!!token);
    setIsLoading(false);

    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'H5',location:'src/contexts/AuthContext.tsx:useEffect:init',message:'auth init from localStorage',data:{hasAccessToken:!!token},timestamp:Date.now()})}).catch(()=>{});
    // #endregion
  }, []);

  const login = async (credentials: LoginRequest) => {
    try {
      // #region agent log
      fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'H5',location:'src/contexts/AuthContext.tsx:login:pre',message:'auth login called',data:{hasUsername:!!credentials.username},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
      const response = await apiService.login(credentials);
      localStorage.setItem('access_token', response.access_token);
      localStorage.setItem('refresh_token', response.refresh_token);
      setIsAuthenticated(true);
      navigate('/games');

      // #region agent log
      fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'H5',location:'src/contexts/AuthContext.tsx:login:success',message:'auth login success',data:{storedAccessToken:true,storedRefreshToken:true},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
    } catch (error) {
      // #region agent log
      fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'pre-fix',hypothesisId:'H5',location:'src/contexts/AuthContext.tsx:login:error',message:'auth login failed',data:{isError: error instanceof Error},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    setIsAuthenticated(false);
    navigate('/login');
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
