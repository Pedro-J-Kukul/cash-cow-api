-- File: 000004_create_users_permissions_table.up.sql

-- This migration script creates the 'users_permissions' table to establish
-- a many-to-many relationship between 'users' and 'permissions'.
CREATE TABLE IF NOT EXISTS "users_permissions" (
    "user_id" BIGINT NOT NULL,
    "permission_id" BIGINT NOT NULL,
    PRIMARY KEY ("user_id", "permission_id")
);

-- Foreign Key Constraints
ALTER TABLE "users_permissions"
ADD CONSTRAINT fk_users_permissions_user_id
FOREIGN KEY ("user_id") REFERENCES "users"("id")
ON DELETE CASCADE;

ALTER TABLE "users_permissions"
ADD CONSTRAINT fk_users_permissions_permission_id
FOREIGN KEY ("permission_id") REFERENCES "permissions"("id")
ON DELETE CASCADE;