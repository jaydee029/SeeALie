-- +goose Up
CREATE TABLE revoked(
    token VARCHAR(100) NOT NULL,
    revoked_at timestamp NOT NULL  
);

-- +goose Down
DROP TABLE revoked;
