-- +goose up
ALTER TABLE recipes
ADD COLUMN created_by TEXT NOT NULL DEFAULT 'necronomnomicon';

-- +goose down
ALTER TABLE recipes
DROP COLUMN created_by;
