export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
}

export interface Game {
  id: string; // UUID
  name: string;
  description?: string;
  created_at?: string;
  updated_at?: string;
}

export interface LeaderboardEntry {
  id?: string;
  username?: string;
  user_id?: string;
  score: number;
  rank?: number;
  game_id?: string;
  game_name?: string;
  created_at?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total?: number;
  offset?: number;
  limit?: number;
}

export interface ApiError {
  message: string;
  code?: string;
}
