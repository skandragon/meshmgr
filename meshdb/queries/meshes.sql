-- name: CreateMesh :one
INSERT INTO meshes (owner_id, name, description)
VALUES (@owner_id, @name, @description)
RETURNING *;

-- name: GetMeshByID :one
SELECT * FROM meshes
WHERE id = @id;

-- name: UpdateMesh :one
UPDATE meshes
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteMesh :exec
DELETE FROM meshes
WHERE id = @id;

-- name: ListMeshesByOwner :many
SELECT * FROM meshes
WHERE owner_id = @owner_id
ORDER BY created_at DESC;

-- name: ListMeshesByUser :many
SELECT DISTINCT m.* FROM meshes m
LEFT JOIN mesh_access ma ON m.id = ma.mesh_id
WHERE m.owner_id = @user_id OR ma.user_id = @user_id
ORDER BY m.created_at DESC;
