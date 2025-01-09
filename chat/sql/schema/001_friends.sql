-- +goose Up
CREATE TABLE friends(
    followed_by VARCHAR(12),
    followed VARCHAR(12),
    room_id uuid NOT NULL,
    connected_at TIMESTAMP NOT NULL,
    PRIMARY KEY (followed_by, followed),
    FOREIGN KEY (followed_by) REFERENCES users(username),
    FOREIGN KEY (followed) REFERENCES users(username)
);

-- +goose Down
DROP TABLE friends;