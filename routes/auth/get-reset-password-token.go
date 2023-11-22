package auth

import (
	"b2g/database"
	"b2g/database/schemas"
	"b2g/tokens"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetResetPasswordToken(c *fiber.Ctx) error {
	var query struct {
		Username string `query:"username"`
	}

	if err := c.QueryParser(&query); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Failed to parse query"})
	}

	var user schemas.User
	db, err := database.Conn()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Failed to check if the user exists I try again later"})
	}

	coll := db.Collection("user")
	if err := coll.FindOne(context.TODO(), bson.M{"username": query.Username}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"message": "There's not an user with that username check it and try again"})
		}
		return c.Status(503).JSON(fiber.Map{"message": "Failed to check if the user exists II try again later"})
	}

	signedString, err := tokens.Make(query.Username, "reset-password")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to generate reset token"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "use the following password-reset-link to reset your password", "password-reset-link": "http://localhost:9090/auth/reset-password?t=" + signedString})
}
