package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TokenType string

const TokenTypeAccess TokenType = "necro-access"

const TokenSecretString string = "SecretString"

const JWTExpireTime time.Duration = 24 * time.Hour

func MakeJWT(userID int64) (string, error) {
	return makeJWT(userID, TokenSecretString, JWTExpireTime)
}

func makeJWT(userID int64, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	userIDString := strconv.FormatInt(userID, 10)
	webToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userIDString,
	})
	return webToken.SignedString(signingKey)
}

func ParseToken(c echo.Context, tokenString string) (interface{}, error) {
	claims := jwt.RegisteredClaims{}
	keyFunk := func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSecretString), nil
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		keyFunk,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ValidateJWT(tokenString string) (int64, error) {
	return validateJWT(tokenString, TokenSecretString)
}

func validateJWT(tokenString, tokenSecret string) (int64, error) {
	claims := jwt.RegisteredClaims{}
	keyFunk := func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		keyFunk,
	)
	if err != nil {
		return -1, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return -1, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return -1, err
	}
	if issuer != string(TokenTypeAccess) {
		return -1, fmt.Errorf("invalid issuer")
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return -1, err
	}

	return userID, nil
}
