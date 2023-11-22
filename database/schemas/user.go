package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Birthday int64              `bson:"birthday"`
	Verified bool               `bson:"verified"`
}
