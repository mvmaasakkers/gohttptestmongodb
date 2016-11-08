package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"gopkg.in/mgo.v2/dbtest"
)

// Server holds the dbtest DBServer
var Server dbtest.DBServer

// Pages fixtures are intentionally setup as map[string]Page so I can easily select them from within the tests
var pages = map[string]Page{
	"ding":      Page{Slug: "ding", Name: "Ding!", Content: "<i>HTML Awesomeness</i>"},
	"dong-ding": Page{Slug: "dong-ding", Name: "Dong! Ding!", Content: "<b>HTML Awesomeness</b>"},
}

// insertFixtures just inserts all pages (and other types) I've defined above.
func insertFixtures() {
	for _, page := range pages {
		if err := Session.DB(MongoDBDatabase).C("pages").Insert(page); err != nil {
			log.Println(err)
		}
	}
}

// reInsertFixtures drops database and re-inserts all fixtures so we can
// make sure every test can start fresh.
func reInsertFixtures() {
	Session.DB(MongoDBDatabase).DropDatabase()
	insertFixtures()
}

// TestMain wraps all tests with the needed initialized mock DB and fixtures
func TestMain(m *testing.M) {
	// The tempdir is created so MongoDB has a location to store its files.
	// Contents are wiped once the server stops
	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	// My main session var is now set to the temporary MongoDB instance
	Session = Server.Session()

	// Make sure to insert my fixtures
	insertFixtures()

	// Run the test suite
	retCode := m.Run()

	// Make sure we DropDatabase so we make absolutely sure nothing is left or locked while wiping the data and
	// close session
	Session.DB(MongoDBDatabase).DropDatabase()
	Session.Close()

	// Stop shuts down the temporary server and removes data on disk.
	Server.Stop()

	// call with result of m.Run()
	os.Exit(retCode)
}
