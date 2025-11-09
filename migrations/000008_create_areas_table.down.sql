-- File: 000008_create_areas_table.down.sql

-- This migration script drops the 'area' table and associated enum type if they exist.

-- Drop Foreign Key Constraint
ALTER TABLE "areas"
DROP CONSTRAINT IF EXISTS fk_areas_region_id;

-- This migration script drops the 'areas' table if it exists.
DROP TABLE IF EXISTS "areas";

-- Drop Area Types Enumeration
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'area_type_enum') THEN
        DROP TYPE area_type_enum;
    END IF;
END $$;
