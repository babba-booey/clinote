package main

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var notesStoreUsername, notesStorePassword, notesStoreURI string = os.Getenv("NOTES_STORE_USERNAME"), os.Getenv("NOTES_STORE_PASSWORD"), os.Getenv("NOTES_STORE_URI")

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

func initLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.Info("Logging initialized")
}

func initDatabaseConnection() error {
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Warn("Unable to connect to MongoDB. Reason: ", err)
		return err
	}
	log.Info("MongoDB connected")
	return nil
}

func main() {
	initLogging()
	err := initDatabaseConnection()
	if err != nil {
		log.Error("Unable to start application because: ", err)
	} else {
		log.Info("Starting application")
	}
}
