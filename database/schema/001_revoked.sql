-- +goose Up
CREATE TABLE revoked(
    token string NOT NULL,
    revoked_at timestamp NOT NULL  
);

-- +goose Down
DROP TABLE revoked;
