-- File: 000003_create_permissions_table.up.sql
-- This migration script creates the 'permissions' table.
CREATE TABLE IF NOT EXISTS "permissions" (
    -- Primary Key
    "id" BIGSERIAL PRIMARY KEY,
    -- Permission Name
    "code" TEXT NOT NULL UNIQUE -- e.g., 'read:users', 'write:tokens'
);