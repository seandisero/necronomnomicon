-- +goose up
CREATE TABLE users(
	id INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose down
DROP TABLE users;

