-- +goose Up
CREATE TABLE friends(
    followed_by uuid PRIMARY KEY REFERENCES users(id_name),
    followed uuid UNIQUE NOT NULL,
    connected_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE friends;