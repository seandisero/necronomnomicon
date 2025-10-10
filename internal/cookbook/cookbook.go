package cookbook

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Cookbook struct {
	Recipes Recipes
}

func NewCookbook() Cookbook {
	return Cookbook{
		Recipes: make([]Recipe, 0),
	}
}

func (c *Cookbook) GetRecipeByName(name string) (Recipe, error) {
	for _, r := range c.Recipes {
		if r.Name == name {
			return r, nil
		}
	}
	return Recipe{}, fmt.Errorf("recipe does not exist")
}

func (c *Cookbook) GetFilteredRecipes(name string) Cookbook {
	cb := Cookbook{
		Recipes: make([]Recipe, 0),
	}
	for _, r := range c.Recipes {
		if strings.Contains(r.Name, name) {
			cb.Recipes = append(cb.Recipes, r)
		}
	}
	return cb
}

func (cb *Cookbook) HandlerGetHome(c echo.Context) error {
	first := make([]Recipe, 20)
	end := min(len(cb.Recipes), 20)
	first = cb.Recipes[:end]
	return c.Render(http.StatusOK, "index", struct {
		Recipes Recipes
		Last    struct {
			Recipe Recipe
			Index  int
		}
	}{
		Recipes: first,
		Last: struct {
			Recipe Recipe
			Index  int
		}{
			Recipe: cb.Recipes[end],
			Index:  21,
		},
	})
}

func (cb *Cookbook) HendlerLoadMoreRecipes(c echo.Context) error {
	index := c.Param("index")
	idx, err := strconv.Atoi(index)
	if err != nil {
		return err
	}
	slog.Info("looking for index", "index found", index)
	first := make([]Recipe, 20)
	end := min(len(cb.Recipes), 21+idx)
	first = cb.Recipes[idx : end-1]

	return c.Render(http.StatusOK, "more-cards", struct {
		Recipes Recipes
		Last    struct {
			Recipe Recipe
			Index  int
		}
	}{
		Recipes: first,
		Last: struct {
			Recipe Recipe
			Index  int
		}{
			Recipe: cb.Recipes[end],
			Index:  21 + idx,
		},
	})
}

func (cb *Cookbook) HandlerGetRecipeForm(c echo.Context) error {
	return c.Render(http.StatusOK, "recipe-form", newRecipeFormData())
}

func (cb *Cookbook) HandlerPostRecipe(c echo.Context) error {
	name := c.FormValue("name")
	ingredients := c.FormValue("ingredients")
	steps := c.FormValue("steps")
	notes := c.FormValue("notes")

	if name == "" {
		c.Error(fmt.Errorf("no name provided"))
		return c.Render(http.StatusBadRequest, "new-recipe-form", newRecipeFormData())
	}

	parsedIngredients, err := ParseIngredients(ingredients)
	if err != nil {
		c.Error(err)
	}

	recipe := NewRecipe(name, parsedIngredients, strings.Split(steps, "\n"), notes)

	cb.Recipes = append(cb.Recipes, recipe)

	var i any
	return c.Render(http.StatusOK, "new-recipe-and-search-bar", i)

}

func (cb *Cookbook) HandlerGetRecipe(c echo.Context) error {
	recipeName := c.FormValue("name")
	slog.Info("requested recipe", "name", recipeName)
	recipe, err := cb.GetRecipeByName(recipeName)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "recipe", recipe)
}

func (cb *Cookbook) HandlerGerRecipeGrid(c echo.Context) error {
	first := make([]Recipe, 20)
	end := min(len(cb.Recipes), 20)
	first = cb.Recipes[:end]
	return c.Render(http.StatusOK, "recipe-grid", struct {
		Recipes Recipes
		Last    struct {
			Recipe Recipe
			Index  int
		}
	}{
		Recipes: first,
		Last: struct {
			Recipe Recipe
			Index  int
		}{
			Recipe: cb.Recipes[end],
			Index:  21,
		},
	})
}

func (cb *Cookbook) HandlerGetSearchBar(c echo.Context) error {
	return c.Render(http.StatusOK, "recipe-search", struct{}{})
}

func (cb *Cookbook) HandlerSearchRecipes(c echo.Context) error {
	name := c.FormValue("search")
	if name == "" {
		first := make([]Recipe, 20)
		end := min(len(cb.Recipes), 20)
		first = cb.Recipes[:end]
		return c.Render(http.StatusOK, "recipe-grid", struct {
			Recipes Recipes
			Last    struct {
				Recipe Recipe
				Index  int
			}
		}{
			Recipes: first,
			Last: struct {
				Recipe Recipe
				Index  int
			}{
				Recipe: cb.Recipes[end],
				Index:  21,
			},
		})
	}
	recipes := cb.GetFilteredRecipes(name)
	return c.Render(http.StatusOK, "recipe-search-grid", recipes)
}
