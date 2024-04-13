-- +goose Up
CREATE TABLE apikey(
    id uuid NOT NULL REFERENCES users(id),
    api_key VARCHAR(64) DEFAULT encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down
DROP TABLE apikey;
