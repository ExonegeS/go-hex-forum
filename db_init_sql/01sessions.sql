CREATE TABLE sessions(
    id SERIAL PRIMARY KEY,
    session_hash TEXT NOT NULL UNIQUE,
    user_id INT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);