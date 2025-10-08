package main

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/seandisero/necronomnomicon/internal/cookbook"
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
	cb := cookbook.MakeMockCookbook()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", cb)
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

		if name == "" {
			c.Error(fmt.Errorf("no name provided"))
			return c.Render(http.StatusBadRequest, "new-recipe-form", newRecipeFormData())
		}

		parsedIngredients, err := cookbook.ParseIngredients(ingredients)
		if err != nil {
			c.Error(err)
		}

		recipe := cookbook.NewRecipe(name, parsedIngredients, strings.Split(steps, "\n"), notes)

		cb.Recipies = append(cb.Recipies, recipe)

		var i any
		return c.Render(http.StatusOK, "getnewrecipeform", i)
	})

	e.GET("/recipe", func(c echo.Context) error {
		recipeName := c.FormValue("name")
		slog.Info("requested recipe", "name", recipeName)
		recipe, err := cb.GetRecipeByName(recipeName)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "recipe", recipe)
	})

	e.GET("/recipe/grid", func(c echo.Context) error {
		return c.Render(http.StatusOK, "recipe-grid", cb)
	})

	e.GET("/search-bar", func(c echo.Context) error {
		var i any
		return c.Render(http.StatusOK, "get-new-recipe-form", i)
	})

	e.POST("/recipe-search", func(c echo.Context) error {
		name := c.FormValue("search")
		recipies := cb.GetFilteredRecipies(name)
		for _, r := range recipies.Recipies {
			fmt.Println(r.Name)
		}
		return c.Render(http.StatusOK, "recipe-grid", recipies)
	})

	e.Logger.Fatal(e.Start(":3030"))
}
