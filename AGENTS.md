# Agent Guidelines and Project Conventions

This file documents important conventions and guidelines for AI agents working on this project.

## Code Licensing

### SQL Files

**Do not add license headers to .sql files.**

SQL files (including migration files in `meshdb/migrations/`) should not include license headers. Keep SQL files minimal and focused on the database changes only.

Example of correct SQL migration file:
```sql
-- Add key_preview column to show first 6 chars of the key for identification
ALTER TABLE user_api_keys ADD COLUMN key_preview TEXT NOT NULL DEFAULT '';
```

### Other Source Files

Go source files, TypeScript/Svelte files, and other code should include the standard AGPL-3.0 license header.
