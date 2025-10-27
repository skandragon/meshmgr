-- name: CreateMesh :one
INSERT INTO meshes (owner_id, name, description, lora_region, modem_preset, frequency_slot)
VALUES (@owner_id, @name, @description, @lora_region, @modem_preset, @frequency_slot)
RETURNING *;

-- name: GetMeshByID :one
SELECT * FROM meshes
WHERE id = @id;

-- name: UpdateMesh :one
UPDATE meshes
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    lora_region = COALESCE(sqlc.narg('lora_region'), lora_region),
    modem_preset = COALESCE(sqlc.narg('modem_preset'), modem_preset),
    frequency_slot = COALESCE(sqlc.narg('frequency_slot'), frequency_slot),
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

-- name: UpdateMeshLoRaConfig :one
-- Update LoRa-specific configuration for a mesh
UPDATE meshes
SET
    lora_region = COALESCE(sqlc.narg('lora_region'), lora_region),
    modem_preset = COALESCE(sqlc.narg('modem_preset'), modem_preset),
    frequency_slot = COALESCE(sqlc.narg('frequency_slot'), frequency_slot),
    hop_limit = COALESCE(sqlc.narg('hop_limit'), hop_limit),
    tx_power = COALESCE(sqlc.narg('tx_power'), tx_power),
    channel_num = COALESCE(sqlc.narg('channel_num'), channel_num),
    use_preset = COALESCE(sqlc.narg('use_preset'), use_preset),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: UpdateMeshConfigDefaults :one
-- Update mesh-wide default configuration
UPDATE meshes
SET
    config_defaults = @config_defaults,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: GetMeshWithDefaults :one
-- Get mesh with all config defaults
SELECT
    id,
    owner_id,
    name,
    description,
    lora_region,
    modem_preset,
    frequency_slot,
    hop_limit,
    tx_power,
    channel_num,
    use_preset,
    config_defaults,
    created_at,
    updated_at
FROM meshes
WHERE id = @id;
