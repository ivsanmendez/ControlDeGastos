-- +goose Up
CREATE TABLE contributions (
    id               BIGSERIAL      PRIMARY KEY,
    house_number     VARCHAR(20)    NOT NULL,
    contributor_name VARCHAR(200)   NOT NULL,
    phone            VARCHAR(30),
    amount           NUMERIC(12,2)  NOT NULL CHECK (amount > 0),
    month            INT            NOT NULL CHECK (month BETWEEN 1 AND 12),
    year             INT            NOT NULL CHECK (year >= 2000),
    payment_date     DATE           NOT NULL,
    payment_method   VARCHAR(20)    NOT NULL CHECK (payment_method IN ('cash', 'transfer', 'other')),
    user_id          BIGINT         NOT NULL REFERENCES users(id),
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    UNIQUE (house_number, month, year)
);

CREATE INDEX idx_contributions_house_year ON contributions(house_number, year);

-- +goose Down
DROP TABLE IF EXISTS contributions;
