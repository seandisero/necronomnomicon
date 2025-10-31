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
		case "unauthorized":
			data.ErrorIngredients = err.Error()
			return c.Render(http.StatusUnauthorized, "recipe-form", data)
		case "internal server error":
			data.ErrorIngredients = err.Error()
			return c.Render(http.StatusInternalServerError, "recipe-form", data)
		default:
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.Render(http.StatusOK, "recipe-form", data)
}

func (cb *Cookbook) HandlerPostRecipe(c echo.Context) error {
	formData := formDataFromContext(c)
	formData.IsNew = true
	userId, err := getIDFromContext(c)
	if err != nil {
		slog.Error("user trying to post recipe without credentials", "error", err)
		ReturnWithFormData(c, formData, fmt.Errorf("unauthorized"))
	}

	recipe, err := recipeFromFormData(formData)
	if err != nil {
		return ReturnWithFormData(c, formData, err)
	}

	if recipeExists := cb.RecipeExists(c, recipe); recipeExists {
		return ReturnWithFormData(c, formData, fmt.Errorf("name already in use"))
	}

	newRecipe, err := cb.AddRecipeToDB(c, recipe, userId)
	if err != nil {
		slog.Error("failed to add recipe to DB after validation", "error", err)
		return ReturnWithFormData(c, formData, err)
	}

	returnParams := struct {
		UserID int64
		Recipe
	}{
		userId,
		newRecipe,
	}
	return c.Render(http.StatusOK, "recipe", returnParams)
}
