-- Add device identification and metadata fields to nodes table
-- These are read-only fields populated from device queries
ALTER TABLE nodes ADD COLUMN node_num BIGINT;
ALTER TABLE nodes ADD COLUMN device_id BYTEA; -- MAC address
ALTER TABLE nodes ADD COLUMN firmware_version TEXT;
ALTER TABLE nodes ADD COLUMN hw_model INTEGER;
ALTER TABLE nodes ADD COLUMN short_name TEXT;

-- Add JSONB column for storing full raw config from device
-- This preserves the complete config structure for reference/debugging
ALTER TABLE nodes ADD COLUMN raw_device_config JSONB;

-- Add JSONB column for node-specific config overrides
-- Only stores values that differ from mesh defaults
ALTER TABLE nodes ADD COLUMN config_overrides JSONB DEFAULT '{}'::jsonb NOT NULL;

-- Add timestamp for last successful config import
ALTER TABLE nodes ADD COLUMN config_imported_at TIMESTAMPTZ;

-- Create indexes for node_num and device_id lookups
CREATE INDEX idx_nodes_node_num ON nodes(node_num);
CREATE INDEX idx_nodes_device_id ON nodes(device_id);

-- Create mesh_channels table for storing all 8 channels
-- Channels are mesh-wide and shared across all nodes
CREATE TABLE mesh_channels (
    id BIGSERIAL PRIMARY KEY,
    mesh_id BIGINT NOT NULL REFERENCES meshes(id) ON DELETE CASCADE,
    channel_index INTEGER NOT NULL CHECK (channel_index >= 0 AND channel_index <= 7),
    channel_role TEXT NOT NULL CHECK (channel_role IN ('PRIMARY', 'SECONDARY', 'DISABLED')),
    psk BYTEA, -- Pre-shared key for the channel
    channel_name TEXT,
    -- Channel settings as JSONB for flexibility
    settings JSONB DEFAULT '{}'::jsonb NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(mesh_id, channel_index)
);

CREATE INDEX idx_mesh_channels_mesh_id ON mesh_channels(mesh_id);

-- Expand meshes table with additional LoRa config that's currently missing
ALTER TABLE meshes ADD COLUMN hop_limit INTEGER DEFAULT 3;
ALTER TABLE meshes ADD COLUMN tx_power INTEGER DEFAULT 30;
ALTER TABLE meshes ADD COLUMN channel_num INTEGER DEFAULT 0;
ALTER TABLE meshes ADD COLUMN use_preset BOOLEAN DEFAULT TRUE NOT NULL;

-- Add JSONB column for mesh-wide default configs
-- These are the default values for all nodes, which can be overridden per-node
ALTER TABLE meshes ADD COLUMN config_defaults JSONB DEFAULT '{}'::jsonb NOT NULL;

-- Alter admin_keys to store the full key data
-- Change public_key from TEXT to BYTEA for binary storage
ALTER TABLE admin_keys ALTER COLUMN public_key TYPE BYTEA USING decode(public_key, 'base64');

-- Add constraints for new mesh config fields
ALTER TABLE meshes ADD CONSTRAINT check_hop_limit
    CHECK (hop_limit >= 0 AND hop_limit <= 7);

ALTER TABLE meshes ADD CONSTRAINT check_tx_power
    CHECK (tx_power >= 0 AND tx_power <= 30);

ALTER TABLE meshes ADD CONSTRAINT check_channel_num
    CHECK (channel_num >= 0 AND channel_num <= 255);
