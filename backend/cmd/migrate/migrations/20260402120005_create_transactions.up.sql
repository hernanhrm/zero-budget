CREATE TABLE budget.transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL REFERENCES identity.organizations(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES budget.accounts(id) ON DELETE RESTRICT,
    category_id UUID REFERENCES budget.categories(id) ON DELETE SET NULL,
    subcategory_id UUID REFERENCES budget.categories(id) ON DELETE SET NULL,
    budget_id UUID REFERENCES budget.budgets(id) ON DELETE SET NULL,
    type VARCHAR(20) NOT NULL,
    amount BIGINT NOT NULL,
    description TEXT,
    external_reference_number VARCHAR(255),
    date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX transactions_organization_id_idx
    ON budget.transactions (organization_id);
CREATE INDEX transactions_account_id_idx
    ON budget.transactions (account_id);
CREATE INDEX transactions_category_id_idx
    ON budget.transactions (category_id);
CREATE INDEX transactions_budget_id_idx
    ON budget.transactions (budget_id);
CREATE INDEX transactions_organization_id_date_idx
    ON budget.transactions (organization_id, date);
CREATE INDEX transactions_organization_id_external_ref_idx
    ON budget.transactions (organization_id, external_reference_number);

ALTER TABLE budget.transactions ENABLE ROW LEVEL SECURITY;

CREATE POLICY transactions_org_scope ON budget.transactions
    FOR ALL USING (
        organization_id = current_setting('app.current_organization_id', true)
    );
