import { useEffect, useState } from 'react';
import { apiService } from '../services/api';
import type { LeaderboardEntry } from '../types/api';
import './Leaderboard.css';

export function MyLeaderboard() {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const load = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const res = await apiService.getMyLeaderboard(0, 10);
        setEntries(Array.isArray(res.data) ? res.data : []);
        // #region agent log
        fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'leaderboard-my-debug',hypothesisId:'H3',location:'src/pages/MyLeaderboard.tsx:load:success',message:'my leaderboard loaded',data:{count:res.data?.length ?? -1,hasTotal:typeof res.total === 'number'},timestamp:Date.now()})}).catch(()=>{});
        // #endregion
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load my leaderboard');
        // #region agent log
        fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'leaderboard-my-debug',hypothesisId:'H3',location:'src/pages/MyLeaderboard.tsx:load:error',message:'my leaderboard failed',data:{isError:err instanceof Error},timestamp:Date.now()})}).catch(()=>{});
        // #endregion
      } finally {
        setIsLoading(false);
      }
    };
    void load();
  }, []);

  return (
    <div className="page-container fade-in">
      <h1>My Leaderboard</h1>
      {error && <div className="error-message">{error}</div>}
      {isLoading && entries.length === 0 ? (
        <div>Loading leaderboard...</div>
      ) : entries.length === 0 ? (
        <div className="empty-state">No leaderboard entries found</div>
      ) : (
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
                <tr key={entry.id ?? `my-${index}`}>
                  <td>{entry.rank ?? index + 1}</td>
                  <td className="username-cell">{entry.username ?? entry.user_id ?? 'N/A'}</td>
                  <td className="score-cell">{entry.score.toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

