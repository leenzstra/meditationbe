package middleware

import (
	"meditationbe/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return jwtware.New(
		jwtware.Config{
			ContextKey:  "user",
			SigningKey:  jwtware.SigningKey{Key: []byte(config.GetConfig().JWTSecret)},
			TokenLookup: "header:Authorization",
			AuthScheme:  "Bearer",
		},
	)
}
