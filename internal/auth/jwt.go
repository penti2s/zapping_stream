package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"time"
	"zapping_stream/internal/model"
)

type JWTClaims struct {
	UserID uint
	Email  string
	Name   string
	jwt.StandardClaims
}

func GenerateToken(user model.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET_KEY")
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	println(tokenString)

	splitToken := strings.Split(tokenString, "Bearer ")
	if len(splitToken) != 2 {
		return nil, errors.New("formato de token inválido")
	}
	jwtToken := splitToken[1]

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
