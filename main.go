package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

// Session is the connection to MongoDB
var Session *mgo.Session

// MongoDBHostname is the connection string of MongoDB
var MongoDBHostname = "localhost:27017"

// MongoDBDatabase is the used database
var MongoDBDatabase = "tmp_test_db"

// Port is used to determine what port to listen on
var Port = "3000"

// init sets basic config from env and start MongoDB Session
func init() {
	if dbHostname := os.Getenv("MONGODB_HOSTNAME"); dbHostname != "" {
		MongoDBHostname = dbHostname
	}

	if dbName := os.Getenv("MONGODB_DATABASE"); dbName != "" {
		MongoDBDatabase = dbName
	}

	if port := os.Getenv("PORT"); port != "" {
		Port = port
	}

	var err error

	Session, err = mgo.Dial(MongoDBHostname)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// main inits http handlers and starts listener
func main() {

	http.HandleFunc("/page", HandleGetPage)
	http.HandleFunc("/", HandleGetAllPages)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", Port), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
