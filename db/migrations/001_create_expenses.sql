-- +goose Up
CREATE TABLE expenses (
    id          BIGSERIAL PRIMARY KEY,
    description TEXT           NOT NULL,
    amount      NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    category    VARCHAR(50)    NOT NULL,
    date        DATE           NOT NULL,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS expenses;
