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
	userId, err := getIDFromContext(c)
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
		userId,
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
	userID, err := getIDFromContext(c)
	idString := c.Param("id")
	recipeId, err := strconv.ParseInt(idString, 10, 64)

	name := c.FormValue("name")
	ingredients := c.FormValue("ingredients")
	steps_string := c.FormValue("steps")
	notes := c.FormValue("notes")

	if err != nil {
		slog.Error("error editing recipe", "error", err)
		return err
	}
	recipe, err := cb.DB.GetRecipeByID(c.Request().Context(), recipeId)
	if err != nil {
		slog.Error("error editing recipe", "error", err)
		return err
	}

	if !recipe.CreatorID.Valid {
		slog.Error("error editing recipe", "error", "unauthorized")
		return c.NoContent(http.StatusUnauthorized)
	}
	if userID != recipe.CreatorID.Int64 {
		slog.Error("error editing recipe", "error", "unauthorized")
		return c.NoContent(http.StatusUnauthorized)
	}

	newFormData := EditRecipeFormData(recipeId, name, ingredients, steps_string, notes)

	if name == "" {
		return ReturnWithFormData(c, newFormData, fmt.Errorf("must have a name"))
	}

	edited_recipe, err := CreateRecipe(name, ingredients, steps_string, notes)
	if err != nil {
		slog.Error("problem creating recipe", "error", err)
		return ReturnWithFormData(c, newFormData, err)
	}

	_, err = cb.DB.GetRecipeByName(c.Request().Context(), edited_recipe.Name)
	if err == nil && edited_recipe.Name != recipe.Name {
		return ReturnWithFormData(c, newFormData, fmt.Errorf("name already exists"))
	}

	ingredientsStr := ingredientsToText(edited_recipe.Ingredients)
	edited_notes := sql.NullString{}
	slog.Info("edted note", "value", notes, "created", edited_recipe.Notes)
	if edited_recipe.Notes != "" {
		slog.Info("notes is not null")
		edited_notes.Valid = true
		edited_notes.String = edited_recipe.Notes
	}

	params := database.EditRecipeParams{
		ID:          recipeId,
		Name:        edited_recipe.Name,
		Ingredients: ingredientsStr,
		Steps:       strings.Join(edited_recipe.Steps, "?"),
		Notes:       edited_notes,
	}
	cb.DB.EditRecipe(c.Request().Context(), params)

	data, err := makePageData(c, cb, userID)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "recipe-grid", data)

}

func (cb *Cookbook) HandlerEditRecipeForm(c echo.Context) error {
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

	notes := ""
	if recipe.Notes.Valid {
		notes = recipe.Notes.String
	}
	formIngredients := strings.ReplaceAll(recipe.Ingredients, "?", "\n")
	formIngredients = strings.ReplaceAll(formIngredients, ":", " ")

	data := EditRecipeFormData(
		recipe.ID,
		recipe.Name,
		formIngredients,
		strings.ReplaceAll(recipe.Steps, "?", "\n"),
		notes,
	)

	return ReturnWithFormData(c, data, nil)
}

func (cb *Cookbook) HandlerStartRecipeNameEdit(c echo.Context) error {
	userID, err := getIDFromContext(c)
	if err != nil {
		return err
	}
	recipeIDString := c.Param("id")
	recipeID, err := strconv.ParseInt(recipeIDString, 10, 64)
	recipe, err := cb.DB.GetRecipeByID(c.Request().Context(), recipeID)

	if !recipe.CreatorID.Valid || userID != recipe.CreatorID.Int64 {
		return c.NoContent(http.StatusUnauthorized)
	}

	data := struct {
		ID      int64
		OldName string
	}{
		ID:      recipe.ID,
		OldName: recipe.Name,
	}
	return c.Render(http.StatusOK, "recipe-name-box", data)
}

func (cb *Cookbook) HandlerStartRecipeIngredientsEdit(c echo.Context) error {
	return c.NoContent(500)
}

func (cb *Cookbook) HandlerStartRecipeStepsEdit(c echo.Context) error {
	return c.NoContent(500)
}

func (cb *Cookbook) HandlerStartRecipeNotesEdit(c echo.Context) error {
	return c.NoContent(500)
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
