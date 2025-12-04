CREATE TABLE IF NOT EXISTS "users"(
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password_hash" TEXT NOT NULL,
    "avatar" TEXT,
    "activated" BOOL NOT NULL DEFAULT FALSE,
    "role" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT now(),
    "updated_at" TIMESTAMPTZ DEFAULT now()
);