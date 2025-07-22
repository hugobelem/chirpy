-- +goose Up
CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    email TEXT UNIQUE NOT NULL
);
-- +goose Down
DROP TABLE users;