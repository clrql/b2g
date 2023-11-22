package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// This endpoint use the MiddlewareAuthed
func Authed(c *fiber.Ctx) error {
	token := c.Locals("auth_token").(jwt.MapClaims)
	return c.Status(200).JSON(fiber.Map{"message": "authed", "username": token["sub"]})
}
