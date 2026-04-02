CREATE TABLE budget.budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL REFERENCES identity.organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    month SMALLINT NOT NULL,
    year SMALLINT NOT NULL,
    currency_code VARCHAR(3) NOT NULL REFERENCES budget.currencies(code) ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT budgets_unique_month_year UNIQUE (organization_id, month, year)
);

CREATE INDEX budgets_organization_id_idx
    ON budget.budgets (organization_id);

ALTER TABLE budget.budgets ENABLE ROW LEVEL SECURITY;

CREATE POLICY budgets_org_scope ON budget.budgets
    FOR ALL USING (
        organization_id = current_setting('app.current_organization_id', true)
    );
