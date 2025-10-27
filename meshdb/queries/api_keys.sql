-- name: CreateAPIKey :one
INSERT INTO user_api_keys (user_id, key_hash, key_name, expires_at)
VALUES (@user_id, @key_hash, @key_name, @expires_at)
RETURNING *;

-- name: GetAPIKeyByHash :one
SELECT * FROM user_api_keys
WHERE key_hash = @key_hash;

-- name: GetAPIKey :one
SELECT * FROM user_api_keys
WHERE id = @id;

-- name: ListAPIKeysByUser :many
SELECT * FROM user_api_keys
WHERE user_id = @user_id
ORDER BY created_at DESC;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE user_api_keys
SET last_used_at = NOW()
WHERE id = @id;

-- name: DeleteAPIKey :exec
DELETE FROM user_api_keys
WHERE id = @id;

-- name: DeleteExpiredAPIKeys :exec
DELETE FROM user_api_keys
WHERE expires_at IS NOT NULL AND expires_at < NOW();

-- name: UpdateAPIKeyHash :one
UPDATE user_api_keys
SET key_hash = @key_hash
WHERE id = @id
RETURNING *;
