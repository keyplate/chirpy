-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIME NOT NULL,
    updated_at TIME NOT NULL,
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
