package auth

import (
	"b2g/database"
	"b2g/database/schemas"
	"b2g/tokens"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday int64  `json:"birthday"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Failed to parse request body"})
	}

	db, err := database.Conn()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Register not available try agan later"})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to hash the password try with a different password"})
	}

	var newuser schemas.User
	newuser.Id = primitive.NewObjectID()
	newuser.Name = body.Name
	newuser.Email = body.Email
	newuser.Password = string(passwordHash)
	newuser.Username = body.Username
	newuser.Birthday = body.Birthday
	newuser.Verified = false

	coll := db.Collection("user")

	var existingUser schemas.User
	if err := coll.FindOne(context.TODO(), bson.M{"$or": []bson.M{
		{"username": newuser.Username},
		{"email": newuser.Email},
	}}).Decode(&existingUser); err != nil {
		if err != mongo.ErrNoDocuments {
			return c.Status(503).JSON(fiber.Map{"message": "Failed to check if the username and email is not in use try again later"})
		}
	}

	if existingUser.Id != primitive.NilObjectID {
		return c.Status(409).JSON(fiber.Map{"message": "There's an user with that username or email change the conflicting field"})
	}

	_, err = coll.InsertOne(context.TODO(), newuser)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"message": "Server failed to register you, try again later"})
	}

	SignedString, _ := tokens.Make(newuser.Username, "verify")

	return c.Status(200).JSON(fiber.Map{"message": "you were registered successfully use the following verify-link to verify your account", "verify-link": "http://localhost:9090/auth/verify?t=" + SignedString})
}
