
ALTER TABLE companies
    ADD COLUMN tenant_id UUID;


ALTER TABLE companies
    ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE companies
    ADD CONSTRAINT fk_companies_tenant
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
    ON DELETE CASCADE;

CREATE INDEX idx_companies_tenant_id ON companies (tenant_id);