package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// Page holds the main test struct
type Page struct {
	ID      bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Slug    string        `json:"slug" bson:"slug"`
	Name    string        `json:"name" bson:"name"`
	Content string        `json:"content" bson:"content"`
}

// HandleGetAllPages simply gets all pages from the database and returns all of them.
// It gives a http.StatusInternalServerError on db or json parsing issues.
func HandleGetAllPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pages := []*Page{}

	errDB := Session.DB(MongoDBDatabase).C("pages").Find(bson.M{}).All(&pages)
	if errDB != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`"%s"`, errDB.Error())))
		return
	}

	jsonData, errJSON := json.Marshal(pages)
	if errJSON != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`"%s"`, errJSON.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// HandleGetPage gets a page based on slug (from query parameter) and returns it.
// If the slug isn't set or when it isn't found in the DB we give http.StatusNotFound.
// It gives a http.StatusInternalServerError on json parsing issues.
func HandleGetPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := r.URL.Query().Get("slug")
	if slug == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`"Not found"`))
		return
	}

	page := &Page{}

	errDB := Session.DB(MongoDBDatabase).C("pages").Find(bson.M{"slug": slug}).One(page)
	if errDB != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`"Not found"`))
		return
	}

	jsonData, errJSON := json.Marshal(page)
	if errJSON != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`"%s"`, errJSON.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
