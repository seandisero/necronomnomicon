package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/seandisero/necronomnomicon/internal/cookbook"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/css", "css")

	e.Renderer = newTemplate()
	cb := cookbook.MakeMockCookbook()

	e.GET("/", cb.HandlerGetHome)
	e.POST("/recipe", cb.HandlerPostRecipe)
	e.GET("/recipe", cb.HandlerGetRecipe)
	e.GET("/recipe/load/:index", cb.HendlerLoadMoreRecipes)
	e.GET("/recipe/grid", cb.HandlerGerRecipeGrid)
	e.GET("/recipe-form", cb.HandlerGetRecipeForm)
	e.GET("/search-bar", cb.HandlerGetSearchBar)
	e.POST("/recipe-search", cb.HandlerSearchRecipes)

	e.Logger.Fatal(e.Start(":3030"))
}
