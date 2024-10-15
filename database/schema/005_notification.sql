-- +goose Up
CREATE TABLE notifications(
    request_by VARCHAR(12) NOT NULL,
    request_to VARCHAR(12) NOT NULL,
    request_status VARCHAR(8) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE notifications;