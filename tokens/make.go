package tokens

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var expTimes = map[string]time.Duration{
	"verify":         time.Minute * 5,
	"auth":           time.Minute * 30,
	"reset-password": time.Hour * 30,
}

func Make(subject string, ttype string) (SignedString string, Error error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["jti"] = 100000 + rand.Intn(900000)
	claims["sub"] = subject

	exptime, ok := expTimes[ttype]
	if !ok {
		return "", fmt.Errorf("Invalid token type")
	}

	claims["exp"] = time.Now().Add(exptime).Unix()
	claims["type"] = ttype

	signedString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedString, nil
}
