-- +goose Up
CREATE TABLE connections(
    request_by VARCHAR(12) NOT NULL,
    request_to VARCHAR(12) NOT NULL,
    connection_id uuid UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    sent_attempts INT DEFAULT 0 CHECK (sent_attempts <=3),
    status_sent BOOLEAN NOT NULL DEFAULT FALSE CHECK (status_sent IN (FALSE, TRUE)),
    FOREIGN KEY (request_by) REFERENCES users (user_id),
    FOREIGN KEY (request_to) REFERENCES users (user_id),
    CHECK (request_by != request_to),
    PRIMARY KEY (request_by, request_to)
    );

CREATE UNIQUE INDEX unique_constraint_conn ON connections (LEAST(request_by, request_to), GREATEST(request_by, request_to));

-- +goose Down
DROP TABLE connection;

DROP INDEX IF EXISTS unique_constraint_conn;