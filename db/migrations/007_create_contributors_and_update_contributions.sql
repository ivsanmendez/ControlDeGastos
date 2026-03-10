-- +goose Up

-- 1. Create contributors table
CREATE TABLE contributors (
    id            BIGSERIAL      PRIMARY KEY,
    house_number  VARCHAR(20)    NOT NULL UNIQUE,
    name          VARCHAR(200)   NOT NULL,
    phone         VARCHAR(30),
    user_id       BIGINT         NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- 2. Migrate existing contributor data from contributions
INSERT INTO contributors (house_number, name, phone, user_id, created_at, updated_at)
SELECT DISTINCT ON (house_number)
    house_number,
    contributor_name,
    phone,
    user_id,
    NOW(),
    NOW()
FROM contributions
ORDER BY house_number, created_at DESC;

-- 3. Add contributor_id column to contributions
ALTER TABLE contributions ADD COLUMN contributor_id BIGINT;

-- 4. Populate contributor_id from migrated data
UPDATE contributions c
SET contributor_id = ct.id
FROM contributors ct
WHERE c.house_number = ct.house_number;

-- 5. Make contributor_id NOT NULL and add FK
ALTER TABLE contributions ALTER COLUMN contributor_id SET NOT NULL;
ALTER TABLE contributions
    ADD CONSTRAINT fk_contributions_contributor
    FOREIGN KEY (contributor_id) REFERENCES contributors(id);

-- 6. Drop old columns and constraint
ALTER TABLE contributions DROP CONSTRAINT contributions_house_number_month_year_key;
DROP INDEX idx_contributions_house_year;
ALTER TABLE contributions DROP COLUMN house_number;
ALTER TABLE contributions DROP COLUMN contributor_name;
ALTER TABLE contributions DROP COLUMN phone;

-- 7. New unique constraint and index
ALTER TABLE contributions ADD CONSTRAINT uq_contributions_contributor_month_year
    UNIQUE (contributor_id, month, year);
CREATE INDEX idx_contributions_contributor_year ON contributions(contributor_id, year);

-- +goose Down

-- Reverse: add back columns, copy data from JOIN, drop contributor_id, drop contributors table
ALTER TABLE contributions ADD COLUMN house_number VARCHAR(20);
ALTER TABLE contributions ADD COLUMN contributor_name VARCHAR(200);
ALTER TABLE contributions ADD COLUMN phone VARCHAR(30);

UPDATE contributions c
SET house_number = ct.house_number,
    contributor_name = ct.name,
    phone = ct.phone
FROM contributors ct
WHERE c.contributor_id = ct.id;

ALTER TABLE contributions ALTER COLUMN house_number SET NOT NULL;
ALTER TABLE contributions ALTER COLUMN contributor_name SET NOT NULL;

ALTER TABLE contributions DROP CONSTRAINT uq_contributions_contributor_month_year;
DROP INDEX idx_contributions_contributor_year;
ALTER TABLE contributions DROP CONSTRAINT fk_contributions_contributor;
ALTER TABLE contributions DROP COLUMN contributor_id;

ALTER TABLE contributions ADD CONSTRAINT contributions_house_number_month_year_key
    UNIQUE (house_number, month, year);
CREATE INDEX idx_contributions_house_year ON contributions(house_number, year);

DROP TABLE IF EXISTS contributors;
