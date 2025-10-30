-- name: CreateUserRecipe :exec
INSERT INTO user_recipes(user_id, recipe_id, data) VALUES(?, ?, ?);
