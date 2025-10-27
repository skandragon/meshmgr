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

-- name: ImportNodeConfig :one
-- Import or update node configuration from device scan
INSERT INTO nodes (
    mesh_id,
    hardware_id,
    node_num,
    device_id,
    name,
    long_name,
    short_name,
    firmware_version,
    hw_model,
    public_key,
    private_key,
    raw_device_config,
    config_imported_at,
    last_seen,
    status
) VALUES (
    @mesh_id,
    @hardware_id,
    @node_num,
    @device_id,
    @name,
    @long_name,
    @short_name,
    @firmware_version,
    @hw_model,
    @public_key,
    @private_key,
    @raw_device_config,
    NOW(),
    NOW(),
    'online'
)
ON CONFLICT (mesh_id, hardware_id)
DO UPDATE SET
    node_num = EXCLUDED.node_num,
    device_id = EXCLUDED.device_id,
    long_name = EXCLUDED.long_name,
    short_name = EXCLUDED.short_name,
    firmware_version = EXCLUDED.firmware_version,
    hw_model = EXCLUDED.hw_model,
    public_key = EXCLUDED.public_key,
    private_key = EXCLUDED.private_key,
    raw_device_config = EXCLUDED.raw_device_config,
    config_imported_at = NOW(),
    last_seen = NOW(),
    status = 'online',
    updated_at = NOW()
RETURNING *;

-- name: UpdateNodeConfigOverrides :one
-- Update node-specific config overrides
UPDATE nodes
SET
    config_overrides = @config_overrides,
    pending_changes = TRUE,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: GetNodeEffectiveConfig :one
-- Get the effective config for a node (merging mesh defaults with node overrides)
-- This is a JSON merge operation
SELECT
    n.id,
    n.hardware_id,
    n.node_num,
    n.device_id,
    n.name,
    n.long_name,
    n.short_name,
    m.config_defaults || n.config_overrides as effective_config,
    n.raw_device_config,
    n.config_imported_at,
    n.config_applied_at,
    n.pending_changes
FROM nodes n
JOIN meshes m ON n.mesh_id = m.id
WHERE n.id = @id;
