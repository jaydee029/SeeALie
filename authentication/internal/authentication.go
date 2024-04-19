package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Revoke struct {
	Token      string    `json:"token"`
	Revoked_at time.Time `json:"revoked_at"`
}

func Tokenize(id uuid.UUID, secret_key string) (string, error) {
	secret_key_byte := []byte(secret_key)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chat-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(60*60) * time.Second)), // 1 hour
		Subject:   id.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret_key_byte)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func RefreshToken(id uuid.UUID, secret_key string) (string, error) {
	secret_key_byte := []byte(secret_key)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chat-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().AddDate(0, 3, 0)), // 90 days
		Subject:   id.String(),
	}

	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rt, err := refresh_token.SignedString(secret_key_byte)
	if err != nil {
		return "", err
	}
	return rt, nil

}

func Hashpassword(passwd string) ([]byte, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, errors.New("couldn't Hash the password")
	}

	return encrypted, nil
}

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

func VerifyRefresh(tokenstring, tokenSecret string) (bool, error) {
	type customClaims struct {
		jwt.RegisteredClaims
	}
	token, err := jwt.ParseWithClaims(tokenstring, &customClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })

	if err != nil {
		return false, errors.New(err.Error()) //"jwt couldn't be parsed"
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		return false, errors.New("issuer couldn't be extracted")
	}

	if issuer == "chat-refresh" {
		return true, nil
	}
	return false, nil
}
