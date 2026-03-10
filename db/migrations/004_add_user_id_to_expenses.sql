-- +goose Up
ALTER TABLE expenses ADD COLUMN user_id BIGINT;

-- Backfill: assign existing expenses to the first user (if any).
UPDATE expenses SET user_id = (SELECT id FROM users ORDER BY id LIMIT 1) WHERE user_id IS NULL;

-- If no users exist yet, delete orphan expenses so NOT NULL can be applied.
DELETE FROM expenses WHERE user_id IS NULL;

ALTER TABLE expenses ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE expenses ADD CONSTRAINT fk_expenses_user FOREIGN KEY (user_id) REFERENCES users(id);
CREATE INDEX idx_expenses_user_id ON expenses(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_expenses_user_id;
ALTER TABLE expenses DROP CONSTRAINT IF EXISTS fk_expenses_user;
ALTER TABLE expenses DROP COLUMN IF EXISTS user_id;
