import { FormEvent, useState } from 'react';
import { apiService } from '../services/api';
import './CreateGame.css';

export function CreateGame() {
  const [name, setName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setIsLoading(true);
    try {
      await apiService.createGame(name);
      setSuccess('Game created');
      setName('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create game');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="page-container">
      <h1>Create Game</h1>
      <div className="card">
        <form onSubmit={handleSubmit} className="form">
          <label className="label" htmlFor="gameName">Game name</label>
          <input
            id="gameName"
            className="input"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
          {error && <div className="error-message">{error}</div>}
          {success && <div className="success-message">{success}</div>}
          <button className="primary-btn" disabled={isLoading} type="submit">
            {isLoading ? 'Creating...' : 'Create'}
          </button>
        </form>
      </div>
    </div>
  );
}

