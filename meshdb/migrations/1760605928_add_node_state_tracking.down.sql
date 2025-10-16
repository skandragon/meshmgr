-- Remove node state tracking fields
DROP INDEX IF EXISTS idx_nodes_pending_changes;

ALTER TABLE nodes DROP COLUMN IF EXISTS pending_changes;
ALTER TABLE nodes DROP COLUMN IF EXISTS config_applied_at;
ALTER TABLE nodes DROP COLUMN IF EXISTS unmessageable;
ALTER TABLE nodes DROP COLUMN IF EXISTS applied_unmessageable;
ALTER TABLE nodes DROP COLUMN IF EXISTS applied_private_key;
ALTER TABLE nodes DROP COLUMN IF EXISTS applied_public_key;
ALTER TABLE nodes DROP COLUMN IF EXISTS applied_role;
ALTER TABLE nodes DROP COLUMN IF EXISTS applied_long_name;
ALTER TABLE nodes DROP COLUMN IF EXISTS applied_name;
