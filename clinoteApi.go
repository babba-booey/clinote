package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var notesStoreUsername, notesStorePassword, notesStoreURI string = os.Getenv("NOTES_STORE_USERNAME"), os.Getenv("NOTES_STORE_PASSWORD"), os.Getenv("NOTES_STORE_URI")

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, options.Client().ApplyURI(notesStoreURI))

// Note describes a note that is to be saved to the database
type Note struct {
	NoteTitle string `json:"noteTitle"`
	NoteBody  string `json:"noteBody"`
}

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

func saveNote(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Info("Looking up note")
		w.WriteHeader(200)
	} else {
		var note Note
		bodyStr, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		json.Unmarshal(bodyStr, &note)
		log.Info("Received note. Title: ", note.NoteTitle, ", Body: ", note.NoteBody)
		w.WriteHeader(200)
	}
}

func registerEndpoints() {
	http.HandleFunc("/note", saveNote)
	log.Info("Registered endpoints")
}

func main() {
	initLogging()
	err := initDatabaseConnection()
	if err != nil {
		log.Error("Unable to start application because: ", err)
		os.Exit(2)
	}
	registerEndpoints()
	log.Info("Starting application")
	http.ListenAndServe(":9090", nil)
}
