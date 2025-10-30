-- +goose up
CREATE TABLE recipes(
	id INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
	data TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose down
DROP TABLE recipes;
