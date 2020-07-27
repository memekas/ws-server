package auth

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

func (tk *Token) Create(UserID uint) (string, error) {
	tk.UserID = UserID
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (tk *Token) Decrypt(tkString string) error {
	_, err := jwt.ParseWithClaims(tkString, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return err
}
