# meshmgr Design Document

**Version:** 1.0
**Last Updated:** 2025-10-15
**Status:** Draft - MVP Definition

---

## Table of Contents

1. [Vision & Goals](#vision--goals)
2. [Core Concepts](#core-concepts)
3. [Architecture Overview](#architecture-overview)
4. [MVP Features](#mvp-features)
5. [Database Design](#database-design)
6. [API Design](#api-design)
7. [Frontend Structure](#frontend-structure)
8. [Security Model](#security-model)
9. [Implementation Phases](#implementation-phases)
10. [Future Features](#future-features)

---

## Vision & Goals

### Vision

Create a web-based, multi-tenant Meshtastic node management system that enables users to organize, configure, and monitor their mesh networks through an intuitive interface.

### Goals

- **Multi-tenant**: Support multiple users managing multiple independent mesh networks
- **Flexible Configuration**: Configure nodes via serial port (local) or radio (remote with admin keys)
- **Collaborative**: Enable users to share mesh management with others
- **Scalable**: Database-driven architecture that can grow with user needs
- **Modern Stack**: Go backend + SvelteKit 5 frontend for maintainability and performance

### Non-Goals (for MVP)

- Real-time mesh topology visualization
- Advanced analytics and metrics
- Mobile app (web-first approach)
- Multi-region deployment
- Message routing/relay functionality

---

## Core Concepts

### User

A person with an account in the system. Users can own and be granted access to multiple meshes.

### Mesh

A logical grouping of Meshtastic nodes that represents a mesh network. Each mesh:

- Has one or more owners (the user who created it, and others can be granted access later)
- Can have multiple collaborators with varying permission levels
- Contains zero or more nodes
- Has configuration settings (e.g., default channel configs)

### Node

A Meshtastic device that belongs to a mesh. Nodes have:

- Unique hardware ID (from Meshtastic device)
- Public and private key pair (base64-encoded, stored in database)
- Configuration settings (name, role, position, etc.)
- Connection method: serial (local) or radio (remote)
- Current status (online, offline, last seen)

### Admin Key

Public key that nodes are configured to trust for remote administration over the mesh network. Up to 3 admin keys can be specified per mesh. The corresponding private key must be held by a node to perform administrative operations. Admin keys are stored as base64-encoded TEXT.

### Access Level

Permissions granted to a user for a specific mesh:

- **Owner**: Full control, can delete mesh, manage all access, add/remove other owners (enables mesh transfer)
- **Admin**: Can configure nodes, add/remove nodes, view all settings
- **Viewer**: Read-only access to mesh and node status

---

## Architecture Overview

```text
┌─────────────────────────────────────────────────────────────┐
│                    SvelteKit 5 Frontend                      │
│  (Authentication, Mesh Management, Node Configuration UI)   │
└─────────────────┬───────────────────────────────────────────┘
                  │ HTTPS/WebSocket
                  │
┌─────────────────▼───────────────────────────────────────────┐
│                      Go Backend API                          │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │   Auth      │  │   Mesh       │  │   Node Config    │  │
│  │   Service   │  │   Service    │  │   Service        │  │
│  └─────────────┘  └──────────────┘  └──────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │          Device Communication Layer                   │  │
│  │  ┌─────────────────┐    ┌─────────────────┐         │  │
│  │  │  Serial Driver  │    │  Radio Driver   │         │  │
│  │  └─────────────────┘    └─────────────────┘         │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────┬───────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────┐
│                    PostgreSQL Database                       │
│              (SQLC queries, gomigrate migrations)            │
└─────────────────────────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────┐
│                   Physical Devices                           │
│  ┌──────────────┐              ┌──────────────┐            │
│  │  USB Serial  │              │ Radio Device │            │
│  │  Local Node  │              │ (Gateway)    │            │
│  └──────────────┘              └──────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

### Technology Stack

**Backend:**

- Language: Go 1.25+
- Database: PostgreSQL 17+
- SQL: SQLC for type-safe queries
- Migrations: golang-migrate/migrate
- HTTP Framework: Chi router (or similar)
- WebSocket: gorilla/websocket for real-time updates

**Frontend:**

- Framework: SvelteKit 5
- Language: TypeScript
- Styling: TailwindCSS (or similar)
- Auth: JWT based

**Database:**

- PostgreSQL with pgx5 driver
- Structure: `meshdb/queries/*.sql` for SQLC
- Migrations: `meshdb/migrations/[timestamp]_[name].up.sql` and `.down.sql`

---

## MVP Features

### Phase 1: Core Infrastructure (Week 1-2)

- [ ] User authentication (register, login, logout)
- [ ] Session management
- [ ] Database schema for users, meshes, nodes
- [ ] SQLC configuration and initial queries
- [ ] Basic API structure (Chi router setup)
- [ ] Frontend skeleton with SvelteKit 5

### Phase 2: Mesh Management (Week 2-3)

- [ ] Create/read/update/delete meshes
- [ ] Mesh ownership
- [ ] View list of user's meshes
- [ ] Basic mesh detail page

### Phase 3: Node Management (Week 3-4)

- [ ] Add nodes to mesh (manual entry for MVP)
- [ ] View node list in a mesh
- [ ] Node detail page
- [ ] Remove nodes from mesh

### Phase 4: Device Communication (Week 4-6)

- [ ] Serial port enumeration and connection
- [ ] Basic Meshtastic protocol implementation (read node info)
- [ ] Configure node via serial port
- [ ] Display connection status

### Phase 5: Multi-User Access (Week 6-7)

- [ ] Invite users to mesh
- [ ] Grant/revoke access levels
- [ ] View collaborators on a mesh
- [ ] Permission enforcement in API

---

## Database Design

### Schema Overview

```sql
-- Core Tables
users              -- User accounts
meshes             -- Mesh networks
nodes              -- Meshtastic devices
mesh_access        -- User permissions for meshes
admin_keys         -- Admin keys for remote configuration
node_config        -- Node configuration history

-- Supporting Tables
sessions           -- User sessions
audit_log          -- Track changes for security
```

### Core Tables Detail

#### `users`

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    display_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### `meshes`

```sql
CREATE TABLE meshes (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Original creator
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Note:** The `owner_id` field tracks the original creator. Current ownership permissions are managed through the `mesh_access` table, which can grant 'owner' access to multiple users.

#### `nodes`

```sql
CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    mesh_id BIGINT NOT NULL REFERENCES meshes(id) ON DELETE CASCADE,
    hardware_id TEXT NOT NULL, -- Meshtastic device ID
    name TEXT NOT NULL,
    long_name TEXT NOT NULL,
    role TEXT, -- router, client, etc.
    public_key TEXT, -- Base64-encoded public key
    private_key TEXT, -- Base64-encoded private key (stored plaintext for now)
    last_seen TIMESTAMPTZ,
    status TEXT, -- 'online', 'offline', 'unknown'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(mesh_id, hardware_id)
);
```

#### `mesh_access`

```sql
CREATE TABLE mesh_access (
    id BIGSERIAL PRIMARY KEY,
    mesh_id BIGINT NOT NULL REFERENCES meshes(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_level TEXT NOT NULL, -- 'owner', 'admin', 'viewer'
    granted_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(mesh_id, user_id)
);
```

**Note:** Multiple users can have 'owner' access level for a mesh. Owners can add/remove other owners, enabling mesh transfer and co-ownership.

#### `admin_keys`

```sql
CREATE TABLE admin_keys (
    id BIGSERIAL PRIMARY KEY,
    mesh_id BIGINT NOT NULL REFERENCES meshes(id) ON DELETE CASCADE,
    key_name TEXT NOT NULL,
    public_key TEXT NOT NULL, -- Base64-encoded public admin key (up to 3 per mesh)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL REFERENCES users(id)
);
```

#### `sessions`

```sql
CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Directory Structure

```text
meshdb/
├── migrations/
│   ├── 1729000000_initial_schema.up.sql
│   ├── 1729000000_initial_schema.down.sql
│   ├── 1729000001_add_nodes.up.sql
│   └── 1729000001_add_nodes.down.sql
├── queries/
│   ├── users.sql
│   ├── meshes.sql
│   ├── nodes.sql
│   ├── mesh_access.sql
│   └── sessions.sql
└── sqlc.yaml
```

---

## API Design

### Authentication Endpoints

```text
POST   /api/auth/register          - Create new user account
POST   /api/auth/login             - Login and create session
POST   /api/auth/logout            - Logout and destroy session
GET    /api/auth/me                - Get current user info
```

### Mesh Endpoints

```text
GET    /api/meshes                 - List user's meshes
POST   /api/meshes                 - Create new mesh
GET    /api/meshes/:id             - Get mesh details
PUT    /api/meshes/:id             - Update mesh
DELETE /api/meshes/:id             - Delete mesh (owner only)

GET    /api/meshes/:id/access      - List mesh collaborators
POST   /api/meshes/:id/access      - Grant access to user
DELETE /api/meshes/:id/access/:uid - Revoke access
```

### Node Endpoints

```text
GET    /api/meshes/:id/nodes       - List nodes in mesh
POST   /api/meshes/:id/nodes       - Add node to mesh
GET    /api/nodes/:id              - Get node details
PUT    /api/nodes/:id              - Update node configuration
DELETE /api/nodes/:id              - Remove node from mesh

POST   /api/nodes/:id/connect      - Connect to node (serial/radio)
POST   /api/nodes/:id/disconnect   - Disconnect from node
GET    /api/nodes/:id/status       - Get real-time node status
POST   /api/nodes/:id/configure    - Send configuration to node
```

### Device Endpoints

```text
POST   /api/devices/scan           - Scan for nodes on radio
```

### Response Format

All API responses follow this structure:

```json
{
    "success": true,
    "data": { ... },
    "error": null
}
```

Error response:

```json
{
    "success": false,
    "data": null,
    "error": {
        "code": "UNAUTHORIZED",
        "message": "Invalid credentials"
    }
}
```

---

## Frontend Structure

### SvelteKit 5 Application

```text
src/
├── routes/
│   ├── +layout.svelte              - Root layout with auth
│   ├── +page.svelte                - Landing/dashboard
│   ├── auth/
│   │   ├── login/+page.svelte
│   │   └── register/+page.svelte
│   ├── meshes/
│   │   ├── +page.svelte            - Mesh list
│   │   ├── [id]/
│   │   │   ├── +page.svelte        - Mesh detail
│   │   │   ├── nodes/+page.svelte  - Node list
│   │   │   ├── access/+page.svelte - Collaborator management
│   │   │   └── settings/+page.svelte
│   │   └── new/+page.svelte
│   └── nodes/
│       └── [id]/
│           ├── +page.svelte        - Node detail
│           └── configure/+page.svelte
├── lib/
│   ├── components/
│   │   ├── MeshCard.svelte
│   │   ├── NodeCard.svelte
│   │   ├── NodeConfigForm.svelte
│   │   └── AccessControlList.svelte
│   ├── stores/
│   │   ├── auth.ts                 - Auth state management
│   │   ├── meshes.ts               - Mesh data store
│   │   └── nodes.ts                - Node data store
│   └── api/
│       ├── client.ts               - API client wrapper
│       ├── auth.ts                 - Auth API calls
│       ├── meshes.ts               - Mesh API calls
│       └── nodes.ts                - Node API calls
└── app.html
```

### Key Pages

1. **Dashboard** (`/`) - Overview of user's meshes
2. **Mesh Detail** (`/meshes/[id]`) - Single mesh with node list
3. **Node Detail** (`/nodes/[id]`) - Node info and configuration
4. **Access Management** (`/meshes/[id]/access`) - Manage collaborators

---

## Security Model

### Authentication

- Bcrypt password hashing (cost 12)
- Session-based auth with httpOnly cookies
- 7-day session expiration (configurable)
- CSRF protection on state-changing endpoints

### Authorization

- Mesh-level permissions enforced at API layer
- Owner can: delete mesh, manage all access (including adding/removing other owners), all admin actions
- Admin can: configure nodes, add/remove nodes
- Viewer can: view mesh and node status only

### Data Protection

- Admin keys are public keys (no encryption needed)
- Node private keys stored in plaintext for MVP (encryption planned for future)
- Database connection uses TLS
- API served over HTTPS only in production

### Input Validation

- All API inputs validated before database operations
- SQL injection prevented by SQLC parameterized queries
- XSS protection via SvelteKit's automatic escaping

---

## Implementation Phases

### Phase 1: Foundation (Days 1-7)

**Goal:** Working authentication and basic mesh CRUD

1. Database setup
   - Create initial migration with users, meshes, sessions tables
   - Configure sqlc.yaml
   - Write basic SQLC queries for auth and meshes

2. Backend API
   - Chi router setup
   - Auth middleware
   - User registration/login endpoints
   - Mesh CRUD endpoints

3. Frontend
   - SvelteKit 5 project setup
   - Auth pages (login/register)
   - Dashboard with mesh list
   - Create/edit mesh forms

**Success Criteria:**

- User can register, login, logout
- User can create and view meshes
- Sessions persist across page reloads

---

### Phase 2: Node Management (Days 8-14)

**Goal:** Add, view, and organize nodes within meshes

1. Database
   - Add nodes table migration
   - SQLC queries for node CRUD

2. Backend
   - Node CRUD endpoints
   - Permission checks for node operations

3. Frontend
   - Node list view in mesh detail
   - Add node form (manual entry)
   - Node detail page

**Success Criteria:**

- User can add nodes to their mesh
- Nodes display with basic info
- User can edit and remove nodes

---

### Phase 3: Device Communication (Days 15-28)

**Goal:** Connect to real Meshtastic devices

1. Backend
   - Serial port library integration (go.bug.st/serial)
   - Meshtastic protobuf protocol implementation
   - Serial device discovery
   - Basic read/write operations

2. Frontend
   - Serial port selection UI
   - Connect/disconnect controls
   - Real-time status display

3. Testing
   - Test with real Meshtastic device via USB

**Success Criteria:**

- List available serial ports
- Connect to device and read node info
- Display connection status in UI

---

### Phase 4: Multi-User Access (Days 29-35)

**Goal:** Share meshes with collaborators

1. Database
   - Add mesh_access table migration
   - SQLC queries for access control

2. Backend
   - Access grant/revoke endpoints
   - Permission enforcement in existing endpoints

3. Frontend
   - Collaborator management page
   - Invite user by email
   - Display access level

**Success Criteria:**

- Owner can invite users to mesh
- Invited users see shared mesh
- Permissions enforced (viewer can't edit)

---

## Future Features

### Short-term (Post-MVP)

- [ ] Radio-based remote node configuration
- [ ] Admin key management UI
- [ ] Node position/location display on map
- [ ] Bulk node operations
- [ ] Export/import mesh configuration
- [ ] Node configuration templates
- [ ] Real-time WebSocket updates for node status

### Medium-term

- [ ] Mesh topology visualization
- [ ] Node health monitoring and alerts
- [ ] Configuration history and rollback
- [ ] Channel configuration management
- [ ] Message log viewer (if enabled on nodes)
- [ ] API keys for programmatic access

### Long-term

- [ ] Mobile app (iOS/Android)
- [ ] Advanced analytics and reporting
- [ ] Mesh performance metrics
- [ ] Integration with external services (MQTT, etc.)
- [ ] Multi-region deployment support
- [ ] Self-hosted option with single-tenant mode

---

## Open Questions

1. **Node Discovery:** Should we support automatic node discovery via radio scan, or manual entry only for MVP?
   - **Decision:** Manual entry for MVP, auto-discovery in Phase 3+

1. **Real-time Updates:** WebSocket for node status, or polling?
   - **Decision:** Polling for MVP, WebSocket in post-MVP

1. **Serial Port Access:** How to handle serial port access in multi-user deployment?
   - **Decision:** Backend server has direct access, exposes API to frontend

---

## Change Log

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-10-15 | 1.0 | Initial draft - MVP definition | Claude |

---

## References

- [Meshtastic Protocol Documentation](https://meshtastic.org/)
- [SQLC Documentation](https://docs.sqlc.dev/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [SvelteKit Documentation](https://kit.svelte.dev/)
