package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/seandisero/necronomnomicon/internal/auth"
	"github.com/seandisero/necronomnomicon/internal/cookbook"
	"github.com/seandisero/necronomnomicon/internal/tmpl"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func getDB() (*sql.DB, error) {
	db_url := os.Getenv("DB_URL")
	db_token := os.Getenv("DB_TOKEN")
	db, err := sql.Open("libsql", db_url+"?authToken="+db_token)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func OptionalJWTMiddeware() echo.MiddlewareFunc {
	config := echojwt.Config{
		Skipper: func(c echo.Context) bool {
			cookie, err := c.Cookie("necro-auth")
			if err != nil {
				return true
			}
			if cookie.Value == "" {
				return true
			}
			return false
		},
		SigningKey:  []byte(auth.TokenSecretString),
		TokenLookup: "cookie:necro-auth",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
	}
	return echojwt.WithConfig(config)
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("could not get port")
	}

	e := echo.New()

	loggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"status":${status},` +
			`"method":"${method}","uri":"${uri}",` +
			`"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}
	e.Use(middleware.LoggerWithConfig(loggerConfig))

	e.Static("/css", "css")
	e.Renderer = tmpl.NewTemplate()

	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("could not enable foreign keys")
	}
	// cb := cookbook.MakeMockCookbook()
	// cb.DB = database.New(db)
	cb := cookbook.NewCookbook(db)

	e.Use(OptionalJWTMiddeware())

	r := e.Group("")
	jwtConfig := echojwt.Config{
		SigningKey:  []byte(auth.TokenSecretString),
		TokenLookup: "cookie:necro-auth",
		ErrorHandler: func(c echo.Context, err error) error {
			return nil
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
	}
	r.Use(echojwt.WithConfig(jwtConfig))

	e.GET("/", cb.HandlerGetHome)
	e.GET("/main", cb.HandlerGetMainPage)

	e.GET("/login", cb.HandlerGetLoginPage)
	e.POST("/login", cb.HandlerLogin)
	e.POST("/logout", cb.HandlerLogout)

	e.GET("/signup", cb.HandlerGetSignUpPage)
	e.POST("/signup", cb.HandlerCreateUser)

	e.GET("/recipe", cb.HandlerGetRecipe)
	e.POST("/recipe", cb.HandlerPostRecipe)

	e.GET("/recipe/load/:index", cb.HendlerLoadMoreRecipes)
	e.GET("/recipe/grid", cb.HandlerGetRecipeGrid)
	e.GET("/recipe-form", cb.HandlerGetRecipeForm)
	e.POST("/recipe-search", cb.HandlerSearchRecipes)

	e.Logger.Fatal(e.Start(":" + port))
}
