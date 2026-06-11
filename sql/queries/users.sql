-- name: CreateUser :one
INSERT INTO users(name, email, password, phone_number, business_name, is_admin)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING id, email, phone_number, business_name, is_admin, created_at, updated_at, verified, blacklisted;

-- name: VerifyUserEmail :one
UPDATE users 
SET verified = TRUE
WHERE email = $1
RETURNING id, email, phone_number, business_name, is_admin, created_at, updated_at, verified, blacklisted;

-- name: CheckUserExistsViaEmail :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE email = $1
);

-- name: ChangePassword :exec
UPDATE users
SET password = $1
WHERE email = $2;

-- name: GetUserDetails :one
SELECT id, email, password, phone_number, business_name, is_admin, created_at, updated_at, verified, blacklisted
FROM users WHERE email = $1;