package auth

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func OptionalJWTErrorHandler(c echo.Context, err error) error {
	if err.Error() == "token has invalid claims: token is expired" {
		return nil
	}
	return nil
}

func OptionalJWTParseTokenFunc(c echo.Context, authString string) (any, error) {
	jwt_secret := os.Getenv("JWT_SECRET")
	newClaimsFunc := func(c echo.Context) jwt.Claims {
		return new(jwt.RegisteredClaims)
	}
	keyFunc := func(token *jwt.Token) (any, error) {
		return []byte(jwt_secret), nil
	}
	token, err := jwt.ParseWithClaims(authString, newClaimsFunc(c), keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func OptionalJWTMiddeware() echo.MiddlewareFunc {
	jwt_secret := os.Getenv("JWT_SECRET")
	config := echojwt.Config{
		Skipper: func(c echo.Context) bool {
			// if the cookie does not exsit or has no value then we need to skip
			if cookie, err := c.Cookie("necro-auth"); err != nil {
				return true
			} else {
				if cookie.Value == "" {
					return true
				}
			}
			// if the cookie has a value then we should check it.
			return false
		},
		SigningKey:             []byte(jwt_secret),
		TokenLookup:            "cookie:necro-auth",
		ErrorHandler:           OptionalJWTErrorHandler,
		ContinueOnIgnoredError: true,
		ParseTokenFunc:         OptionalJWTParseTokenFunc,
	}
	return echojwt.WithConfig(config)
}

func NonOptionalJWTMiddleware() echo.MiddlewareFunc {
	jwt_secret := os.Getenv("JWT_SECRET")
	jwtConfig := echojwt.Config{
		SigningKey:  []byte(jwt_secret),
		TokenLookup: "cookie:necro-auth",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
	}
	return echojwt.WithConfig(jwtConfig)
}
