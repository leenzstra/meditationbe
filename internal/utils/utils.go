package utils

import (
	"fmt"
	"meditationbe/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func TokenFromLocals(c *fiber.Ctx) (*jwt.Token, error) {
	token := c.Locals("user").(*jwt.Token)
	if token == nil {
		return nil, fmt.Errorf("token not found")
	}

	return token, nil
}

func TokenFromHeaders(c *fiber.Ctx) (*jwt.Token, error) {
	authHeader, ok := c.GetReqHeaders()["Authorization"]
	if !ok {
		return nil, fmt.Errorf("no Authorization header")
	}

	header := strings.Split(authHeader[0], " ")
	if len(header) != 2 {
		return nil, fmt.Errorf("invalid Authorization header")
	}	
	
	tokenString := header[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.GetConfig().JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func GenerateJWT(claims jwt.MapClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}