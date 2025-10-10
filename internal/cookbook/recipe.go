package cookbook

type Recipe struct {
	Name        string
	Ingredients []Ingredient
	Steps       []string
	Notes       string
}

type RecipeFormData struct {
	Name        string
	Ingredients string
	Steps       string
	Notes       string
}

func NewRecipe(name string, ingredients []Ingredient, steps []string, notes string) Recipe {
	return Recipe{
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
		Notes:       notes,
	}
}

type Recipes = []Recipe

func newRecipeFormData() RecipeFormData {
	return RecipeFormData{}
}
