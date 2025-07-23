-- +goose Up
ALTER TABLE chirps
ALTER COLUMN user_id SET NOT NULL;