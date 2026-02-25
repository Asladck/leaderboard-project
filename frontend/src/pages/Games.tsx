import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiService } from '../services/api';
import type { Game } from '../types/api';
import './Games.css';

export function Games() {
  const navigate = useNavigate();
  const [games, setGames] = useState<Game[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [limit] = useState(10);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadGames();
  }, [offset]);

  const loadGames = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiService.getGames(offset, limit);
      setGames(response.data);
      setTotal(response.total || 0);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load games');
    } finally {
      setIsLoading(false);
    }
  };

  const handlePrevious = () => {
    if (offset > 0) {
      setOffset(Math.max(0, offset - limit));
    }
  };

  const handleNext = () => {
    if (offset + limit < total) {
      setOffset(offset + limit);
    }
  };

  if (isLoading && games.length === 0) {
    return <div className="page-container">Loading games...</div>;
  }

  return (
    <div className="page-container">
      <h1>Games</h1>
      {error && <div className="error-message">{error}</div>}
      {games.length === 0 && !isLoading ? (
        <div className="empty-state">No games found</div>
      ) : (
        <>
          <div className="games-grid">
            {games.map((game) => (
              <button
                key={game.id}
                className="game-card game-card-button"
                type="button"
                onClick={() => navigate(`/games/${game.id}`)}
              >
                <h3>{game.name}</h3>
                {game.description && <p>{game.description}</p>}
                {game.created_at && (
                  <div className="game-meta">
                    Created: {new Date(game.created_at).toLocaleDateString()}
                  </div>
                )}
              </button>
            ))}
          </div>
          <div className="pagination">
            <button
              onClick={handlePrevious}
              disabled={offset === 0 || isLoading}
              className="pagination-btn"
            >
              Previous
            </button>
            <span className="pagination-info">
              Showing {offset + 1}-{Math.min(offset + limit, total)} of {total}
            </span>
            <button
              onClick={handleNext}
              disabled={offset + limit >= total || isLoading}
              className="pagination-btn"
            >
              Next
            </button>
          </div>
        </>
      )}
    </div>
  );
}
