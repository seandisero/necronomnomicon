package cookbook

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cb *Cookbook) HandlerGetRecipeForm(c echo.Context) error {
	return ReturnWithFormData(c, EmptyRecipeFormData(), nil)
}

func ReturnWithFormData(c echo.Context, data RecipeFormData, err error) error {
	if err != nil {
		switch err.Error() {
		case "must provide a name":
			data.ErrorName = err.Error()
			return c.Render(http.StatusOK, "recipe-form", data)
		case "name already exists":
			data.ErrorName = err.Error()
			return c.Render(http.StatusOK, "recipe-form", data)
		case "malformed ingredient: must be number measure name":
			data.ErrorIngredients = err.Error()
			return c.Render(http.StatusOK, "recipe-form", data)
		case "start or line must be a number":
			data.ErrorIngredients = err.Error()
			return c.Render(http.StatusOK, "recipe-form", data)
		case "must provide ingredients list":
			data.ErrorIngredients = err.Error()
			return c.Render(http.StatusOK, "recipe-form", data)
		}
	}
	return c.Render(http.StatusOK, "recipe-form", data)
}

func (cb *Cookbook) HandlerPostRecipe(c echo.Context) error {
	userID, err := getIDFromContext(c)
	if err != nil {
		slog.Error("could nto get user id from context", "error", err)
		userID = -1
	}
	name := c.FormValue("name")
	ingredients := c.FormValue("ingredients")
	steps_string := c.FormValue("steps")
	notes := c.FormValue("notes")

	newFormData := NewRecipeFormData(name, ingredients, steps_string, notes)

	if name == "" {
		return ReturnWithFormData(c, newFormData, fmt.Errorf("must have a name"))
	}

	recipe, err := CreateRecipe(name, ingredients, steps_string, notes)
	if err != nil {
		slog.Error("problem creating recipe", "error", err)
		return ReturnWithFormData(c, newFormData, err)
	}

	_, err = cb.DB.GetRecipeByName(c.Request().Context(), recipe.Name)
	if err == nil {
		return ReturnWithFormData(c, newFormData, fmt.Errorf("name already exists"))
	}

	cb.AddRecipeToDB(c, recipe, userID)

	data, err := makePageData(c, cb, userID)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "recipe-grid", data)
}
