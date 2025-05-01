CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(id),
    parent_id INTEGER REFERENCES comments(id), -- NULL for top-level
    session_id INTEGER REFERENCES sessions(id),
    content TEXT NOT NULL,
    image_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);