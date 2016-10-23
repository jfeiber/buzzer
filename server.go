package main

import (
    "log"
    "net/http"
    "os"
    "time"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    "github.com/jinzhu/gorm"
    "github.com/gorilla/sessions"
    "math/rand"
    _ "github.com/jinzhu/gorm/dialects/postgres"
  )

var db *gorm.DB
var sessionStore *sessions.CookieStore
var buzzerNameGenerator BuzzerNameGenerator

func main() {
  log.SetPrefix("[main] ")
  rand.Seed(time.Now().UnixNano())

  buzzerNameGenerator = NewBuzzerNameGenerator(time.Now().UnixNano())

  authKey := os.Getenv("SESSION_AUTHENTICATION_KEY")
  if authKey == "" {
    log.Fatal("$SESSION_AUTHENTICATION_KEY needs to be set")
  }
  sessionStore = sessions.NewCookieStore([]byte(authKey))

  //Read the server port as an ENV variable
  var port string
  if port = os.Getenv("PORT"); port == ""{
    port = "8000"
  }
  log.Println("Using port: " + port)

  var database_url string
  if database_url = os.Getenv("DATABASE_URL"); database_url == ""{
    database_url = "host="+os.Getenv("POSTGRES_PORT_5432_TCP_ADDR")+" user=ra dbname=ra password=password sslmode=disable"
  }
  log.Println("Using DB URL: " + database_url)

  var err error
  db, err = gorm.Open("postgres", database_url)
  defer db.Close()
  db.LogMode(true)
  if err != nil{
    log.Fatal(err)
  }

  //Setup the routes
  router := mux.NewRouter()
  router.HandleFunc("/", RootHandler)
  router.HandleFunc("/login", LoginHandler)
  router.HandleFunc("/add_user", AddUserHandler)
  router.HandleFunc("/wait_list", WaitListHandler)
  router.HandleFunc("/wait_temp", WaitListTempHandler)
  router.HandleFunc("/buzzer_api/get_new_buzzer_name", GetNewBuzzerNameHandler)
  router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

  //Tell the router to server the assets folder as static files
  router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

  //Setup the middleware
  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(router)

  //Run the server
  log.Fatal(http.ListenAndServe(":"+port, n))
}
