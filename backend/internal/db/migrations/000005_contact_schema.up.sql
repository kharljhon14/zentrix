CREATE TABLE IF NOT EXISTS "contacts" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "company_id" UUID,
    "title" VARCHAR(255) NOT NULL,
    "status" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT fk_company_id 
        FOREIGN KEY ("company_id") REFERENCES "companies"(id)
);

CREATE INDEX idx_contacts_company_id ON contacts(company_id);
