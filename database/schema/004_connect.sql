-- +goose Up
CREATE TABLE connections(
    request_by VARCHAR(12) NOT NULL,
    connection_id uuid UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE connection;