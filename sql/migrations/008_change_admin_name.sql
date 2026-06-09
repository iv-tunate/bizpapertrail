-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
RENAME COLUMN "isadmin" TO is_admin;
-- +goose StatementEnd