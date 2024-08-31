package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func BearerHeader(headers http.Header) (string, error) {

	token := headers.Get("Authorization")

	if token == "" {
		return "", errors.New("auth header not found")
	}

	splitToken := strings.Split(token, " ")

	if len(splitToken) < 2 || splitToken[0] != "Bearer" {
		return "", errors.New("auth header not what expected")
	}

	return splitToken[1], nil
}

func ValidateToken(tokenstring, tokenSecret string) (string, error) {
	type customClaims struct {
		jwt.RegisteredClaims
	}
	token, err := jwt.ParseWithClaims(tokenstring, &customClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })

	if err != nil {
		return "", errors.New(err.Error()) //"jwt couldn't be parsed"
	}

	userId, err := token.Claims.GetSubject()

	if err != nil {
		return "", errors.New("user id couldn't be extracted")
	}

	return userId, nil
}
