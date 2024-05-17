package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(UserID string) (string, error) {
	expiredToken := time.Now().Add(1 * time.Hour)

	claimsToken := &Claims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredToken),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsToken)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func TokenVerify(tokenCode string) (string, error) {
	claimsToken := &Claims{}
	token, err := jwt.ParseWithClaims(tokenCode, claimsToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claimsToken, nil
}
