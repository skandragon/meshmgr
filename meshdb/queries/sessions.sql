-- name: CreateSession :one
INSERT INTO sessions (user_id, token, expires_at)
VALUES (@user_id, @token, @expires_at)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE token = @token AND expires_at > NOW();

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = @token;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = @user_id;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= NOW();

-- name: GetUserSessions :many
SELECT * FROM sessions
WHERE user_id = @user_id AND expires_at > NOW()
ORDER BY created_at DESC;
