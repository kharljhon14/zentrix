CREATE TABLE IF NOT EXISTS "users"(
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "tenant_id" UUID NOT NULL REFERENCES tenants(id),
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password_hash" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT now(),
    "updated_at" TIMESTAMPTZ DEFAULT now()
);