-- File: 000002_create_tokens_table.up.sql

-- This migration script creates the 'tokens' table.
CREATE TABLE IF NOT EXISTS "tokens" (
    -- Primary Key
    "hash" BYTEA PRIMARY KEY,
    -- Foreign Key to Users Table
    "user_id" BIGINT NOT NULL,
    -- Token Info
    "expires_at" TIMESTAMPTZ NOT NULL,
    "scope" TEXT NOT NULL
);

-- Foreign Key Constraint
ALTER TABLE "tokens"
ADD CONSTRAINT fk_tokens_user_id
FOREIGN KEY ("user_id") REFERENCES "users"("id")
ON DELETE CASCADE;