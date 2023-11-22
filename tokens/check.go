package tokens

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Check(signedString string) (Token jwt.MapClaims, Error error) {
	token, err := jwt.Parse(signedString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Unix(int64(claims["exp"].(float64)), 0)

	return claims, nil
}
