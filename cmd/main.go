package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/seandisero/necronomnomicon/internal/auth"
	"github.com/seandisero/necronomnomicon/internal/cookbook"
	"github.com/seandisero/necronomnomicon/internal/tmpl"

	"github.com/tursodatabase/go-libsql"
)

// func getDB() (*sql.DB, error) {
// 	dbName := "necro.db"
// 	dir, err := os.MkdirTemp("", "libsql-*")
// 	if err != nil {
// 		slog.Error("could not make temp dir for necro db", "error", err)
// 		return nil, err
// 	}
// 	db_url := os.Getenv("DB_URL")
// 	db_token := os.Getenv("DB_TOKEN")
// 	db, err := sql.Open("libsql", db_url+"?authToken="+db_token)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

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

func NonOptionalJWTMiddleware() echo.MiddlewareFunc {
	jwtConfig := echojwt.Config{
		SigningKey:  []byte(auth.TokenSecretString),
		TokenLookup: "cookie:necro-auth",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
	}
	return echojwt.WithConfig(jwtConfig)
}

func EchoLogger() echo.MiddlewareFunc {
	loggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"status":${status},` +
			`"method":"${method}","uri":"${uri}",` +
			`"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}
	return middleware.LoggerWithConfig(loggerConfig)
}

func SetupRouting(e *echo.Echo, cb *cookbook.Cookbook) {
	e.Static("/css", "css")
	e.Renderer = tmpl.NewTemplate()

	e.Use(EchoLogger())
	e.Use(OptionalJWTMiddeware())

	r := e.Group("")
	r.Use(NonOptionalJWTMiddleware())

	recipeEdit := e.Group("/recipe/edit/:id")
	recipeEdit.Use(NonOptionalJWTMiddleware())

	e.GET("/", cb.HandlerGetHome)
	e.GET("/main", cb.HandlerGetMainPage)

	e.GET("/login", cb.HandlerGetLoginPage)
	e.POST("/login", cb.HandlerLogin)
	e.POST("/logout", cb.HandlerLogout)

	e.GET("/signup", cb.HandlerGetSignUpPage)
	e.POST("/signup", cb.HandlerCreateUser)

	r.POST("/recipe", cb.HandlerPostRecipe)
	e.GET("/recipe/:id", cb.HandlerGetRecipe)
	r.DELETE("/recipe/:id", cb.HandlerDeleteRecipe)
	r.GET("/recipe/edit/:id", cb.HandlerEditRecipeForm)
	r.PUT("/recipe/:id", cb.HandlerEditRecipe)

	// recipeEdit.GET("/name", cb.HandlerStartRecipeNameEdit)
	// recipeEdit.GET("/ingredients", cb.HandlerStartRecipeIngredientsEdit)
	// recipeEdit.GET("/steps", cb.HandlerStartRecipeStepsEdit)
	// recipeEdit.GET("/notes", cb.HandlerStartRecipeNotesEdit)

	e.GET("/recipe/load/:index", cb.HendlerLoadMoreRecipes)
	e.GET("/recipe/grid", cb.HandlerGetRecipeGrid)
	e.GET("/recipe-form", cb.HandlerGetRecipeForm)
	e.POST("/recipe-search", cb.HandlerSearchRecipes)
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("could not get port")
	}

	dbName := "necro.db"
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		slog.Error("could not make temp dir for necro db", "error", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	slog.Info("file path for db", "path", dbPath)

	db_url := os.Getenv("DB_URL")
	db_token := os.Getenv("DB_TOKEN")

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, db_url,
		libsql.WithAuthToken(db_token),
		libsql.WithSyncInterval(30*time.Minute),
		libsql.WithEncryption("ENCRYPTION_STRING"),
	)
	if err != nil {
		slog.Error("error creating connector", "error", err)
		os.Exit(1)
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("could not enable foreign keys")
	}

	cb := cookbook.NewCookbook(db)
	e := echo.New()
	SetupRouting(e, &cb)

	e.Logger.Fatal(e.Start(":" + port))
}
