package auth

import (
	"context"
	"errors"
	"os"
	"time"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(UserID string, rememberMe bool) (string, error) {
	// expiredToken := time.Now().Add(1 * time.Hour)

	var expiredToken time.Time
    if rememberMe {
        expiredToken = time.Now().Add(72 * time.Hour) // Token berlaku selama 3 hari jika rememberMe diatur
    } else {
        expiredToken = time.Now().Add(1 * time.Hour) // Token berlaku selama 1 jam jika rememberMe tidak diatur
    }
	
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

func TokenVerify(tokenCode string) (*Claims, error) {
	claimsToken := &Claims{}
	token, err := jwt.ParseWithClaims(tokenCode, claimsToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	} else {
		return nil, errors.New("Invalid claims")
	}
}


func  GetTokenUserLogin(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == ""{
		
		return "", http.ErrNoCookie
	}
	return userID, nil
}