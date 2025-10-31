-- +goose up
ALTER TABLE recipes
ADD COLUMN creator_id INTEGER
REFERENCES users(id);

-- +goose down
ALTER TABLE recipes
DROP COLUMN creator_id;
