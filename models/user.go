package models

import (
	"crud-app/database"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	// ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ID    string `bson:"id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Age   int    `bson:"age" json:"age"`
	// IsDeleted bool       `bson:"isDeleted" json:"isDeleted"`
	CreatedAt time.Time  `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time  `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	DeletedAt *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

func GetUserCollection() *mongo.Collection {
	if database.MongoDatabase == nil {
		panic("MongoDatabase not initialized. Call database.Connect() first.")
	}
	return database.MongoDatabase.Collection("users")
}
