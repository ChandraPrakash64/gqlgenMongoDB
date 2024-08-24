package database

import (
	"context"
	"fmt"
	"log"

	"example.com/gqlmongo/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

const mongoURI = "mongodb://localhost:27017"

func Connect() *DB {
	clientOption := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	db := &DB{
		client: client,
	}
	return db
}

func (db *DB) CreateUser(name *string, email *string) *model.User {
	coll := db.client.Database("mydbs").Collection("user")
	returnedId, err := coll.InsertOne(context.Background(), bson.M{"name": name, "email": email})
	if err != nil {
		fmt.Println(err)
	}
	insertedId := returnedId.InsertedID.(primitive.ObjectID).Hex()
	newUser := model.User{
		ID:    insertedId,
		Name:  name,
		Email: email,
	}
	return &newUser
}

func (db *DB) GetUser(id string) *model.User {
	_id, _ := primitive.ObjectIDFromHex(id)
	coll := db.client.Database("mydbs").Collection("user")
	filter := bson.M{"_id": _id}
	var reqUser model.User
	err := coll.FindOne(context.Background(), filter).Decode(reqUser)
	if err != nil {
		log.Fatal(err)
	}
	return &reqUser
}

func (db *DB) GetUsers() []*model.User {
	coll := db.client.Database("mydbs").Collection("user")
	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var users []*model.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		log.Fatal(err)
	}
	return users
}
