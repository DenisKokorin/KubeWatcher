package jwt

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.RegisteredClaims
	UUID     string
	Name     string
	Email    string
	Position string
	Team     string
	Role     string
}
