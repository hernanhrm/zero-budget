CREATE TABLE budget.accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL REFERENCES identity.organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    currency_code VARCHAR(3) NOT NULL REFERENCES budget.currencies(code) ON DELETE RESTRICT,
    current_balance BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX accounts_organization_id_idx
    ON budget.accounts (organization_id);
CREATE INDEX accounts_organization_id_type_idx
    ON budget.accounts (organization_id, type);

ALTER TABLE budget.accounts ENABLE ROW LEVEL SECURITY;

CREATE POLICY accounts_org_scope ON budget.accounts
    FOR ALL USING (
        organization_id = current_setting('app.current_organization_id', true)
    );
