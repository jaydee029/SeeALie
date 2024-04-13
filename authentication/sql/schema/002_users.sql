-- +goose Up
CREATE TABLE users(
    id uuid PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    passwd bytea NOT NULL,
    username VARCHAR(12) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;