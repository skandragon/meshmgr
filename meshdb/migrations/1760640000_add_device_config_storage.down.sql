-- Revert admin_keys public_key to TEXT
ALTER TABLE admin_keys ALTER COLUMN public_key TYPE TEXT USING encode(public_key, 'base64');

-- Drop mesh config fields
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_channel_num;
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_tx_power;
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_hop_limit;

ALTER TABLE meshes DROP COLUMN IF EXISTS config_defaults;
ALTER TABLE meshes DROP COLUMN IF EXISTS use_preset;
ALTER TABLE meshes DROP COLUMN IF EXISTS channel_num;
ALTER TABLE meshes DROP COLUMN IF EXISTS tx_power;
ALTER TABLE meshes DROP COLUMN IF EXISTS hop_limit;

-- Drop mesh_channels table
DROP INDEX IF EXISTS idx_mesh_channels_mesh_id;
DROP TABLE IF EXISTS mesh_channels;

-- Drop node indexes
DROP INDEX IF EXISTS idx_nodes_device_id;
DROP INDEX IF EXISTS idx_nodes_node_num;

-- Drop node config fields
ALTER TABLE nodes DROP COLUMN IF EXISTS config_imported_at;
ALTER TABLE nodes DROP COLUMN IF EXISTS config_overrides;
ALTER TABLE nodes DROP COLUMN IF EXISTS raw_device_config;
ALTER TABLE nodes DROP COLUMN IF EXISTS short_name;
ALTER TABLE nodes DROP COLUMN IF EXISTS hw_model;
ALTER TABLE nodes DROP COLUMN IF EXISTS firmware_version;
ALTER TABLE nodes DROP COLUMN IF EXISTS device_id;
ALTER TABLE nodes DROP COLUMN IF EXISTS node_num;
