package security

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SecretKeySigned = []byte("94647177_Mc")

func GenerateJsonWebToken(userId int) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userId"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString(SecretKeySigned)
}

func ValidToken(r *http.Request) error {
	tokenString := extractToken(r)
	_, err := jwt.Parse(tokenString, returnVerificationKey)
	if err != nil {
		return err
	}
	return nil
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func returnVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Método de assinatura inesperado: %v", token.Header["alg"])
	}

	return SecretKeySigned, nil
}
