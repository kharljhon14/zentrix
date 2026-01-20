CREATE TABLE IF NOT EXISTS projects(
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "company_id" UUID NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "status" VARCHAR(255) NOT NULL,
    "owner_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_company_id
        FOREIGN KEY ("company_id") REFERENCES "companies"(id),

    CONSTRAINT fk_onwer_id
        FOREIGN KEY ("owner_id") REFERENCES "users"(id)
);

CREATE INDEX idx_projects_company_id ON companies(id);
CREATE INDEX idx_projects_owner_id ON users(id);