-- Create junction table to track which admin keys have been applied to which nodes
-- This is essential for safe key rotation
CREATE TABLE node_admin_keys (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    admin_key_id BIGINT NOT NULL REFERENCES admin_keys(id) ON DELETE CASCADE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_current BOOLEAN DEFAULT TRUE NOT NULL,
    UNIQUE(node_id, admin_key_id)
);

-- Create indexes for efficient queries
CREATE INDEX idx_node_admin_keys_node_id ON node_admin_keys(node_id);
CREATE INDEX idx_node_admin_keys_admin_key_id ON node_admin_keys(admin_key_id);
CREATE INDEX idx_node_admin_keys_current ON node_admin_keys(is_current) WHERE is_current = TRUE;
