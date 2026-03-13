-- +goose Up
ALTER TABLE contributions ADD COLUMN updated_at TIMESTAMPTZ;
UPDATE contributions SET updated_at = created_at;
ALTER TABLE contributions ALTER COLUMN updated_at SET NOT NULL;
ALTER TABLE contributions ALTER COLUMN updated_at SET DEFAULT NOW();

-- +goose Down
ALTER TABLE contributions DROP COLUMN IF EXISTS updated_at;
