package cookbook

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (cb *Cookbook) HendlerLoadMoreRecipes(c echo.Context) error {
	userID, err := getIDFromContext(c)
	if err != nil {
		userID = -1
	}

	index := c.Param("index")
	idx, err := strconv.Atoi(index)
	if err != nil {
		return err
	}

	data, err := cb.drawRecipeData(c, idx+1, userID)
	if err != nil {
		c.Render(http.StatusBadRequest, "more-cards", data)
		return fmt.Errorf("no more recipe data")
	}

	return c.Render(http.StatusOK, "more-cards", data)
}

func (cb *Cookbook) HandlerGetRecipeGrid(c echo.Context) error {
	userID, err := getIDFromContext(c)
	if err != nil {
		userID = -1
	}
	data, err := makePageData(c, cb, userID)
	if err != nil {
		return fmt.Errorf("could not get recipes")
	}
	return c.Render(http.StatusOK, "recipe-grid", data)
}
