package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "html/template"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
  )

func RandomURLHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RandomURLHandler] ")
  log.Println("hallo from the random URL handler")
  vars := mux.Vars(r)
  name := vars["name"]
  t, err := template.ParseFiles("assets/templates/index.html.tmpl")
  if err != nil{
    //deal with 500s later
    log.Fatal("this is a problem")
    log.Fatal(err)
  } else {
    t.Execute(w, map[string] string {"Name": name})
  }
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RootHandler] ")
  log.Println("hallo from the root handler")
  fmt.Fprintln(w, "Welcome to restaur-anteater")
}

func main() {
  log.SetPrefix("[main] ")
  log.Println("hallo")

  //Read the server port as an ENV variable
  var port string
  if port = os.Getenv("PORT"); port == ""{
    port = "8000"
  }
  log.Println("Using port: " + port)

  //Setup the routes
  router := mux.NewRouter()
  router.HandleFunc("/", RootHandler)
  router.HandleFunc("/test/{name}", RandomURLHandler)

  //Tell the router to server the assets folder as static files
  router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

  //Setup the middleware
  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(router)

  //Run the server
  log.Fatal(http.ListenAndServe(":"+port, n))
}
