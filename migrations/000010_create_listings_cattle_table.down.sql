-- File: 000010_create_listings_cattle_table.down.sql

-- This migration script drops the 'listings_cattle' table if it exists.

-- Drop Foreign Key Constraints
ALTER TABLE "listings_cattle"
DROP CONSTRAINT IF EXISTS fk_listing_id;
ALTER TABLE "listings_cattle"
DROP CONSTRAINT IF EXISTS fk_cattle_id;

DROP TABLE IF EXISTS listings_cattle;