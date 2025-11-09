-- File: 000011_listing_prices_table.up.sql

-- This migration script holds the unit prices per age bracket of cattle listings.

-- Create Enum Type for Age Bracket
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cattle_class_enum') THEN
        CREATE TYPE cattle_class_enum AS ENUM (
            'male_calf', -- Young male cattle
            'female_calf', -- Young female cattle
            'steer', -- Castrated male
            'heifer', -- Young female cattle that has not borne a calf
            'cow', -- Adult female cattle
            'bull'  -- Adult male cattle
        );
    END IF;
END $$;

-- Create Listing Prices Table
CREATE TABLE IF NOT EXISTS "listing_prices" (
    "listing_id" BIGINT NOT NULL, -- Foreign Key to listings table
    "cattle_class" cattle_class_enum NOT NULL, -- Age Bracket of Cattle
    "price_per_kg" NUMERIC(10, 2) NOT NULL, -- Price per Kilogram
    "quantity" INT NOT NULL, -- Quantity of Cattle in this Bracket
    PRIMARY KEY ("listing_id", "cattle_class")
);

-- Add Foreign Key Constraint
ALTER TABLE "listing_prices"
ADD CONSTRAINT fk_listing_id
FOREIGN KEY ("listing_id") REFERENCES "listings"("id") ON DELETE CASCADE;

