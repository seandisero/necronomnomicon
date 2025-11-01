package cookbook

import (
	"database/sql"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/seandisero/necronomnomicon/internal/database"
)

type Recipe struct {
	ID          int64
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Notes       string       `json:"notes"`
	CreatedBy   string       `json:"created_by"`
	CreatorID   int64        `json:"creator_id"`
}

type Recipes = []Recipe

func (cb *Cookbook) RecipeExists(c echo.Context, recipe Recipe) bool {
	_, err := cb.DB.GetRecipeByName(c.Request().Context(), recipe.Name)
	return err == nil
}

func NewRecipe(name string, ingredients []Ingredient, steps []string, notes string) Recipe {
	return Recipe{
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
		Notes:       notes,
	}
}

func RecipeFromDBRecipe(dbRecipe database.Recipe) (Recipe, error) {
	recipeIngredients, err := ParseIngredientsFromDB(dbRecipe.Ingredients)
	if err != nil {
		return Recipe{}, err
	}
	recipeSteps := strings.Split(dbRecipe.Steps, "?")
	notes := ""
	if dbRecipe.Notes.Valid {
		notes = dbRecipe.Notes.String
	}

	var creatorId int64
	if dbRecipe.CreatorID.Valid {
		creatorId = dbRecipe.CreatorID.Int64
	}

	r := Recipe{
		ID:          dbRecipe.ID,
		CreatedBy:   dbRecipe.CreatedBy,
		CreatorID:   creatorId,
		Name:        dbRecipe.Name,
		Ingredients: recipeIngredients,
		Steps:       recipeSteps,
		Notes:       notes,
	}
	return r, nil

}

func (cb *Cookbook) AddRecipeToDB(c echo.Context, recipe Recipe, id int64) (Recipe, error) {
	ingredients := ComposeIngredientsForDB(recipe.Ingredients)
	steps := strings.Join(recipe.Steps, "?")
	user, err := cb.DB.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return Recipe{}, err
	}

	var notes sql.NullString

	if recipe.Notes != "" {
		notes.Valid = true
		notes.String = recipe.Notes
	}

	userId := sql.NullInt64{
		Int64: id,
		Valid: true,
	}

	params := database.CreateRecipeParams{
		Name:        recipe.Name,
		Ingredients: ingredients,
		Steps:       steps,
		Notes:       notes,
		CreatedBy:   user.Username,
		CreatorID:   userId,
	}

	recipeRow, err := cb.DB.CreateRecipe(c.Request().Context(), params)
	if err != nil {
		return Recipe{}, err
	}
	r, err := RecipeFromDBRecipe(recipeRow)
	if err != nil {
		return Recipe{}, err
	}
	return r, nil
}
