-- +migrate Down
DROP TABLE IF EXISTS time_entries;
ALTER TABLE tasks DROP COLUMN IF EXISTS estimated_minutes;
