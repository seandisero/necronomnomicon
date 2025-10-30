-- name: CreateRecipe :one
INSERT INTO recipes(name, ingredients, steps, notes, created_by, created_at, updated_at) 
VALUES (?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
RETURNING *;

-- name: GetAllRecipes :many
SELECT * FROM recipes;

-- name: GetRecipeByName :one
SELECT * FROM recipes
WHERE name = ?;
