package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
  )

func RandomURLHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RandomURLHandler] ")
  log.Println("hallo from the random URL handler")
  vars := mux.Vars(r)
  name := vars["name"]
  fmt.Fprintf(w, "Hello, %s", name)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RootHandler] ")
  log.Println("hallo from the root handler")
  fmt.Fprintln(w, "Welcome to restaur-anteater")
}

func main() {
  log.SetPrefix("[main] ")
  log.Println("hallo")
  var port string
  if port = os.Getenv("PORT"); port == ""{
    port = "8000"
  }
  log.Println("Using port: " + port)
  router := mux.NewRouter()
  router.HandleFunc("/", RootHandler)
  router.HandleFunc("/test/{name}", RandomURLHandler)


  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(router)
  log.Fatal(http.ListenAndServe(":"+port, n))
}
