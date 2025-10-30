package cookbook

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/seandisero/necronomnomicon/internal/database"
)

func ingredientsToText(ingredients []Ingredient) string {
	ingredientList := make([]string, len(ingredients))
	for i, ingredient := range ingredients {
		parts := make([]string, 3)
		amount := strconv.Itoa(ingredient.Amount)
		parts[0] = amount
		parts[1] = ingredient.Measure
		parts[2] = ingredient.Name
		joined := strings.Join(parts, ":")
		ingredientList[i] = joined
	}
	return strings.Join(ingredientList, "?")
}

func textToIngredients(text string) ([]Ingredient, error) {
	ingredientList := strings.Split(text, "?")
	ingredients := make([]Ingredient, len(ingredientList))
	for i, ingstr := range ingredientList {
		parts := strings.SplitN(ingstr, ":", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf(
				"malformed ingredient in list for ingredient: %s",
				ingstr,
			)
		}

		amount, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("malformed ingredient amount %v", err)
		}
		ingredient := Ingredient{
			Amount:  amount,
			Measure: parts[1],
			Name:    parts[2],
		}
		ingredients[i] = ingredient
	}
	return ingredients, nil
}

func RecipeFromDBRecipe(dbRecipe database.Recipe) (Recipe, error) {
	recipeIngredients, err := textToIngredients(dbRecipe.Ingredients)
	if err != nil {
		return Recipe{}, err
	}
	recipeSteps := strings.Split(dbRecipe.Steps, "?")
	notes := ""
	if dbRecipe.Notes.Valid {
		notes = dbRecipe.Notes.String
	}

	r := Recipe{
		Name:        dbRecipe.Name,
		Ingredients: recipeIngredients,
		Steps:       recipeSteps,
		Notes:       notes,
	}
	return r, nil

}

func (cb *Cookbook) AddRecipeToDB(c echo.Context, recipe Recipe, id int64) (Recipe, error) {
	ingredients := ingredientsToText(recipe.Ingredients)
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

	params := database.CreateRecipeParams{
		Name:        recipe.Name,
		Ingredients: ingredients,
		Steps:       steps,
		Notes:       notes,
		CreatedBy:   user.Username,
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

func (cb *Cookbook) GetAllRecipes(c echo.Context) ([]Recipe, error) {
	dbRecipes, err := cb.DB.GetAllRecipes(c.Request().Context())
	if err != nil {
		return nil, err
	}
	recipes := make([]Recipe, len(dbRecipes))
	for i, dbRecipe := range dbRecipes {
		recipe, err := RecipeFromDBRecipe(dbRecipe)
		if err != nil {
			return nil, err
		}
		recipes[i] = recipe
	}
	cb.Recipes = recipes
	return recipes, nil
}
