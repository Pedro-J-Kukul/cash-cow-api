-- File 000007_create_regions_table.up.sql

-- This migration script creates the 'regions' table.
CREATE TABLE IF NOT EXISTS "regions" (
    -- Primary Key
    "id" BIGSERIAL PRIMARY KEY,
    -- Region Info
    "name" TEXT NOT NULL UNIQUE,
    "code" TEXT NOT NULL UNIQUE -- e.g., BZ, CZL, OW, CY,
);