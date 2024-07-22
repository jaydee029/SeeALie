-- +goose Up
CREATE TABLE id_name(
    user_id uuid PRIMARY KEY REFERENCES users(id),
    username VARCHAR(12) 
);

-- +goose Down
DROP TABLE id_name;
