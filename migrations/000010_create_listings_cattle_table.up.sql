-- File: 000010_create_listings_cattle_table.up.sql

-- This migration script creates the 'listings_cattle' table to have a list of cattle associated with each listing.

-- Create Listings_Cattle Table
CREATE TABLE IF NOT EXISTS "listings_cattle" (
    "listing_id" BIGINT NOT NULL, -- Foreign Key to listings table
    "cattle_id" BIGINT NOT NULL,  -- Foreign Key to cattle table
    PRIMARY KEY ("listing_id", "cattle_id")
);

-- Add Foreign Key Constraints
ALTER TABLE "listings_cattle"
ADD CONSTRAINT fk_listing_id
FOREIGN KEY ("listing_id") REFERENCES "listings"("id") ON DELETE CASCADE;

ALTER TABLE "listings_cattle"
ADD CONSTRAINT fk_cattle_id
FOREIGN KEY ("cattle_id") REFERENCES "cattle"("id") ON DELETE CASCADE;