package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
  )

func RootSubdomainHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RootSubdomainHandler] ")
  log.Println("hallo from the root subdomain handler")
  vars := mux.Vars(r)
  subdomain := vars["subdomain"]
  fmt.Fprintf(w, "Hello, %s", subdomain)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RootHandler] ")
  log.Println("hallo from the root handler")
  fmt.Fprintln(w, "Welcome to restaur-anteater")
}

func main() {
  log.SetPrefix("[main] ")
  log.Println("hallo")
  var domain, port string
  if domain = os.Getenv("DOMAIN"); domain == ""{
    domain = "localhost"
  }
  if port = os.Getenv("PORT"); port == ""{
    port = "8000"
  }
  log.Println("Using domain: " + domain)
  log.Println("Using port: " + port)
  mux := mux.NewRouter()
  subdomain_match := "{subdomain:[a-z0-9 -]+}." + domain
  log.Println("Using subdomain match: " + subdomain_match)

  mux.Host(domain).Path("/").HandlerFunc(RootHandler)
  mux.Host(subdomain_match).Path("/").HandlerFunc(RootSubdomainHandler)

  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(mux)
  log.Fatal(http.ListenAndServe(":"+port, n))
}

