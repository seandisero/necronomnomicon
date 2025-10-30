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

type LastItem struct {
	CanLoadMore bool
	Recipe      Recipe
	Index       int
}

func MakeLastItem(canLoadMore bool, recipe Recipe, idx int) LastItem {
	return LastItem{
		CanLoadMore: canLoadMore,
		Recipe:      recipe,
		Index:       idx,
	}
}

func NewCookbook(db *sql.DB) Cookbook {
	return Cookbook{
		Recipes: make([]Recipe, 0),
		DB:      database.New(db),
	}
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

func getIDFromContext(c echo.Context) (int64, error) {
	if c.Get("user") == nil {
		slog.Info("c.Get('user') is nil, user is not logged in")
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
	if userID != -1 {
		slog.Info("got user ID", "id", userID)
	} else {
		slog.Info("did not get user id", "id", userID)
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
