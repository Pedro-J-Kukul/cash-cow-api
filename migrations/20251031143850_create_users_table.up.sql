CREATE TABLE IF NOT EXISTS "users" (
    "id" BIGSERIAL PRIMARY KEY,
    -- optional farmer_id for verification purposes
    "farmer_id" TEXT UNIQUE,
    -- user info
    "email" TEXT UNIQUE NOT NULL,
    "phone_number" TEXT UNIQUE,
    "password_hash" BYTEA NOT NULL,
    -- name
    "first_name" TEXT NOT NULL,
    "last_name" TEXT NOT NULL,
    "middle_name" TEXT,
    -- system related
    "is_activated" BOOLEAN NOT NULL DEFAULT false,
    "is_deleted" BOOLEAN NOT NULL DEFAULT false,
    "is_verified" BOOLEAN NOT NULL DEFAULT false,
    "version" INTEGER NOT NULL DEFAULT 1,
    -- dates
    "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);