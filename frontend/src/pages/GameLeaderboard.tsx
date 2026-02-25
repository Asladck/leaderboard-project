import { FormEvent, useEffect, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { apiService } from '../services/api';
import type { LeaderboardEntry } from '../types/api';
import './Leaderboard.css';

export function GameLeaderboard() {
  const { gameId } = useParams<{ gameId: string }>();
  const resolvedGameId = useMemo(() => gameId ?? '', [gameId]);

  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [limit] = useState(10);
  const [total, setTotal] = useState(0);

  const [score, setScore] = useState<number>(0);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitSuccess, setSubmitSuccess] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const load = async () => {
      if (!resolvedGameId) return;
      setIsLoading(true);
      setError(null);
      try {
        const res = await apiService.getGameLeaderboardTop({ gameId: resolvedGameId, offset, limit });
        setEntries(res.data);
        setTotal(res.total ?? 0);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load game leaderboard');
      } finally {
        setIsLoading(false);
      }
    };
    void load();
  }, [resolvedGameId, offset, limit]);

  const handleSubmitScore = async (e: FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    setSubmitSuccess(null);
    setIsSubmitting(true);
    try {
      await apiService.submitScore({ gameId: resolvedGameId, score });
      setSubmitSuccess('Score submitted');
      // Refresh leaderboard after submit
      const res = await apiService.getGameLeaderboardTop({ gameId: resolvedGameId, offset, limit });
      setEntries(res.data);
      setTotal(res.total ?? 0);
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'Failed to submit score');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handlePrevious = () => setOffset((v) => Math.max(0, v - limit));
  const handleNext = () => {
    if (offset + limit < total) setOffset((v) => v + limit);
  };

  if (!resolvedGameId) {
    return <div className="page-container">Invalid game id</div>;
  }

  return (
    <div className="page-container fade-in">
      <h1>Game Leaderboard</h1>

      <div className="leaderboard-table" style={{ marginBottom: '1.5rem' }}>
        <div style={{ padding: '1rem' }}>
          <form onSubmit={handleSubmitScore} style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
            <label style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
              <span style={{ fontWeight: 600, color: '#2c3e50' }}>Submit score</span>
              <input
                type="number"
                min={0}
                step={1}
                value={Number.isFinite(score) ? score : 0}
                onChange={(e) => setScore(Number(e.target.value))}
                style={{ padding: '0.6rem', border: '1px solid #ddd', borderRadius: 6 }}
                required
              />
            </label>
            <button
              className="pagination-btn"
              style={{ height: 40, alignSelf: 'flex-end' }}
              type="submit"
              disabled={isSubmitting}
            >
              {isSubmitting ? 'Submitting...' : 'Submit'}
            </button>
          </form>
          {submitError && <div className="error-message" style={{ marginTop: '0.75rem' }}>{submitError}</div>}
          {submitSuccess && <div className="success-message" style={{ marginTop: '0.75rem' }}>{submitSuccess}</div>}
        </div>
      </div>

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
                  <th>Date</th>
                </tr>
              </thead>
              <tbody>
                {entries.map((entry, index) => (
                  <tr key={entry.id ?? `${offset}-${index}`}>
                    <td>{entry.rank ?? offset + index + 1}</td>
                    <td className="username-cell">{entry.username ?? entry.user_id ?? 'N/A'}</td>
                    <td className="score-cell">{entry.score.toLocaleString()}</td>
                    <td>{entry.created_at ? new Date(entry.created_at).toLocaleDateString() : 'N/A'}</td>
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

