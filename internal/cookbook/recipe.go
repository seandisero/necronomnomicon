package cookbook

import (
	"strings"
)

type Recipe struct {
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Notes       string       `json:"notes"`
}

type Recipes = []Recipe

type RecipeFormData struct {
	Name             string
	ErrorName        string
	Ingredients      string
	ErrorIngredients string
	Steps            string
	ErrorSteps       string
	Notes            string
	ErrorNotes       string
}

func NewRecipe(name string, ingredients []Ingredient, steps []string, notes string) Recipe {
	return Recipe{
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
		Notes:       notes,
	}
}

func CreateRecipe(name, ingredients, steps, notes string) (Recipe, error) {
	recipe := Recipe{}
	recipe.Name = name
	parsedIngredients, err := ParseIngredients(ingredients)
	if err != nil {
		return Recipe{}, err
	}
	parsedSteps := strings.Split(steps, "\n")

	recipe.Ingredients = parsedIngredients
	recipe.Steps = parsedSteps

	return recipe, nil
}

func NewRecipeFormData(name, ingredients, steps, notes string) RecipeFormData {
	return RecipeFormData{
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
		Notes:       notes,
	}
}

func NewRecipeErrorFormData(name, errorName, ingredients, errorIngredients, steps, errorSteps, notes, errorNotes string) RecipeFormData {
	return RecipeFormData{
		Name:             name,
		ErrorName:        errorName,
		Ingredients:      ingredients,
		ErrorIngredients: errorIngredients,
		Steps:            steps,
		ErrorSteps:       errorSteps,
		Notes:            notes,
		ErrorNotes:       errorNotes,
	}
}

func EmptyRecipeFormData() RecipeFormData {
	return RecipeFormData{}
}
