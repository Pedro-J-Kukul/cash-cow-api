-- File: 000008_create_areas_table.up.sql

-- This migration script creates the 'areas' table.

-- Area Types Enumeration
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'area_type_enum') THEN
        CREATE TYPE area_type_enum AS ENUM ('city', 'town', 'village');
    END IF;
END $$;

-- Create Areas Table
CREATE TABLE IF NOT EXISTS "areas" (
    -- Primary Key
    "id" BIGSERIAL PRIMARY KEY,
    -- Area Info
    "region_id" BIGINT NOT NULL,
    "name" TEXT NOT NULL,
    "area_type" area_type_enum NOT NULL, -- e.g., city, town, village
    "latitude" FLOAT,
    "longitude" FLOAT,
    -- System Fields
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE, -- Active status
    -- Timestamps
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Foreign Key Constraint
ALTER TABLE "areas"
ADD CONSTRAINT fk_areas_region_id
FOREIGN KEY ("region_id") REFERENCES "regions"("id")
ON DELETE RESTRICT;