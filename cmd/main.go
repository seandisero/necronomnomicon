package main

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Ingredient struct {
	Amount  int
	Measure string
	Name    string
}

func newIngredient(name string, amount int, measure string) Ingredient {
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

func newRecipe(name string, ingredients []Ingredient, steps []string, notes string) Recipe {
	return Recipe{
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
		Notes:       notes,
	}
}

type Recipies = []Recipe

type Cookbook struct {
	Recipies Recipies
}

func (c *Cookbook) getRecipeByName(name string) (Recipe, error) {
	for _, r := range c.Recipies {
		if r.Name == name {
			return r, nil
		}
	}
	return Recipe{}, fmt.Errorf("recipe does not exist")
}

func newCookbook() Cookbook {
	return Cookbook{
		Recipies: []Recipe{
			newRecipe(
				"cookies",
				[]Ingredient{
					newIngredient("sugar", 10, "cups"),
					newIngredient("flour", 1, "cups"),
					newIngredient("milk", 100, "ml"),
				},
				[]string{
					"mix the shit",
					"bake it",
				},
				"they tast good",
			),
			newRecipe(
				"apple pie",
				[]Ingredient{
					newIngredient("egg", 1, "whole"),
					newIngredient("flour", 1, "cups"),
					newIngredient("lard", 1, "lb"),
				},
				[]string{
					"cut in the lard",
					"mix in the egg, vinigar and water",
					"bake it",
				},
				"dont over do the water or mixing",
			),
			newRecipe(
				"pumpkin pie",
				[]Ingredient{
					newIngredient("egg", 1, "whole"),
					newIngredient("flour", 1, "cups"),
					newIngredient("lard", 1, "lb"),
				},
				[]string{
					"cut in the lard",
					"mix in the egg, vinigar and water",
					"bake it",
				},
				"dont over do the water or mixing",
			),
			newRecipe(
				"meetballs",
				[]Ingredient{
					newIngredient("egg", 1, "whole"),
					newIngredient("flour", 1, "cups"),
					newIngredient("lard", 1, "lb"),
				},
				[]string{
					"cut in the lard",
					"mix in the egg, vinigar and water",
					"bake it",
				},
				"dont over do the water or mixing",
			),
		},
	}
}

func parseIngredients(ingredients string) ([]Ingredient, error) {
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
		ret = append(ret, newIngredient(ingredientSplit[2], intAmount, ingredientSplit[1]))
	}

	return ret, nil
}

type RecipeFormData struct {
	Name        string
	Ingredients string
	Steps       string
	Notes       string
}

func newRecipeFormData() RecipeFormData {
	return RecipeFormData{}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/css", "css")

	e.Renderer = newTemplate()
	cookbook := newCookbook()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", cookbook)
	})

	e.GET("/recipe-form", func(c echo.Context) error {
		return c.Render(http.StatusOK, "new-recipe-form", newRecipeFormData())
	})

	e.POST("/recipe", func(c echo.Context) error {
		fmt.Println("## getting data from post ##")
		name := c.FormValue("name")
		ingredients := c.FormValue("ingredients")
		steps := c.FormValue("steps")
		notes := c.FormValue("notes")

		fmt.Printf("## name ->\t\t\t%s\n", name)
		fmt.Printf("## ingredients ->\t\t\t%s\n", ingredients)
		fmt.Printf("## steps ->\t\t\t%s\n", steps)
		fmt.Printf("## notes ->\t\t\t%s\n", notes)

		parsedIngredients, err := parseIngredients(ingredients)
		if err != nil {
			c.Error(err)
		}

		recipe := newRecipe(name, parsedIngredients, strings.Split(steps, "\n"), notes)

		cookbook.Recipies = append(cookbook.Recipies, recipe)

		var i any
		return c.Render(http.StatusOK, "getnewrecipeform", i)
	})

	e.GET("/recipe", func(c echo.Context) error {
		recipeName := c.FormValue("name")
		slog.Info("requested recipe", "name", recipeName)
		recipe, err := cookbook.getRecipeByName(recipeName)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "recipe", recipe)
	})

	e.GET("/recipe/grid", func(c echo.Context) error {
		return c.Render(http.StatusOK, "recipe-grid", cookbook)
	})

	e.Logger.Fatal(e.Start(":3030"))
}
