-- +goose Up
CREATE TABLE connections(
    request_by VARCHAR(12) NOT NULL,
    request_to VARCHAR(12) NOT NULL,
    connection_id uuid UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    sent_attempts INT DEFAULT 0 NOT NULL CHECK (request_status <=3),
    status_sent VARCHAR(7) DEFAULT "PENDING" NOT NULL,
    status_accepted BOOLEAN DEFAULT FALSE NOT NULL
);

-- +goose Down
DROP TABLE connection;