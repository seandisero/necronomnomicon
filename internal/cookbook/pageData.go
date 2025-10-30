package cookbook

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type PageData struct {
	IsAuthenticated bool
	GridData        RecipeGridData
}

type RecipeGridData struct {
	DrawLastCard bool
	Current      int
	LastCard     Recipe
	Cards        Recipes
}

func NewPageData(authenticated bool, recipeData RecipeGridData) *PageData {
	return &PageData{
		IsAuthenticated: authenticated,
		GridData:        recipeData,
	}
}

func (cb *Cookbook) drawRecipeData(c echo.Context, startIndex int, id int64) (RecipeGridData, error) {
	// TODO: filter out recipes for user id
	userRecipes, err := cb.GetAllRecipes(c)
	if err != nil {
		return RecipeGridData{}, err
	}
	if len(userRecipes) == 0 {
		slog.Error("no recipes found for user", "user id", id)
		return RecipeGridData{}, nil
	}

	if startIndex > len(userRecipes) {
		slog.Error("reached the end of the list")
		return RecipeGridData{}, nil
	}

	endIndex := min(len(userRecipes), startIndex+21)
	recipes := userRecipes[startIndex:endIndex]

	if endIndex >= len(userRecipes)-1 {
		return RecipeGridData{
			Cards:        recipes,
			Current:      endIndex,
			DrawLastCard: false,
		}, nil
	}

	return RecipeGridData{
		Cards:        recipes,
		DrawLastCard: true,
		Current:      endIndex,
		LastCard:     userRecipes[endIndex],
	}, nil
}

func makePageData(c echo.Context, cb *Cookbook, id int64) (*PageData, error) {
	data, err := makePageDataFromIndex(c, cb, 0, id)
	if data == nil {
		slog.Error("error getting page data", "error", err)
		return nil, fmt.Errorf("data is nil in makePageData %v", err)
	}
	return data, err
}

func makePageDataFromIndex(c echo.Context, cb *Cookbook, index int, id int64) (*PageData, error) {
	authenticated := false
	if id != -1 {
		authenticated = true
	}

	recipeData, err := cb.drawRecipeData(c, index, id)
	if err != nil {
		return nil, err
	}

	pageData := NewPageData(authenticated, recipeData)
	return pageData, nil
}

func getUserRecipes(userID int) Recipes {
	// TODO: get users recipes from db
	// TODO: have last viewed first in the list
	// TODO: cache users recipes for quick access.
	return Recipes{}
}
