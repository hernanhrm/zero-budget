CREATE SCHEMA IF NOT EXISTS budget;

-- Global reference table for ISO 4217 currencies (no RLS)
CREATE TABLE budget.currencies (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(5) NOT NULL,
    decimal_places SMALLINT NOT NULL DEFAULT 2
);

-- Seed common currencies
INSERT INTO budget.currencies (code, name, symbol, decimal_places) VALUES
    ('USD', 'US Dollar', '$', 2),
    ('EUR', 'Euro', '€', 2),
    ('GBP', 'British Pound', '£', 2),
    ('COP', 'Colombian Peso', '$', 2),
    ('MXN', 'Mexican Peso', '$', 2),
    ('BRL', 'Brazilian Real', 'R$', 2),
    ('ARS', 'Argentine Peso', '$', 2),
    ('CLP', 'Chilean Peso', '$', 0),
    ('PEN', 'Peruvian Sol', 'S/', 2),
    ('CAD', 'Canadian Dollar', '$', 2),
    ('JPY', 'Japanese Yen', '¥', 0),
    ('CNY', 'Chinese Yuan', '¥', 2);
