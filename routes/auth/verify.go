package auth

import (
	"b2g/database"
	"b2g/database/schemas"
	"b2g/tokens"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Verify(c *fiber.Ctx) error {
	var query struct {
		T string `query:"t"`
	}

	if err := c.QueryParser(&query); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Failed to parse request url query"})
	}

	token, err := tokens.Check(query.T)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"message": "Query field \"t\" is not a valid json web token"})
	}

	if token["type"] != "verify" {
		return c.Status(400).JSON(fiber.Map{"message": "The token is valid but is not a verify token"})
	}

	db, err := database.Conn()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Verify is not available try agan later"})
	}

	coll := db.Collection("user")

	var user schemas.User

	if err := coll.FindOne(context.TODO(), bson.M{"username": token["sub"]}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"message": "The user doesn't exists"})
		}
		return c.Status(503).JSON(fiber.Map{"message": "Failed to check if the user is already verified try again later"})
	}

	if user.Verified {
		return c.Status(409).JSON(fiber.Map{"message": "You're already verified"})
	}

	_, err = coll.UpdateByID(context.TODO(), user.Id, bson.M{"$set": bson.M{"verified": true}})
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "The server was not able to verify you try again later"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Your account was verified successfully"})
}
