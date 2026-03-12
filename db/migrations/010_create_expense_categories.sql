-- +goose Up

-- 1. Create expense_categories table
CREATE TABLE expense_categories (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(500),
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 2. Seed the 4 legacy categories (using first user as owner)
INSERT INTO expense_categories (name, description, is_active, user_id)
SELECT name, description, TRUE, (SELECT id FROM users ORDER BY id LIMIT 1)
FROM (VALUES
    ('Comida', 'Gastos de comida'),
    ('Transporte', 'Gastos de transporte'),
    ('Vivienda', 'Gastos de vivienda'),
    ('Otro', 'Gastos varios')
) AS v(name, description);

-- 3. Add category_id column to expenses
ALTER TABLE expenses ADD COLUMN category_id BIGINT;

-- 4. Backfill category_id from the old category string
UPDATE expenses SET category_id = (SELECT id FROM expense_categories WHERE name = 'Comida') WHERE category = 'food';
UPDATE expenses SET category_id = (SELECT id FROM expense_categories WHERE name = 'Transporte') WHERE category = 'transport';
UPDATE expenses SET category_id = (SELECT id FROM expense_categories WHERE name = 'Vivienda') WHERE category = 'housing';
UPDATE expenses SET category_id = (SELECT id FROM expense_categories WHERE name = 'Otro') WHERE category = 'other';

-- Fallback: any unmatched rows get "Otro"
UPDATE expenses SET category_id = (SELECT id FROM expense_categories WHERE name = 'Otro') WHERE category_id IS NULL;

-- 5. Make NOT NULL + FK
ALTER TABLE expenses ALTER COLUMN category_id SET NOT NULL;
ALTER TABLE expenses ADD CONSTRAINT fk_expenses_category FOREIGN KEY (category_id) REFERENCES expense_categories(id);

-- 6. Drop old category column
ALTER TABLE expenses DROP COLUMN category;

-- 7. Index
CREATE INDEX idx_expenses_category ON expenses(category_id);

-- +goose Down
ALTER TABLE expenses ADD COLUMN category VARCHAR(50);
UPDATE expenses SET category = 'other';
ALTER TABLE expenses ALTER COLUMN category SET NOT NULL;
DROP INDEX IF EXISTS idx_expenses_category;
ALTER TABLE expenses DROP CONSTRAINT IF EXISTS fk_expenses_category;
ALTER TABLE expenses DROP COLUMN IF EXISTS category_id;
DROP TABLE IF EXISTS expense_categories;
