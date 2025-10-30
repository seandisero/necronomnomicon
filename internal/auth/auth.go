package auth

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type argonParams struct {
	threads   uint8
	saltLen   uint8
	time      uint32
	memory    uint32
	keyLength uint32
}

var authParams = argonParams{
	threads:   1,
	saltLen:   16,
	time:      2,
	memory:    16 * 1024,
	keyLength: 32,
}

func AuthValidator(key string, c echo.Context) (bool, error) {
	userID, err := ValidateJWT(key)
	if err != nil {
		return false, err
	}

	ctx := context.WithValue(c.Request().Context(), "userID", userID)
	c.Request().WithContext(ctx)
	return true, nil
}

func AuthSkipper(c echo.Context) bool {
	return true
}

func GetAuthConfig() middleware.KeyAuthConfig {
	config := middleware.KeyAuthConfig{
		Skipper:   AuthSkipper,
		KeyLookup: "cookie:necro-auth",
		Validator: AuthValidator,
	}
	return config
}
