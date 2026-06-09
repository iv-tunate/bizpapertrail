-- +goose Up
ALTER TABLE users 
ADD COLUMN IsAdmin BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE users 
DROP COLUMN is_admin;