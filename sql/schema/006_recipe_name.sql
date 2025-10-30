-- +goose up
ALTER TABLE recipes
ADD COLUMN name TEXT NOT NULL;

-- +goose down
ALTER TABLE recipes
DROP COLUMN name;
