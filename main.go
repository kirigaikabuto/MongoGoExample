package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	Collection     *mongo.Collection
	Database       string = "lab3"
	CollectionName string = "posts"
)

type Post struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}

func (p *Post) Add() {
	insertResult, err := Collection.InsertOne(context.TODO(), p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}

func GetPostById(id int64) (*Post, error) {
	filter := bson.D{
		{"id", id},
	}
	post := &Post{}
	err := Collection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func CreateIndex() error {
	indexName, err := Collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"name": 1,
				"body": 1,
			},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		return err
	}
	fmt.Print("Create Index " + indexName)
	return nil
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	Collection = client.Database(Database).Collection(CollectionName)
	if err != nil {
		log.Fatal(err)
	}
	//err = CreateIndex()
	//if err != nil {
	//	log.Fatal(err)
	//}
	post1 := &Post{
		4,
		"231321321",
		"1232131233",
	}
	post1.Add()
	currentPost, err := GetPostById(1)
	fmt.Print(currentPost)
	defer client.Disconnect(context.TODO())
}
