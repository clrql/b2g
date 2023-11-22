package middlewares

import (
	"b2g/tokens"

	"github.com/gofiber/fiber/v2"
)

func Authed(c *fiber.Ctx) error {

	session__auth := c.Cookies("session__auth")
	if session__auth == "" {
		return c.Redirect("/auth/login")
	}

	token, err := tokens.Check(session__auth)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Your session is not valid"})
	}

	c.Locals("auth_token", token)

	return c.Next()
}
