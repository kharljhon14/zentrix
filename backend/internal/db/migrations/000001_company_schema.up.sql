CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE "companies"(
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    "sales_owner" UUID REFERENCES users(id),
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "company_size" VARCHAR(255) NOT NULL,
    "industry" VARCHAR(255) NOT NULL,
    "business_type" VARCHAR(255) NOT NULL,
    "country" VARCHAR(255) NOT NULL,
    "image" TEXT,
    "website" TEXT,
    "created_at" TIMESTAMPTZ DEFAULT now(),
    "updated_at" TIMESTAMPTZ DEFAULT now()
);

