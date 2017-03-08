package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
)

func TestIsLoggedIn(t *testing.T) {
	var session *sessions.Session
	var err error

	t.Log("Start with no session values")

	store := sessions.NewCookieStore([]byte("secret-key"))

	req, _ := http.NewRequest("GET", "http://localhost:8000/", nil)
	if session, err = store.Get(req, "session-key"); err != nil {
		t.Fatalf("Error getting session: %v", err)
	}

	if IsUserLoggedIn(session) {
		t.Error("There should have been no isuserloggedin")
	}

	session.Values["username"] = "something"

	if !IsUserLoggedIn(session) {
		t.Error("There is a login value, so it should be logged in")
	}
}

func TestDoesHandlerWork(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// check if logged in
	var session *sessions.Session
	var err error

	t.Log("Start with no session values")

	store := sessions.NewCookieStore([]byte("secret-key"))
	if session, err = store.Get(req, "session-key"); err != nil {
		t.Fatalf("Error getting session: %v", err)
	}

	// userhandler should add "Could not add user. Did you forget a field? to flash
	AddUserHandler(w, req)

}
