package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestHandleGetPage does a call to the HandleGetPage with the pages["ding"] (from fixtures) slug as query param.
// This should give a 200 as this is inserted already and should be parsed as JSON. Also the name we get from
// the JSON should the same as the one from pages["ding"]
func TestHandleGetPage(t *testing.T) {
	// Get the page with slug "ding" from fixtures so we can easily check it.
	page := pages["ding"]

	// Create the requestURI
	requestURI := url.URL{
		Path: "/page",
	}
	q := requestURI.Query()
	q.Set("slug", page.Slug)
	requestURI.RawQuery = q.Encode()

	// Create a new HTTP Get request with the above generated uri.
	req, err := http.NewRequest("GET", requestURI.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to get results.
	rr := httptest.NewRecorder()
	// Set the handler to test (HandleGetPage)
	handler := http.HandlerFunc(HandleGetPage)
	// Run the above request
	handler.ServeHTTP(rr, req)

	// Status should be http.StatusOK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Let's just parse in map[string]interface{} for now.
	testResults := make(map[string]interface{})

	// The response body should be Unmarshaled
	errJSON := json.Unmarshal(rr.Body.Bytes(), &testResults)
	if errJSON != nil {
		t.Errorf("Error unmarshalling JSON %s",
			errJSON.Error())
	}

	// The name element should be set
	if _, ok := testResults["name"]; !ok {
		t.Fatalf(`Expected element "name", didn't get it: "%s"`, rr.Body.String())
	}

	// The name element should be the same as pages["ding"].Name
	if testResults["name"] != page.Name {
		t.Errorf(`Expected name "%s", got "%s"`, page.Name, testResults["name"])
	}
}

// TestHandleGetPage404 does a call to the HandleGetPage with a nonexisting slug as query param.
// This should return a http.StatusNotFound
func TestHandleGetPage404(t *testing.T) {
	// Create the requestURI
	requestURI := url.URL{
		Path: "/page",
	}
	q := requestURI.Query()
	q.Set("slug", "nonexisting")
	requestURI.RawQuery = q.Encode()

	// Create a new HTTP Get request with the above generated uri.
	req, err := http.NewRequest("GET", requestURI.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to get results.
	rr := httptest.NewRecorder()
	// Set the handler to test (HandleGetPage)
	handler := http.HandlerFunc(HandleGetPage)
	// Run the above request
	handler.ServeHTTP(rr, req)

	// Status should be http.StatusNotFound
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestHandleGetPageNoSlug does a call to the HandleGetPage without a slug as query param.
// This should return a http.StatusNotFound
func TestHandleGetPageNoSlug(t *testing.T) {
	// Create the requestURI
	requestURI := url.URL{
		Path: "/page",
	}

	// Create a new HTTP Get request with the above generated uri.
	req, err := http.NewRequest("GET", requestURI.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to get results.
	rr := httptest.NewRecorder()
	// Set the handler to test (HandleGetPage)
	handler := http.HandlerFunc(HandleGetPage)
	// Run the above request
	handler.ServeHTTP(rr, req)

	// Status should be http.StatusNotFound
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestHandleGetAllPages does a call to the HandleGetAllPages.
// This should return a http.StatusOK and should return len(pages) pages.
func TestHandleGetAllPages(t *testing.T) {
	// Create the requestURI
	requestURI := url.URL{
		Path: "/",
	}

	// Create a new HTTP Get request with the above generated uri.
	req, err := http.NewRequest("GET", requestURI.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to get results.
	rr := httptest.NewRecorder()
	// Set the handler to test (HandleGetAllPages)
	handler := http.HandlerFunc(HandleGetAllPages)
	// Run the above request
	handler.ServeHTTP(rr, req)

	// Status should be http.StatusOK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Let's try to unmarshal body to []*Page{}, but I am currently only really interested in the amount of results
	testResults := []*Page{}

	// The response body should be Unmarshaled
	errJSON := json.Unmarshal(rr.Body.Bytes(), &testResults)
	if errJSON != nil {
		t.Fatalf("Error unmarshalling JSON %s",
			errJSON.Error())
	}

	// The amount of returned pages should be the same as the amount of page fixtures.
	if len(testResults) != len(pages) {
		t.Errorf(`Expected "%d" results, got "%d"`, len(pages), len(testResults))
	}
}
