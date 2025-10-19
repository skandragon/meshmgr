-- name: CreateUser :one
INSERT INTO users (email, password_hash, display_name)
VALUES (@email, @password_hash, @display_name)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = @id;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = @email;

-- name: UpdateUser :one
UPDATE users
SET
    display_name = COALESCE(sqlc.narg('display_name'), display_name),
    password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = @id;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;
