-- +goose up
ALTER TABLE recipes
ADD COLUMN ingredients TEXT NOT NULL;

-- +goose down
ALTER TABLE recipes
DROP COLUMN ingredients;
