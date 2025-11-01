package cookbook

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/seandisero/necronomnomicon/internal/database"
)

type Cookbook struct {
	Recipes Recipes
	DB      *database.Queries
}

func NewCookbook(db *sql.DB) Cookbook {
	return Cookbook{
		Recipes: make([]Recipe, 0),
		DB:      database.New(db),
	}
}

func (cb *Cookbook) GetAllRecipesFromDB(c echo.Context) error {
	dbRecipes, err := cb.DB.GetAllRecipes(c.Request().Context())
	if err != nil {
		return err
	}
	recipes := make([]Recipe, len(dbRecipes))
	for i, dbRecipe := range dbRecipes {
		recipe, err := RecipeFromDBRecipe(dbRecipe)
		if err != nil {
			return err
		}
		recipes[i] = recipe
	}
	cb.Recipes = recipes
	return nil
}

func (cb *Cookbook) GetFilteredRecipesByName(c echo.Context, name string) Recipes {
	recipes := make([]Recipe, 0)
	err := cb.GetAllRecipesFromDB(c)
	if err != nil {
		return nil
	}
	for _, r := range cb.Recipes {
		if strings.Contains(r.Name, name) {
			recipes = append(recipes, r)
		}
	}
	return recipes
}

func (cb *Cookbook) GetFilteredRecipesByUserAndName(c echo.Context, name string, userID int64) Recipes {
	recipes := make([]Recipe, 0)
	err := cb.GetAllRecipesFromDB(c)
	if err != nil {
		return nil
	}
	for _, r := range cb.Recipes {
		if strings.Contains(r.Name, name) && r.CreatorID == userID {
			recipes = append(recipes, r)
		}
	}
	return recipes
}

func getIDFromContext(c echo.Context) (int64, error) {
	if c.Get("user") == nil {
		slog.Warn("c.Get('user') is nil, user is not logged in")
		return -1, nil
	}
	user := c.Get("user").(*jwt.Token)
	userIDString, err := user.Claims.GetSubject()
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(userIDString, 10, 64)
}

func (cb *Cookbook) HandlerGetHome(c echo.Context) error {
	userID, err := getIDFromContext(c)
	if err != nil {
		slog.Error("error getting id from context", "error", err)
	}

	data, err := makePageData(c, cb, userID)
	if err != nil {
		slog.Error("logging error in Handler get home", "error", err)
		return c.Render(http.StatusOK, "index", data)
	}
	if data == nil {
		slog.Error("page data is nil")
	}
	return c.Render(http.StatusOK, "index", data)
}

func (cb *Cookbook) HandlerGetMainPage(c echo.Context) error {
	userID, err := getIDFromContext(c)
	if err != nil {
		slog.Error("error getting id from context", "error", err)
	}
	data, err := makePageData(c, cb, userID)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "main-page", data)
}
