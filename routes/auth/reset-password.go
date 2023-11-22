package auth

import (
	"b2g/database"
	"b2g/database/schemas"
	"b2g/tokens"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(c *fiber.Ctx) error {
	var body struct {
		NewPassword string `json:"newPassword"`
		Token       string `json:"token"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Failed to parse body"})
	}

	token, err := tokens.Check(body.Token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Password reset token is not valid or is expired"})
	}

	db, err := database.Conn()
	coll := db.Collection("user")

	var user schemas.User
	if err := coll.FindOne(context.TODO(), bson.M{"username": token["sub"]}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"message": "The token is not valid there's not an user with that username"})
		}
		return c.Status(503).JSON(fiber.Map{"message": "Failed to reset the password try again later"})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 12)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to parse the new password try with a different password"})
	}

	_, err = coll.UpdateByID(context.TODO(), user.Id, bson.M{"$set": bson.M{"password": passwordHash}})
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Failed to save the new password try again later"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Your new password is set, try to login"})
}
