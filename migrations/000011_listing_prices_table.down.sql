-- File: 000011_listing_prices_table.down.sql
-- This migration script drops the 'listing_prices' table if it exists.

-- Drop Foreign Key Constraints
ALTER TABLE listing_prices
DROP CONSTRAINT IF EXISTS fk_listing_id;

-- Drop Listing Prices Table
DROP TABLE IF EXISTS listing_prices;

-- Drop Price Types Enumeration
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cattle_class_enum') THEN
        DROP TYPE cattle_class_enum;
    END IF;
END $$;
