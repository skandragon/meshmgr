-- Remove node admin keys mapping table
DROP INDEX IF EXISTS idx_node_admin_keys_current;
DROP INDEX IF EXISTS idx_node_admin_keys_admin_key_id;
DROP INDEX IF EXISTS idx_node_admin_keys_node_id;
DROP TABLE IF EXISTS node_admin_keys;
