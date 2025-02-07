-- +goose Up
CREATE TABLE users(
    username VARCHAR(12) NOT NULL,
    email VARCHAR(100) NOT NULL,
    last_updated_at TIMESTAMP NOT NULL
    );

-- +goose Down
DROP TABLE users;