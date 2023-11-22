package auth

import (
	"b2g/database"
	"b2g/database/schemas"
	"b2g/tokens"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	var body struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Failed to parse request body"})
	}

	db, err := database.Conn()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Login not available try again later"})
	}

	coll := db.Collection("user")

	var user schemas.User
	if err := coll.FindOne(context.TODO(), bson.M{"$or": []bson.M{
		{"username": body.User},
		{"email": body.User},
	}}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"message": "There's not an user with that username/email check your credentials"})
		}
		return c.Status(503).JSON(fiber.Map{"message": "Login is not working try again later"})
	}

	if !user.Verified {
		return c.Status(403).JSON(fiber.Map{"message": "Your account need to be verified before login in"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Your password is not correct check it and try again"})
	}

	signedString, err := tokens.Make(user.Username, "auth")
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Server was not able to init your session try again later"})
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session__auth",
		Value:   signedString,
		Expires: time.Now().Add(time.Hour * 24),
	})

	return c.Status(200).JSON(fiber.Map{"message": "Logged successfully! welcome " + user.Name})
}
