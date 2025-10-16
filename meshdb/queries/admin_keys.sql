-- name: CreateAdminKey :one
INSERT INTO admin_keys (mesh_id, public_key, key_name, added_by)
VALUES (@mesh_id, @public_key, @key_name, @added_by)
RETURNING *;

-- name: GetAdminKey :one
SELECT * FROM admin_keys
WHERE id = @id;

-- name: ListAdminKeysByMesh :many
SELECT * FROM admin_keys
WHERE mesh_id = @mesh_id
ORDER BY created_at DESC;

-- name: DeleteAdminKey :exec
DELETE FROM admin_keys
WHERE id = @id;

-- name: CountAdminKeysByMesh :one
SELECT COUNT(*) FROM admin_keys
WHERE mesh_id = @mesh_id;
