-- File: 000005_create_breeds_table.up.sql

-- This migration script creates the 'breeds' table.
CREATE TABLE IF NOT EXISTS "breeds" (
    -- Primary Key
    "id" BIGSERIAL PRIMARY KEY,
    -- Breed Info
    "name" TEXT NOT NULL UNIQUE,
    "description" TEXT,
    -- System Fields
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE, -- Active status
    -- Timestamps
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);