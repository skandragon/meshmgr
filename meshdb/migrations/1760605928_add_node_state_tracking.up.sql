-- Add applied state tracking fields for nodes
-- These track what's currently on the device vs what's desired
ALTER TABLE nodes ADD COLUMN applied_name TEXT;
ALTER TABLE nodes ADD COLUMN applied_long_name TEXT;
ALTER TABLE nodes ADD COLUMN applied_role TEXT;
ALTER TABLE nodes ADD COLUMN applied_public_key TEXT;
ALTER TABLE nodes ADD COLUMN applied_private_key TEXT;
ALTER TABLE nodes ADD COLUMN applied_unmessageable BOOLEAN;

-- Add unmessageable field (desired state)
ALTER TABLE nodes ADD COLUMN unmessageable BOOLEAN DEFAULT FALSE NOT NULL;

-- Add timestamp for when config was last applied to device
ALTER TABLE nodes ADD COLUMN config_applied_at TIMESTAMPTZ;

-- Add pending_changes flag to quickly identify nodes needing updates
ALTER TABLE nodes ADD COLUMN pending_changes BOOLEAN DEFAULT FALSE NOT NULL;

-- Create index for querying nodes with pending changes
CREATE INDEX idx_nodes_pending_changes ON nodes(pending_changes) WHERE pending_changes = TRUE;
