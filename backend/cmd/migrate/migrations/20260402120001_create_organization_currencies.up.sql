CREATE TABLE budget.organization_currencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL REFERENCES identity.organizations(id) ON DELETE CASCADE,
    currency_code VARCHAR(3) NOT NULL REFERENCES budget.currencies(code) ON DELETE RESTRICT,
    is_base BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT organization_currencies_unique UNIQUE (organization_id, currency_code)
);

-- Enforce exactly one base currency per organization
CREATE UNIQUE INDEX organization_currencies_base_uidx
    ON budget.organization_currencies (organization_id) WHERE is_base = true;

CREATE INDEX organization_currencies_organization_id_idx
    ON budget.organization_currencies (organization_id);

ALTER TABLE budget.organization_currencies ENABLE ROW LEVEL SECURITY;

CREATE POLICY organization_currencies_org_scope ON budget.organization_currencies
    FOR ALL USING (
        organization_id = current_setting('app.current_organization_id', true)
    );
