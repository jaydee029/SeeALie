-- +goose Up
CREATE TABLE contacts(
    id uuid PRIMARY KEY,
    username VARCHAR(12) NOT NULL,
    room_id uuid NOT NULL,
    connected_on TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE contacts;