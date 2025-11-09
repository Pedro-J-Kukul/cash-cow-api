-- File: 000006_create_cattle_table.up.sql

-- Enum for Cattle Sex
CREATE TYPE cattle_sex AS ENUM ('male', 'female', 'unknown');



-- This migration script creates the 'cattle' table.
CREATE TABLE IF NOT EXISTS "cattle" (
    -- Primary Key
    "id" BIGSERIAL PRIMARY KEY,
    -- Foreign Keys
    "owner_id" BIGINT NOT NULL, -- owner of the cattle
    "breed_id" BIGINT NOT NULL, -- breed of the cattle
    -- Cattle Info
    "tag_number" TEXT NOT NULL UNIQUE,
    "sex" cattle_sex NOT NULL,
    "age_months" INT NOT NULL,
    "weight_kg" FLOAT, -- optional incase of unknown weight
    -- Health Info
    "vaccinations" TEXT,
    "medical_history" TEXT,
    "is_pregnant" BOOLEAN, -- for female cattle
    "is_castrated" BOOLEAN, -- for male cattle
    -- System Fields
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE, -- Active status
    -- Timestamps
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Foreign Key Constraint
ALTER TABLE "cattle"
ADD CONSTRAINT fk_cattle_breed_id
FOREIGN KEY ("breed_id") REFERENCES "breeds"("id")
ON DELETE RESTRICT;

ALTER TABLE "cattle"
ADD CONSTRAINT fk_cattle_owner_id
FOREIGN KEY ("owner_id") REFERENCES "users"("id")
ON DELETE CASCADE;