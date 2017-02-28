package main

import (
	"testing"
    "net/http"
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