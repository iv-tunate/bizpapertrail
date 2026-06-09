-- +goose Up

ALTER TABLE users ALTER COLUMN is_admin SET NOT NULL;