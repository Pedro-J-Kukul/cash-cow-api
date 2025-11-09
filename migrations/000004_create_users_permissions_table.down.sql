-- File: 000004_create_users_permissions_table.down.sql
-- This migration script drops the 'users_permissions' table.

-- Drop Foreign Key Constraints
ALTER TABLE IF EXISTS users_permissions DROP CONSTRAINT IF EXISTS fk_users_permissions_user_id;
ALTER TABLE IF EXISTS users_permissions DROP CONSTRAINT IF EXISTS fk_users_permissions_permission_id;

-- Drop the users_permissions table
DROP TABLE IF EXISTS users_permissions;