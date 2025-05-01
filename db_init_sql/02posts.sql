CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES sessions(id),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expiry_time TIMESTAMPTZ NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT false
);