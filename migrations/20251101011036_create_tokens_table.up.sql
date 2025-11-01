CREATE TABLE IF NOT EXISTS "tokens" (
  "hash" bytea PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "expiry" timestamp NOT NULL,
  "scope" text NOT NULL
);