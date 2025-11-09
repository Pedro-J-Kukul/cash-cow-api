-- File: 000006_create_cattle_table.down.sql

-- This migration script drops the 'cattle' table.



-- Drop Foreign Key Constraints
ALTER TABLE IF EXISTS cattle DROP CONSTRAINT IF EXISTS fk_cattle_breed_id;
ALTER TABLE IF EXISTS cattle DROP CONSTRAINT IF EXISTS fk_cattle_owner_id;


-- Drop the cattle table
DROP TABLE IF EXISTS cattle;
DROP TYPE IF EXISTS cattle_sex;