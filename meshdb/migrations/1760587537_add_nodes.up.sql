CREATE TABLE admin_keys (
    id BIGSERIAL PRIMARY KEY,
    mesh_id BIGINT NOT NULL REFERENCES meshes(id) ON DELETE CASCADE,
    public_key TEXT NOT NULL,
    key_name TEXT,
    added_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_keys_mesh_id ON admin_keys(mesh_id);

CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    mesh_id BIGINT NOT NULL REFERENCES meshes(id) ON DELETE CASCADE,
    hardware_id TEXT NOT NULL,
    name TEXT NOT NULL,
    long_name TEXT NOT NULL,
    role TEXT,
    public_key TEXT,
    private_key TEXT,
    last_seen TIMESTAMPTZ,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(mesh_id, hardware_id)
);

CREATE INDEX idx_nodes_mesh_id ON nodes(mesh_id);
CREATE INDEX idx_nodes_hardware_id ON nodes(hardware_id);
CREATE INDEX idx_nodes_status ON nodes(status);
