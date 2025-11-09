-- File: 000002_create_tokens_table.down.sql

-- This migration script drops the 'tokens' table.

-- Drop Foreign Key Constraint
ALTER TABLE IF EXISTS tokens DROP CONSTRAINT IF EXISTS fk_tokens_user_id;

-- Drop the tokens table
DROP TABLE IF EXISTS tokens;