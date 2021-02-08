package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	SECRETKEY = "243223ffslsfsldfl412fdsfsdf"
)

type CustomClaims struct {
	UserId int64
	jwt.StandardClaims
}

// Generate token
func GenerateToken(userId int64, email string) string {
	// seconds of one day
	maxAge := 60 * 60 * 24
	customClaims := &CustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(maxAge) * time.Second).Unix(),
			Issuer:    email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, _ := token.SignedString([]byte(SECRETKEY))
	return tokenString
}

// Parse token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
