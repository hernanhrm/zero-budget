CREATE TABLE budget.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL REFERENCES identity.organizations(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES budget.categories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(7),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT categories_no_self_ref CHECK (parent_id IS DISTINCT FROM id),
    CONSTRAINT categories_unique_name UNIQUE (organization_id, parent_id, name)
);

CREATE INDEX categories_organization_id_idx
    ON budget.categories (organization_id);
CREATE INDEX categories_parent_id_idx
    ON budget.categories (parent_id);

ALTER TABLE budget.categories ENABLE ROW LEVEL SECURITY;

CREATE POLICY categories_org_scope ON budget.categories
    FOR ALL USING (
        organization_id = current_setting('app.current_organization_id', true)
    );
