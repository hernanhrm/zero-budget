ALTER TABLE budget.organization_currencies
    ADD COLUMN rate NUMERIC(20, 10) NOT NULL DEFAULT 1;

UPDATE budget.organization_currencies SET rate = 1 WHERE is_base = true;
