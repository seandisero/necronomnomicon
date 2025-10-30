-- +goose up
ALTER TABLE recipes
ADD COLUMN notes TEXT;

-- +goose down
ALTER TABLE recipes
DROP COLUMN notes; 
