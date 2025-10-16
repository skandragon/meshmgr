-- name: CreateNode :one
INSERT INTO nodes (mesh_id, hardware_id, name, long_name, role, public_key, private_key, status, unmessageable)
VALUES (@mesh_id, @hardware_id, @name, @long_name, @role, @public_key, @private_key, @status, @unmessageable)
RETURNING *;

-- name: GetNode :one
SELECT * FROM nodes
WHERE id = @id;

-- name: GetNodeByHardwareID :one
SELECT * FROM nodes
WHERE mesh_id = @mesh_id AND hardware_id = @hardware_id;

-- name: ListNodesByMesh :many
SELECT * FROM nodes
WHERE mesh_id = @mesh_id
ORDER BY name ASC;

-- name: UpdateNode :one
UPDATE nodes
SET
    name = COALESCE(sqlc.narg('name'), name),
    long_name = COALESCE(sqlc.narg('long_name'), long_name),
    role = COALESCE(sqlc.narg('role'), role),
    public_key = COALESCE(sqlc.narg('public_key'), public_key),
    private_key = COALESCE(sqlc.narg('private_key'), private_key),
    status = COALESCE(sqlc.narg('status'), status),
    unmessageable = COALESCE(sqlc.narg('unmessageable'), unmessageable),
    last_seen = COALESCE(sqlc.narg('last_seen'), last_seen),
    pending_changes = COALESCE(sqlc.narg('pending_changes'), pending_changes),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: UpdateNodeStatus :one
UPDATE nodes
SET
    status = @status,
    last_seen = NOW(),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteNode :exec
DELETE FROM nodes
WHERE id = @id;

-- name: CountNodesByMesh :one
SELECT COUNT(*) FROM nodes
WHERE mesh_id = @mesh_id;

-- name: UpdateNodeAppliedState :one
UPDATE nodes
SET
    applied_name = @applied_name,
    applied_long_name = @applied_long_name,
    applied_role = @applied_role,
    applied_public_key = @applied_public_key,
    applied_private_key = @applied_private_key,
    applied_unmessageable = @applied_unmessageable,
    config_applied_at = NOW(),
    pending_changes = FALSE,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: ListNodesWithPendingChanges :many
SELECT * FROM nodes
WHERE mesh_id = @mesh_id AND pending_changes = TRUE
ORDER BY name ASC;
