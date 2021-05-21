package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(username string, email string, isAdmin bool) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		email,
		isAdmin,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "easycache",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
