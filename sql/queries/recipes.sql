-- name: CreateRecipe :one
INSERT INTO recipes(name, ingredients, steps, notes, created_by, creator_id, created_at, updated_at) 
VALUES (?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
RETURNING *;

-- name: GetAllRecipes :many
SELECT * FROM recipes;

-- name: GetRecipeByName :one
SELECT * FROM recipes
WHERE name = ?;

-- name: GetRecipeByID :one
SELECT * FROM recipes
WHERE id = ?;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = ?;

-- name: EditRecipe :one
UPDATE recipes
SET name = ?, ingredients = ?, steps = ?, notes = ?, updated_at = DATETIME('now')
WHERE id = ?
RETURNING *;
