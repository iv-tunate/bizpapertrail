-- +goose Up

ALTER TABLE users
ADD CONSTRAINT unique_business_name UNIQUE (business_name);