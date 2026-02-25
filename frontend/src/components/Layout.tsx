  import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import './Layout.css';

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const { logout } = useAuth();

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="layout">
      <header className="header">
        <div className="header-content">
          <h1 className="logo">Online Leadership</h1>
          <nav className="nav">
            <Link to="/games" className="nav-link">Games</Link>
            <Link to="/leaderboard/global" className="nav-link">Global Leaderboard</Link>
            <Link to="/admin/create" className="nav-link">Create Game</Link>
            <button onClick={handleLogout} className="logout-btn">Logout</button>
          </nav>
        </div>
      </header>
      <main className="main-content">
        {children}
      </main>
    </div>
  );
}
