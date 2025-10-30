-- +goose up
CREATE TABLE user_recipes(
	user_id INTEGER NOT NULL,
	recipe_id INTEGER NOT NULL,
	data TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (recipe_id) REFERENCES recipes(id)
);

-- +goose down
DROP TABLE user_recipes;
