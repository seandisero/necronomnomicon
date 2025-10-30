-- +goose up
ALTER TABLE recipes
DROP COLUMN data;

-- +goose down
ALTER TABLE recipes
ADD data TEXT NOT NULL DEFAULT 'default';
