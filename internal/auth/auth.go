package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// User represents a morpheus client information.
type User struct {
	Username string `json:"username"`
}

func (c User) String() string {
	return c.Username
}

// Claims of jwt token.
type Claims struct {
	Username string
	Iat      int64
	Exp      int64
}

// Valid checks claims issuer.
func (c Claims) Valid() error {
	if c.Exp != 0 && c.Exp < time.Now().Unix() {
		return errors.New("token has been expired")
	}

	return nil
}

// Validate given token with given secret and if it is valid it returns client information from its claims.
func Validate(tkn string, secret string) (User, error) {
	// Validating and parsing the tokenString
	token, err := jwt.ParseWithClaims(tkn, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validating if algorithm used for signing is same as the algorithm in token
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return User{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return User{}, jwt.ErrInvalidKey
	}

	return User{
		Username: claims.Username,
	}, nil
}
