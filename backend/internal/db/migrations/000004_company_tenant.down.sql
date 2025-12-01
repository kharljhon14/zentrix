ALTER TABLE companies
    DROP CONSTRAINT fk_companies_tenant;

ALTER TABLE companies
    DROP COLUMN tenant_id;