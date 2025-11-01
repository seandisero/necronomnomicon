package cookbook

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/seandisero/necronomnomicon/internal/auth"
)

func (cb *Cookbook) HandlerGetLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login-page", nil)
}

func (cb *Cookbook) HandlerLogin(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")
	user, err := cb.DB.GetUserByName(c.Request().Context(), name)
	if err != nil {
		slog.Error("could not get user by name", "name", name, "error", err)
		return err
	}

	valid, err := auth.ValidatePassword(password, user.HashedPassword)
	if !valid {
		if err != nil {
			slog.Error("error validating password", "error", err)
		}
		// TODO: this section should keep the user on the loging page with a
		// error that the password was wrong.
		return c.Render(http.StatusBadRequest, "index", nil)
	}

	jwt, err := auth.MakeJWT(user.ID)
	if err != nil {
		slog.Error("error making jwt", "error", err)
		return c.Render(http.StatusInternalServerError, "index", nil)
	}
	cookie := new(http.Cookie)
	cookie.Name = "necro-auth"
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = 1
	cookie.Value = jwt
	c.SetCookie(cookie)

	data, err := makePageData(c, cb, user.ID)
	if err != nil {
		return c.Render(http.StatusOK, "index", nil)
	}
	return c.Render(http.StatusOK, "index", data)
}

func (cb *Cookbook) HandlerLogout(c echo.Context) error {
	emptyCookie := new(http.Cookie)
	emptyCookie.Name = "necro-auth"
	c.SetCookie(emptyCookie)
	data, err := makePageData(c, cb, -1)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "index", data)
}
