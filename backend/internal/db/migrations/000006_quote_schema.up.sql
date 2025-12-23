CREATE TABLE IF NOT EXISTS "quotes" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "company_id" UUID NOT NULL,
    "total_amount" INTEGER NOT NULL,
    "sales_tax" INTEGER DEFAULT 0,
    "stage" VARCHAR(255),
    "notes" TEXT,
    "prepared_by" UUID NOT NULL,
    "prepared_for" UUID NOT NULL,

    CONSTRAINT fk_prepared_by 
        FOREIGN KEY ("prepared_by") REFERENCES "users"(id),

    CONSTRAINT fk_prepared_for 
        FOREIGN KEY ("prepared_for") REFERENCES "contacts"(id)
    
);

CREATE INDEX idx_quotes_users_id ON users(id);
CREATE INDEX idx_quotes_contacts_id ON contacts(id);
