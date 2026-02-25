import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  Game,
  LeaderboardEntry,
  PaginatedResponse,
  ApiError,
} from '../types/api';

type ApiScope = 'api' | 'auth' | 'admin';

class ApiService {
  private getAccessToken(): string | null {
    return localStorage.getItem('access_token');
  }

  private buildHeaders(optionsHeaders?: HeadersInit): Headers {
    const headers = new Headers(optionsHeaders);
    if (!headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json');
    }
    return headers;
  }

  private basePath(scope: ApiScope): string {
    if (scope === 'api') return '/api';
    if (scope === 'auth') return '/auth';
    return '/admin';
  }

  private normalizeLeaderboardResponse(
    raw: unknown,
    ctx: { path: string; offset: number; limit: number }
  ): PaginatedResponse<LeaderboardEntry> {
    if (Array.isArray(raw)) {
      return { data: raw as LeaderboardEntry[], total: raw.length, offset: ctx.offset, limit: ctx.limit };
    }

    if (raw && typeof raw === 'object') {
      const obj = raw as Record<string, unknown>;

      // Common paginated shape
      if (Array.isArray(obj.data)) {
        return {
          data: obj.data as LeaderboardEntry[],
          total: typeof obj.total === 'number' ? obj.total : (obj.data as unknown[]).length,
          offset: typeof obj.offset === 'number' ? obj.offset : ctx.offset,
          limit: typeof obj.limit === 'number' ? obj.limit : ctx.limit,
        };
      }

      // Sometimes APIs return items instead of data
      if (Array.isArray(obj.items)) {
        return {
          data: obj.items as LeaderboardEntry[],
          total: typeof obj.total === 'number' ? obj.total : (obj.items as unknown[]).length,
          offset: typeof obj.offset === 'number' ? obj.offset : ctx.offset,
          limit: typeof obj.limit === 'number' ? obj.limit : ctx.limit,
        };
      }

      // Sentinel shape observed in logs: {"rank":-1} meaning "no entries"
      if (typeof obj.rank === 'number' && (obj.rank as number) === -1) {
        // #region agent log
        fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'leaderboard-my-debug',hypothesisId:'H4',location:'src/services/api.ts:normalizeLeaderboardResponse:sentinel',message:'leaderboard returned sentinel rank=-1; treating as empty list',data:{path:ctx.path},timestamp:Date.now()})}).catch(()=>{});
        // #endregion
        return { data: [], total: 0, offset: ctx.offset, limit: ctx.limit };
      }

      // #region agent log
      fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'leaderboard-my-debug',hypothesisId:'H4',location:'src/services/api.ts:normalizeLeaderboardResponse:unexpected',message:'leaderboard returned unexpected shape; treating as empty list',data:{path:ctx.path,keys:Object.keys(obj).slice(0,12)},timestamp:Date.now()})}).catch(()=>{});
      // #endregion
    }

    return { data: [], total: 0, offset: ctx.offset, limit: ctx.limit };
  }

  private async request<T>(scope: ApiScope, path: string, options: RequestInit = {}): Promise<T> {
    const token = this.getAccessToken();
    const headers = this.buildHeaders(options.headers);
    if (token) headers.set('Authorization', `Bearer ${token}`);

    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'leaderboard-my-debug',hypothesisId:'H1',location:'src/services/api.ts:request:pre',message:'api.request start',data:{scope,path,method:(options.method ?? 'GET'),hasAuthToken:!!token,hasBody:typeof options.body === 'string' ? (options.body as string).length : !!options.body},timestamp:Date.now()})}).catch(()=>{});
    // #endregion

    const response = await fetch(`${this.basePath(scope)}${path}`, {
      ...options,
      headers,
    });

    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'leaderboard-my-debug',hypothesisId:'H1',location:'src/services/api.ts:request:post',message:'api.request response',data:{scope,path,status:response.status,ok:response.ok,contentType:response.headers.get('content-type')},timestamp:Date.now()})}).catch(()=>{});
    // #endregion

    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({
        message: `HTTP ${response.status}: ${response.statusText}`,
      }));
      throw new Error(error.message || `HTTP ${response.status}`);
    }

    const text = await response.text();
    // #region agent log
    fetch('http://127.0.0.1:7242/ingest/ef1b0bce-0a5e-4a32-a929-81ff24b1566f',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({sessionId:'debug-session',runId:'leaderboard-my-debug',hypothesisId:'H2',location:'src/services/api.ts:request:body',message:'api.request body stats',data:{scope,path,textLength:text.length,startsWith:text.slice(0,32)},timestamp:Date.now()})}).catch(()=>{});
    // #endregion
    return (text ? JSON.parse(text) : (undefined as unknown as T));
  }

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    return this.request<LoginResponse>('auth', '/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }
  
  async getGames(offset = 0, limit = 10): Promise<PaginatedResponse<Game>> {
    return this.request<PaginatedResponse<Game>>('admin', `/games?offset=${offset}&limit=${limit}`);
  }

  async getGlobalLeaderboard(
    offset = 0,
    limit = 10
  ): Promise<PaginatedResponse<LeaderboardEntry>> {
    const path = `/leaderboard/global?offset=${offset}&limit=${limit}`;
    const raw = await this.request<unknown>(
      'api',
      path
    );
    return this.normalizeLeaderboardResponse(raw, { path, offset, limit });
  }

  async getMyLeaderboard(
    offset = 0,
    limit = 10
  ): Promise<PaginatedResponse<LeaderboardEntry>> {
    const path = `/leaderboard/my?offset=${offset}&limit=${limit}`;
    const raw = await this.request<unknown>(
      'api',
      path
    );
    return this.normalizeLeaderboardResponse(raw, { path, offset, limit });
  }

  async register(payload: RegisterRequest): Promise<void> {
    await this.request<void>('auth', '/register', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  }

  async createGame(name: string): Promise<void> {
    await this.request<void>('admin', '/create', {
      method: 'POST',
      body: JSON.stringify({ name }),
    });
  }

  async getGameLeaderboardTop(params: {
    gameId: string;
    offset?: number;
    limit?: number;
  }): Promise<PaginatedResponse<LeaderboardEntry>> {
    const { gameId, offset = 0, limit = 10 } = params;
    const path = '/leaderboard/top';
    const raw = await this.request<unknown>('api', path, {
      method: 'POST',
      body: JSON.stringify({ game_id: gameId, offset, limit }),
    });
    return this.normalizeLeaderboardResponse(raw, { path, offset, limit });
  }

  async submitScore(params: { gameId: string; score: number }): Promise<void> {
    const { gameId, score } = params;
    await this.request<void>('api', '/leaderboard/top', {
      method: 'POST',
      body: JSON.stringify({ score, game_id: gameId }),
    });
  }
}

export const apiService = new ApiService();
