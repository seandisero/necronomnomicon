package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/seandisero/necronomnomicon/internal/auth"
	"github.com/seandisero/necronomnomicon/internal/cookbook"
	"github.com/seandisero/necronomnomicon/internal/tmpl"

	"github.com/tursodatabase/go-libsql"
)

func EchoLogger() echo.MiddlewareFunc {
	// TODO: newer config is RequestLoggerConfig{}
	loggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"status":${status},` +
			`"method":"${method}","uri":"${uri}",` +
			`"error_msg":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}
	// TODO: this should be changed to middleware.RequestLoggerWithConfig()
	return middleware.LoggerWithConfig(loggerConfig)
}

func SetupRouting(e *echo.Echo, cb *cookbook.Cookbook) {
	e.Static("/css", "css")
	e.Static("/images", "images")
	e.Renderer = tmpl.NewTemplate()

	e.Use(EchoLogger())
	e.Use(auth.OptionalJWTMiddeware())

	r := e.Group("")
	r.Use(auth.NonOptionalJWTMiddleware())

	e.GET("/", cb.HandlerGetHome)
	e.GET("/main", cb.HandlerGetMainPage)

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.File("images/necronomnom_icon.png")
	})

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

	e.GET("/recipe/load/:index", cb.HendlerLoadMoreRecipes)
	e.GET("/recipe/grid", cb.HandlerGetRecipeGrid)
	e.GET("/recipe-form", cb.HandlerGetRecipeForm)
	e.POST("/recipe-search", cb.HandlerSearchRecipes)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("no .env file found, using global env variables", "error", err)
	}
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
