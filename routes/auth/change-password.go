package auth

import (
	"b2g/database"
	"b2g/database/schemas"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// This endpoint use the MiddlewareAuthed
func ChangePassword(c *fiber.Ctx) error {
	token := c.Locals("auth_token").(jwt.MapClaims)

	var body struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Failed to parse body"})
	}

	db, err := database.Conn()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Password change is not available try again later"})
	}

	coll := db.Collection("user")

	var user schemas.User
	if err := coll.FindOne(context.TODO(), bson.M{"username": token["sub"]}).Decode(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Unexpected error happened try again later"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Old password wasn't correct check it and try again"})
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 12)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Failed to hash the new password try with a different password"})
	}

	_, err = coll.UpdateByID(context.TODO(), user.Id, bson.M{"$set": bson.M{"password": string(newPasswordHash)}})
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Failed to save the new password try again later"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "New password was saved succesfully"})
}
