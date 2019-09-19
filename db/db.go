package db

import (
	"context"
	"log"
	"time"

	"github.com/hanstar17/wordy/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var uri string
var connectTimeout time.Duration
var client *mongo.Client
var db *mongo.Database

func Init() {
	uri = env.GetString("DB_URI")
	connectTimeout = time.Duration(env.GetInt("DB_CONNECT_TIMEOUT_SEC")) * time.Second
}

func Connect() error {
	if client != nil {
		log.Fatalln("db is already connected.")
	}

	ctx, _ := context.WithTimeout(context.Background(), connectTimeout)

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	db = client.Database("wordthy")
	return nil
}

func Disconnect() error {
	if client == nil {
		log.Fatalln("db has not been connected.")
	}
	return client.Disconnect(context.TODO())
}
