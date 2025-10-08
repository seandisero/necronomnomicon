package cookbook

import (
	"fmt"
	"strconv"
	"strings"
)

type Ingredient struct {
	Amount  int
	Measure string
	Name    string
}

func NewIngredient(name string, amount int, measure string) Ingredient {
	return Ingredient{
		Name:    name,
		Amount:  amount,
		Measure: measure,
	}
}

type Recipe struct {
	Name        string
	Ingredients []Ingredient
	Steps       []string
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

type Recipies = []Recipe

func ParseIngredients(ingredients string) ([]Ingredient, error) {
	ret := make([]Ingredient, 0)
	ingredientsList := strings.Split(ingredients, "\n")

	for _, ingredient := range ingredientsList {
		ingredientSplit := strings.SplitN(ingredient, " ", 3)
		if len(ingredientSplit) != 3 {
			return []Ingredient{}, fmt.Errorf("malformed ingredient")
		}
		intAmount, err := strconv.Atoi(ingredientSplit[0])
		if err != nil {
			return []Ingredient{}, err
		}
		ret = append(ret, NewIngredient(ingredientSplit[2], intAmount, ingredientSplit[1]))
	}

	return ret, nil
}

type Cookbook struct {
	Recipies Recipies
}

func NewCookbook() Cookbook {
	return Cookbook{}
}

func (c *Cookbook) GetRecipeByName(name string) (Recipe, error) {
	for _, r := range c.Recipies {
		if r.Name == name {
			return r, nil
		}
	}
	return Recipe{}, fmt.Errorf("recipe does not exist")
}

func (c *Cookbook) GetFilteredRecipies(name string) Cookbook {
	cb := Cookbook{
		Recipies: make([]Recipe, 0),
	}
	for _, r := range c.Recipies {
		if strings.Contains(r.Name, name) {
			cb.Recipies = append(cb.Recipies, r)
		}
	}
	return cb
}
