package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const userCollectionKey = "users"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Eaddress string `json:"eaddress"`
}

func FindUser(filter bson.D) *User {
	coll := db.Collection(userCollectionKey)
	result := coll.FindOne(context.TODO(), filter)
	if result == nil {
		return nil
	}

	user := &User{}
	err := result.Decode(user)
	if err != nil {
		return nil
	}
	return user
}

func InsertUser(user User) primitive.ObjectID {
	coll := db.Collection(userCollectionKey)
	result, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid
	}
	return primitive.NilObjectID
}
