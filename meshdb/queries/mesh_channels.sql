-- name: UpsertMeshChannel :one
-- Insert or update a mesh channel
INSERT INTO mesh_channels (
    mesh_id,
    channel_index,
    channel_role,
    psk,
    channel_name,
    settings
) VALUES (
    @mesh_id,
    @channel_index,
    @channel_role,
    @psk,
    @channel_name,
    @settings
)
ON CONFLICT (mesh_id, channel_index)
DO UPDATE SET
    channel_role = EXCLUDED.channel_role,
    psk = EXCLUDED.psk,
    channel_name = EXCLUDED.channel_name,
    settings = EXCLUDED.settings,
    updated_at = NOW()
RETURNING *;

-- name: GetMeshChannel :one
SELECT * FROM mesh_channels
WHERE mesh_id = @mesh_id AND channel_index = @channel_index;

-- name: ListMeshChannels :many
SELECT * FROM mesh_channels
WHERE mesh_id = @mesh_id
ORDER BY channel_index ASC;

-- name: DeleteMeshChannel :exec
DELETE FROM mesh_channels
WHERE mesh_id = @mesh_id AND channel_index = @channel_index;

-- name: ImportMeshChannels :exec
-- Import all 8 channels from device config, replacing existing ones
-- This should be called within a transaction
DELETE FROM mesh_channels WHERE mesh_id = @mesh_id;

-- name: GetPrimaryChannel :one
-- Get the primary (channel 0) for a mesh
SELECT * FROM mesh_channels
WHERE mesh_id = @mesh_id AND channel_index = 0;

-- name: CountMeshChannels :one
SELECT COUNT(*) FROM mesh_channels
WHERE mesh_id = @mesh_id;
