-- File: 000009_create_listings_table.down.sql

-- This migration script drops the 'listings' table if it exists.

-- Drop Foreign Key Constraints
ALTER TABLE "listings"
DROP CONSTRAINT IF EXISTS fk_user_id;

ALTER TABLE "listings"
DROP CONSTRAINT IF EXISTS fk_area_id;

ALTER TABLE "listings"
DROP CONSTRAINT IF EXISTS fk_region_id;

-- Drop Listings Table
DROP TABLE IF EXISTS "listings";