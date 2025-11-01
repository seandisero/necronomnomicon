package cookbook

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/seandisero/necronomnomicon/internal/auth"
	"github.com/seandisero/necronomnomicon/internal/database"
)

func (cb *Cookbook) HandlerGetSignUpPage(c echo.Context) error {
	return c.Render(http.StatusOK, "signup-form", nil)
}

func (cb *Cookbook) HandlerCreateUser(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")
	if name == "" {
		return fmt.Errorf("no name parameter")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	createUserParams := database.CreateUserParams{
		Username:       name,
		HashedPassword: hashedPassword,
	}

	user, err := cb.DB.CreateUser(c.Request().Context(), createUserParams)
	if err != nil {
		return err
	}
	slog.Info("created user", "username", user.Username, "id", user.ID)

	jwt, err := auth.MakeJWT(user.ID)
	if err != nil {
		slog.Error("error making jwt during signup", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	cookie := new(http.Cookie)
	cookie.Name = "necro-auth"
	cookie.Value = jwt
	c.SetCookie(cookie)

	return c.Render(http.StatusOK, "welcome-page", user)
}

func (cb *Cookbook) HandlerWelcomePage(c echo.Context) error {
	data, err := makePageData(c, cb, -1)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "main-page", data)
}
