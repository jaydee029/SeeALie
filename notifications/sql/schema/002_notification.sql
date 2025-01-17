-- +goose Up
CREATE TABLE notifications(
    request_init_by VARCHAR(12) NOT NULL,
    request_to VARCHAR(12) NOT NULL,
    request_status VARCHAR(8) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    sent_attempts INT DEFAULT 0 CHECK (sent_attempts <=3),
    status_sent BOOLEAN DEFAULT FALSE CHECK (status_sent IN (FALSE, TRUE)),
    FOREIGN KEY (request_init_by) REFERENCES users (user_id),
    FOREIGN KEY (request_to) REFERENCES users (user_id),
    CHECK (request_init_by != request_to),
    PRIMARY KEY (request_init_by, request_to)
    );

CREATE UNIQUE INDEX unique_constraint_ntf ON notifications (LEAST(request_init_by, request_to), GREATEST(request_init_by, request_to));

-- +goose Down
DROP TABLE notifications;

DROP INDEX IF EXISTS unique_constraint_ntf;