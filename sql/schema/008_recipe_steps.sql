-- +goose up
ALTER TABLE recipes
ADD COLUMN steps TEXT NOT NULL;

-- +goose down
ALTER TABLE recipes
DROP COLUMN steps;
