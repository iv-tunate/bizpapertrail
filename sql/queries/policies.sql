-- name: AddPolicies :many
INSERT INTO policies(title, description, category)
SELECT
    unnest(@titles::text[]),
    unnest(@descriptions::text[]),
    unnest(@categories::text[])
RETURNING *;

-- name: AddPolicy :one
INSERT INTO policies(title, description, category)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetPolicies :many
SELECT * FROM policies
ORDER BY category;

-- name: GetPoliciesByCategory :many
SELECT * FROM policies
WHERE category = $1;

-- name: GetPolicyByID :one
SELECT  * FROM policies
WHERE id = $1;

-- name: UpdatePolicy :one
UPDATE policies
SET 
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    category = COALESCE($4, category),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePolicy :exec
DELETE FROM policies
WHERE id = $1;

-- name: GetUserPolicies :many
SELECT p.* FROM policies p
INNER JOIN user_policies up ON up.policy_id = p.id
WHERE up.user_id = $1
ORDER BY category;

-- name: AddUserPolicy :one
INSERT INTO user_policies(user_id, policy_id)
VALUES($1, $2)
RETURNING *;

-- name: AddUserPolicies :many
INSERT INTO user_policies(user_id, policy_id)
SELECT
    unnest(@user_id::UUID[]),
    unnest(@policy_id::UUID[])
RETURNING *;

-- name: RemoveUserPolicy :exec
DELETE FROM user_policies
WHERE user_id = $1 AND policy_id = $2;

-- name: GetUserPolicyIDs :many
SELECT policy_id FROM user_policies
WHERE user_id = $1;

