-- File: 000001_create_users_table.up.sql

-- This migration script creates the 'users' table.
CREATE TABLE "users" (
    -- Primary key
    "id" BIGSERIAL PRIMARY KEY,
    -- User Login Info
    "email" TEXT NOT NULL UNIQUE,
    "password_hash" TEXT NOT NULL,
    -- User Profile Info
    "first_name" TEXT NOT NULL,
    "last_name" TEXT NOT NULL,
    -- Optional Info
    "farmer_id" INTEGER, -- For BAHA farmer identification
    "phone_number" TEXT,
    -- System Fields
    "is_activated" BOOLEAN NOT NULL DEFAULT FALSE, -- Email activation status
    "is_verified" BOOLEAN NOT NULL DEFAULT FALSE, -- Farmer verification status
    "is_deleted" BOOLEAN NOT NULL DEFAULT FALSE, -- Soft delete flag
    -- Timestamps
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
    