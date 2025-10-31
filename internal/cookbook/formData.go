package cookbook

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/seandisero/necronomnomicon/internal/database"
)

type RecipeFormData struct {
	ID          int64
	Name        string
	Ingredients string
	Steps       string
	Notes       string

	ErrorName        string
	ErrorIngredients string
	ErrorSteps       string
	ErrorNotes       string

	IsNew  bool
	IsEdit bool
}

func MakeRecipeFormData(name, ingredients, steps, notes string) RecipeFormData {
	return RecipeFormData{
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
		Notes:       notes,
	}
}

func EmptyRecipeFormData() RecipeFormData {
	// Is new because we will never send down empty data for an edit form
	return RecipeFormData{
		IsNew: true,
	}
}

func formDataFromContext(c echo.Context) RecipeFormData {
	name := c.FormValue("name")
	ingredients := c.FormValue("ingredients")
	steps := c.FormValue("steps")
	notes := c.FormValue("notes")

	formData := MakeRecipeFormData(name, ingredients, steps, notes)
	formData.IsNew = true
	return formData
}

func formDataFromDBRecipe(recipe database.Recipe) (RecipeFormData, error) {
	notes := ""
	if recipe.Notes.Valid {
		notes = recipe.Notes.String
	}
	formIngredients := strings.ReplaceAll(recipe.Ingredients, "?", "\n")
	formIngredients = strings.ReplaceAll(formIngredients, ":", " ")

	data := MakeRecipeFormData(
		recipe.Name,
		formIngredients,
		strings.ReplaceAll(recipe.Steps, "?", "\n"),
		notes,
	)
	err := validateRecipeFormData(data)
	if err != nil {
		slog.Error("error validating recipe from database!", "error", err)
		return EmptyRecipeFormData(), fmt.Errorf("internal server error")
	}
	return data, nil
}

func recipeFromFormData(data RecipeFormData) (Recipe, error) {
	if err := validateRecipeFormData(data); err != nil {
		return Recipe{}, err
	}
	recipe := Recipe{}
	recipe.Name = data.Name
	parsedIngredients, err := ParseIngredientsFromForm(data.Ingredients)
	if err != nil {
		return Recipe{}, err
	}
	parsedSteps := strings.Split(data.Steps, "\n")

	recipe.Ingredients = parsedIngredients
	recipe.Steps = parsedSteps
	recipe.Notes = data.Notes

	return recipe, nil
}

func validateRecipeFormData(data RecipeFormData) error {
	if data.Name == "" {
		return fmt.Errorf("recipe must have a name")
	}
	if data.Ingredients == "" {
		return fmt.Errorf("recipe must have ingredients")
	}
	if data.Steps == "" {
		return fmt.Errorf("recipe must have steps")
	}

	return nil
}
