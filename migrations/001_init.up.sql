CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- USERS
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT now()
);

-- GAMES / ACTIVITIES
CREATE TABLE games (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       name TEXT UNIQUE NOT NULL
);

-- SCORE HISTORY
CREATE TABLE score_history (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                               game_id UUID REFERENCES games(id) ON DELETE CASCADE,
                               score INT NOT NULL,
                               created_at TIMESTAMP DEFAULT now()
);

-- INDEXES
CREATE INDEX idx_score_user ON score_history(user_id);
CREATE INDEX idx_score_game ON score_history(game_id);
CREATE INDEX idx_score_created ON score_history(created_at);
