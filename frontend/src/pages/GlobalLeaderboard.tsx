import { useEffect, useState } from 'react';
import { apiService } from '../services/api';
import type { LeaderboardEntry } from '../types/api';
import './Leaderboard.css';

export function GlobalLeaderboard() {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [limit] = useState(10);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    const load = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const res = await apiService.getGlobalLeaderboard(offset, limit);
        setEntries(res.data);
        setTotal(res.total ?? 0);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load global leaderboard');
      } finally {
        setIsLoading(false);
      }
    };
    void load();
  }, [offset, limit]);

  const handlePrevious = () => setOffset((v) => Math.max(0, v - limit));
  const handleNext = () => {
    if (offset + limit < total) setOffset((v) => v + limit);
  };

  return (
    <div className="page-container">
      <h1>Global Leaderboard</h1>
      {error && <div className="error-message">{error}</div>}
      {isLoading && entries.length === 0 ? (
        <div>Loading leaderboard...</div>
      ) : entries.length === 0 ? (
        <div className="empty-state">No leaderboard entries found</div>
      ) : (
        <>
          <div className="leaderboard-table">
            <table>
              <thead>
                <tr>
                  <th>Rank</th>
                  <th>User</th>
                  <th>Score</th>
                </tr>
              </thead>
              <tbody>
                {entries.map((entry, index) => (
                  <tr key={entry.id ?? `${offset}-${index}`}>
                    <td>{entry.rank ?? offset + index + 1}</td>
                    <td className="username-cell">{entry.username ?? entry.user_id ?? 'N/A'}</td>
                    <td className="score-cell">{entry.score.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
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

