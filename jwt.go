package env

import (
	"os"

	"github.com/golang-jwt/jwt"
)

// ================================================================
// Env: Jwt
// ================================================================
type Jwt struct {
	SigningKey    []byte
	SigningMethod *jwt.SigningMethodHMAC
}

func NewJwt() (*Jwt, error) {
	var jwtSigningMethod *jwt.SigningMethodHMAC
	switch os.Getenv("JWT_SIGNING_METHOD") {
	case "HS256":
		jwtSigningMethod = jwt.SigningMethodHS256
	case "HS384":
		jwtSigningMethod = jwt.SigningMethodHS384
	default: // case "HS512":
		jwtSigningMethod = jwt.SigningMethodHS512
	}

	return &Jwt{
		SigningKey:    []byte(os.Getenv("JWT_SIGNING_KEY")),
		SigningMethod: jwtSigningMethod,
	}, nil
}

func (e Jwt) GenToken(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(e.SigningMethod, claims).SignedString(e.SigningKey)
}

func (e Jwt) ParseWithClaims(tokenStr string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return e.SigningKey, nil
	})
}

func (e Jwt) Parse(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return e.SigningKey, nil
	})
}
