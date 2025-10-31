package cookbook

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/seandisero/necronomnomicon/internal/database"
)

func (cb *Cookbook) HandlerGetRecipe(c echo.Context) error {
	userID, err := getIDFromContext(c)
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	dbRecipe, err := cb.DB.GetRecipeByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	recipe, err := RecipeFromDBRecipe(dbRecipe)
	if err != nil {
		return err
	}
	returnParams := struct {
		UserID int64
		Recipe
	}{
		userID,
		recipe,
	}
	return c.Render(http.StatusOK, "recipe", returnParams)
}

func (cb *Cookbook) HandlerDeleteRecipe(c echo.Context) error {
	// TODO: this still needs work
	userID, err := getIDFromContext(c)
	idString := c.Param("id")
	recipeId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		slog.Error("error deleting recipe", "error", err)
		return err
	}
	recipe, err := cb.DB.GetRecipeByID(c.Request().Context(), recipeId)
	if err != nil {
		slog.Error("error deleting recipe", "error", err)
		return err
	}

	if !recipe.CreatorID.Valid {
		slog.Error("error deleting recipe", "error", "unauthorized")
		return c.NoContent(http.StatusUnauthorized)
	}
	if userID != recipe.CreatorID.Int64 {
		slog.Error("error deleting recipe", "error", "unauthorized")
		return c.NoContent(http.StatusUnauthorized)
	}

	recipeName := recipe.Name

	err = cb.DB.DeleteRecipe(c.Request().Context(), recipeId)
	if err != nil {
		slog.Error("error deleting recipe", "error", err)
		return err
	}

	deleteConfirmation := struct {
		RecipeName string
	}{
		RecipeName: recipeName,
	}

	return c.Render(
		http.StatusAccepted,
		"recipe-deleted-confirm-page",
		deleteConfirmation,
	)
}

func (cb *Cookbook) HandlerEditRecipe(c echo.Context) error {
	formData := formDataFromContext(c)
	formData.IsEdit = true

	idString := c.Param("id")
	recipeId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return ReturnWithFormData(c, formData, fmt.Errorf("internal server error"))
	}

	userID, err := getIDFromContext(c)
	if err != nil {
		slog.Error("user trying to post recipe without credentials", "error", err)
		ReturnWithFormData(c, formData, fmt.Errorf("unauthorized"))
	}

	editedRecipe, err := recipeFromFormData(formData)
	if err != nil {
		return ReturnWithFormData(c, formData, err)
	}

	params := database.GetUsersRecipeByIDParams{
		ID:        recipeId,
		CreatorID: sql.NullInt64{Int64: userID, Valid: true},
	}
	dbRecipe, err := cb.DB.GetUsersRecipeByID(c.Request().Context(), params)
	if err != nil {
		slog.Error("error editing recipe", "error", err)
		return ReturnWithFormData(c, formData, fmt.Errorf("internal server error"))
	}

	if cb.RecipeExists(c, editedRecipe) && dbRecipe.Name != editedRecipe.Name {
		return ReturnWithFormData(c, formData, fmt.Errorf("name already exists"))
	}

	editRecipeParams := database.EditRecipeParams{
		ID:          recipeId,
		Name:        editedRecipe.Name,
		Ingredients: ComposeIngredientsForDB(editedRecipe.Ingredients),
		Steps:       strings.Join(editedRecipe.Steps, "?"),
		Notes: sql.NullString{
			Valid:  editedRecipe.Notes != "",
			String: editedRecipe.Notes,
		},
	}

	dbRecipeEdit, err := cb.DB.EditRecipe(c.Request().Context(), editRecipeParams)
	if err != nil {
		slog.Error("could not edit recipe", "error", err)
		return ReturnWithFormData(c, formData, fmt.Errorf("internal server error"))
	}

	finalRecipe, err := RecipeFromDBRecipe(dbRecipeEdit)

	returnParams := struct {
		UserID int64
		Recipe
	}{
		userID,
		finalRecipe,
	}

	return c.Render(http.StatusOK, "recipe", returnParams)
}

func (cb *Cookbook) HandlerEditRecipeForm(c echo.Context) error {
	userID, err := getIDFromContext(c)
	idString := c.Param("id")
	recipeId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		slog.Error("error deleting recipe", "error", err)
		return err
	}

	params := database.GetUsersRecipeByIDParams{
		ID:        recipeId,
		CreatorID: sql.NullInt64{Int64: userID, Valid: true},
	}

	recipe, err := cb.DB.GetUsersRecipeByID(c.Request().Context(), params)
	if err != nil {
		slog.Error("error deleting recipe", "error", err)
		return err
	}

	formData, err := formDataFromDBRecipe(recipe)
	formData.IsEdit = true
	if err != nil {
		return ReturnWithFormData(c, formData, err)
	}

	return ReturnWithFormData(c, formData, err)
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
	recipes := cb.GetFilteredRecipesByName(c, name)
	if len(recipes) > 21 {
		data, err := makePageData(c, cb, userID)
		if err != nil {
			log.Fatal("something went wrong making page data during serach")
		}
		return c.Render(http.StatusOK, "recipe-grid", data)
	}
	return c.Render(http.StatusOK, "recipe-search-grid",
		struct{ Recipes Recipes }{Recipes: recipes})
}
