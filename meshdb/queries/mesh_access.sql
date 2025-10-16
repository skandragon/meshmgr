-- name: GrantMeshAccess :one
INSERT INTO mesh_access (mesh_id, user_id, access_level, granted_by)
VALUES (@mesh_id, @user_id, @access_level, @granted_by)
RETURNING *;

-- name: GetMeshAccess :one
SELECT * FROM mesh_access
WHERE mesh_id = @mesh_id AND user_id = @user_id;

-- name: UpdateMeshAccess :one
UPDATE mesh_access
SET
    access_level = @access_level
WHERE mesh_id = @mesh_id AND user_id = @user_id
RETURNING *;

-- name: RevokeMeshAccess :exec
DELETE FROM mesh_access
WHERE mesh_id = @mesh_id AND user_id = @user_id;

-- name: ListMeshAccessByMesh :many
SELECT ma.*, u.email, u.display_name
FROM mesh_access ma
JOIN users u ON ma.user_id = u.id
WHERE ma.mesh_id = @mesh_id
ORDER BY ma.created_at DESC;

-- name: ListMeshAccessByUser :many
SELECT ma.*, m.name as mesh_name
FROM mesh_access ma
JOIN meshes m ON ma.mesh_id = m.id
WHERE ma.user_id = @user_id
ORDER BY ma.created_at DESC;

-- name: CheckUserMeshAccess :one
SELECT access_level FROM mesh_access
WHERE mesh_id = @mesh_id AND user_id = @user_id;
