-- +goose Up
BEGIN;

ALTER TABLE users
ADD COLUMN hashed_password TEXT;

UPDATE users
SET hashed_password = 'unset'
WHERE hashed_password IS NULL;

ALTER TABLE users
ALTER COLUMN hashed_password SET NOT NULL;

COMMIT;