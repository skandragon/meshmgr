-- name: AssignAdminKeyToNode :one
INSERT INTO node_admin_keys (node_id, admin_key_id)
VALUES (@node_id, @admin_key_id)
ON CONFLICT (node_id, admin_key_id) DO UPDATE
SET applied_at = NOW(), is_current = TRUE
RETURNING *;

-- name: ListAdminKeysForNode :many
SELECT nak.*, ak.key_name, ak.public_key
FROM node_admin_keys nak
JOIN admin_keys ak ON nak.admin_key_id = ak.id
WHERE nak.node_id = @node_id
ORDER BY nak.applied_at DESC;

-- name: ListNodesForAdminKey :many
SELECT nak.*, n.name as node_name, n.hardware_id
FROM node_admin_keys nak
JOIN nodes n ON nak.node_id = n.id
WHERE nak.admin_key_id = @admin_key_id
ORDER BY nak.applied_at DESC;

-- name: MarkAdminKeyNotCurrent :exec
UPDATE node_admin_keys
SET is_current = FALSE
WHERE node_id = @node_id AND admin_key_id = @admin_key_id;

-- name: GetCurrentAdminKeysForNode :many
SELECT nak.*, ak.key_name, ak.public_key
FROM node_admin_keys nak
JOIN admin_keys ak ON nak.admin_key_id = ak.id
WHERE nak.node_id = @node_id AND nak.is_current = TRUE
ORDER BY nak.applied_at DESC;

-- name: DeleteNodeAdminKeyMapping :exec
DELETE FROM node_admin_keys
WHERE node_id = @node_id AND admin_key_id = @admin_key_id;
