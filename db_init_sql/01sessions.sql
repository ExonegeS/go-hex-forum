CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    session_hash CHAR(64) UNIQUE NOT NULL, -- SHA256
    avatar_url TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);