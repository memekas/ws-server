package auth

import (
	"github.com/dgrijalva/jwt-go"
)

// Token JWT
type Token struct {
	jwt.StandardClaims
	UserID uint
}
