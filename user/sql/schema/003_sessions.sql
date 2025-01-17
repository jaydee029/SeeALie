-- +goose Up
CREATE TABLE sessions(
    session_id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    jwt VARCHAR(100) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) on DELETE CASCADE
);

-- +goose Down
DROP TABLE sessions;