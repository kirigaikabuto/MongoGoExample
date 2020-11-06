package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)
//create global variables
var (
	Collection     *mongo.Collection
	Database       string = "lab3"
	CollectionName string = "posts"
)

//create post struct
type Post struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}
//create method to the struct of Post
func (p *Post) Add() {
	insertResult, err := Collection.InsertOne(context.TODO(), p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}

//create function for getting Post by id
func GetPostById(id int64) (*Post, error) {
	//create filter of type Document where exists condition by id
	filter := bson.D{
		{"id", id},
	}
	//create empty post instance in order to upload to this decoded data from mongodb
	post := &Post{}
	//find element  by the filter and set it to post variable
	err := Collection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func CreateIndex() error {
	//create one index for name and body fields and also set options unique to true for not duplication of data
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
	fmt.Println("Create Index " + indexName)
	return nil
}

func main() {
	//create options for connection to mongodb
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	//create connection to the mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	//select database and also collection
	Collection = client.Database(Database).Collection(CollectionName)
	//create index for name and body
	err = CreateIndex()
	if err != nil {
		log.Fatal(err)
	}
	//create post instance
	post1 := &Post{
		4,
		"231321321",
		"1232131233",
	}
	//run add method to add post to the mongodb
	post1.Add()
	//get post by the id
	currentPost, err := GetPostById(1)
	fmt.Println(currentPost)
	//close connection after all actions will be done
	defer client.Disconnect(context.TODO())
}
