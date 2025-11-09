-- File: 000009_create_listings_table.up.sql
-- This migration script creates the 'listings' table.
CREATE TABLE IF NOT EXISTS "listings" (
    -- Primary Key
    "id" BIGSERIAL PRIMARY KEY,
    -- Foreign Keys
    "user_id" BIGINT NOT NULL, -- who created the listing
    "area_id" BIGINT NOT NULL, -- area where the listing is located
    "region_id" BIGINT NOT NULL, -- region where the listing is located
    -- Listing Info
    "title" TEXT NOT NULL,
    "description" TEXT,
    "latitude" FLOAT,
    "longitude" FLOAT,
    -- Status Info
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE, -- Active status
    -- Timestamps
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Foreign Key Constraints
ALTER TABLE "listings"
ADD CONSTRAINT fk_listings_user_id
FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "listings"
ADD CONSTRAINT fk_listings_area_id
FOREIGN KEY ("area_id") REFERENCES "areas"("id") ON DELETE CASCADE;
ALTER TABLE "listings"
ADD CONSTRAINT fk_listings_region_id
FOREIGN KEY ("region_id") REFERENCES "regions"("id") ON DELETE CASCADE;