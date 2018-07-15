package main

import (
	"encoding/json"
	"fmt"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

var mongoClient MongoClient

// MongoOptions init options
type MongoOptions struct {
	uri string
	db string
}

// MongoClient class
type MongoClient struct {
	options MongoOptions
	client *mongo.Client
	db *mongo.Database
}

// NewMongo create MongoClient class
func NewMongo(mongoOptions MongoOptions) (MongoClient, error){
	var mongoClient MongoClient
	if mongoOptions.uri == "" {
		mongoOptions.uri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.Background(), mongoOptions.uri, nil)
	if err != nil {
		log.Error("connected mongodb failed")
		return mongoClient, err
	}

	mongoClient.client = client
	mongoClient.db = client.Database(mongoOptions.db)

	return mongoClient, nil
}

func (mongoClient *MongoClient) insertOne(collection string, jsonStr string) (*mongo.InsertOneResult, error){
	coll := mongoClient.db.Collection(collection)

	var f interface{}
    arr := []byte(jsonStr)
	err := json.Unmarshal(arr, &f)
	if err != nil {}

	result, err := coll.InsertOne(context.Background(), f)
	if err != nil {
		log.Error("insertOne error")
		fmt.Println(err, result)
	}

	return result, nil
}

func (mongoClient *MongoClient) insertMany(collection string, jsonStrs []string) (*mongo.InsertManyResult, error) {
	coll := mongoClient.db.Collection(collection)

	var fs []interface{}
	
	for _, jsonStr := range jsonStrs {
		var f interface{}
		arr := []byte(jsonStr)
		err := json.Unmarshal(arr, &f)
		if err != nil {}

		f = append(fs, f)
	}

	result, err := coll.InsertMany(context.Background(), fs)	
	if err != nil {
		log.Error("insertOne error")
		fmt.Println(err, result)
	}

	return result, nil
}

func (mongoClient *MongoClient) changeDB(dbName string) bool {
	mongoClient.db = mongoClient.client.Database(dbName)
	return true
}
