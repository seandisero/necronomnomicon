package cookbook

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cb *Cookbook) HandlerGetRecipe(c echo.Context) error {
	recipeName := c.FormValue("name")
	dbRecipe, err := cb.DB.GetRecipeByName(c.Request().Context(), recipeName)
	if err != nil {
		return err
	}
	recipe, err := RecipeFromDBRecipe(dbRecipe)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "recipe", recipe)
}

func (cb *Cookbook) HandlerSearchRecipes(c echo.Context) error {
	userID, err := getIDFromContext(c)
	if err != nil {
		slog.Error("error getting id from context", "error", err)
	}
	name := c.FormValue("search")
	if name == "" {
		data, err := makePageData(c, cb, userID)
		if err != nil {
			log.Fatal("something went wrong making page data during serach")
		}
		return c.Render(http.StatusOK, "recipe-grid", data)
	}
	recipes := cb.GetFilteredRecipes(name)
	return c.Render(http.StatusOK, "recipe-search-grid", recipes)
}

