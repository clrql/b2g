package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    "session__auth",
		Value:   "",
		Expires: time.Now().Add(time.Second * 30),
	})
	c.ClearCookie("session_auth")
	return c.Status(200).JSON(fiber.Map{"message": "Logged-out successfully"})
}
