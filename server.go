package main

import (
    // "fmt"
    "log"
    // "strconv"
    "net/http"
    "os"
    "time"
    "html/template"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    "github.com/jinzhu/gorm"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/sessions"
    "math/rand"
    _ "github.com/jinzhu/gorm/dialects/postgres"
  )

var db *gorm.DB
var sessionStore *sessions.CookieStore

func MakeRandAlphaNumericStr(n int) string {
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
  b := make([]rune, n)
  for i := range b {
      b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}

func Handle500Error(w http.ResponseWriter) {
  http.Error(w, http.StatusText(500), 500)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[LoginURLHandler] ")
  session, err := sessionStore.Get(r, "buzzer-session")
  if err != nil {
    log.Println(err)
    Handle500Error(w)
  }
  if r.Method == "POST" {
    username := r.FormValue("username")
    password := r.FormValue("password")

    var user User;
    db.First(&user, "Username = ?", username)
    if (user != (User{})) {
      passSalt := user.PassSalt
      if (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+passSalt)) == nil) {
          session.Values["username"] = username
          session.Save(r, w)
          http.Redirect(w, r, "/wait_list", 302)
      }
    } else {
      session.AddFlash("Username or Password not correct")
      session.Save(r, w)
    }
  }
  t, err := template.ParseFiles("assets/templates/login.html.tmpl")
  if err != nil{
    log.Println(err)
    Handle500Error(w)
  }
  t.Execute(w, nil)
}

func WaitListHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[WaitListHandler] ")
  session, err := sessionStore.Get(r, "buzzer-session")
  if err != nil {
    log.Fatal("this is a problem")
  }
  username, found := session.Values["username"]
  if !found || username == "" {
    http.Redirect(w, r, "/login", 302)
    log.Println("User not logged in")
    return
  }
  t, err := template.ParseFiles("assets/templates/waitlist.html.tmpl")
  if err != nil{
    log.Println(err)
    Handle500Error(w)
  } else {
    t.Execute(w, nil)
  }
}


func RootHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RootHandler] ")
  log.Println("hallo from the root handler")
  t, err := template.ParseFiles("assets/templates/login.html.tmpl")
  if err != nil{
    //deal with 500s later
    log.Println("this is a problem")
    log.Fatal(err)
  } else {
    t.Execute(w, nil)
  }
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[AddUserHandler] ")
  session, err := sessionStore.Get(r, "buzzer-session")
  if err != nil {
    log.Fatal("this is a problem")
  }
  if r.Method == "POST" {
    username := r.FormValue("username")
    password := r.FormValue("password")
    restaurantName := r.FormValue("restaurant_name")
    if username != "" && password != "" && restaurantName != "" {

      //salt and hash the password
      passSalt := MakeRandAlphaNumericStr(50)
      hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+passSalt), bcrypt.DefaultCost)
      if err != nil {
        log.Fatal(err)
      }

      var restaurant Restaurant;
      db.First(&restaurant, "name = ?", restaurantName)

      //make a restaurant if there isn't one
      if restaurant == (Restaurant{}) {
        restaurant = Restaurant{Name: restaurantName}
        db.NewRecord(restaurant)
        db.Create(&restaurant)
      }

      var user User
      db.First(&user, "username = ?", username)

      if user != (User{}) {
        session.AddFlash("Username already exists")
      } else {
        //add the user
        user = User{RestaurantID: restaurant.ID, Username: username, Password: string(hashedPassword), PassSalt: passSalt}
        db.NewRecord(user)
        db.Create(&user)
        session.AddFlash("User successfully added")
      }
    } else {
      session.AddFlash("Could not add user. Did you forget a field?")
    }
    session.Save(r, w)
    http.Redirect(w, r, "/add_user", 302)
  } else {
    template_data := map[string]interface{}{}
    if flashes := session.Flashes(); len(flashes) > 0 {
      template_data["flash"] = flashes[0]
    }
    session.Save(r, w)
    t, err := template.ParseFiles("assets/templates/adduser.html.tmpl")
    if err != nil{
      //deal with 500s later
      log.Println("this is a problem")
      log.Fatal(err)
    } else {
      t.Execute(w, template_data)
    }
  }
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[NotFoundHandler] ")
  log.Println("hit the not found handler")
  t, _ := template.ParseFiles("assets/templates/404.html")
  t.Execute(w, nil)
}

func main() {
  log.SetPrefix("[main] ")
  rand.Seed(time.Now().UnixNano())

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
    panic("error connecting to db")
  }

  //Setup the routes
  router := mux.NewRouter()
  router.HandleFunc("/", RootHandler)
  router.HandleFunc("/login", LoginHandler)
  router.HandleFunc("/add_user", AddUserHandler)
  router.HandleFunc("/wait_list", WaitListHandler)
  router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

  //Tell the router to server the assets folder as static files
  router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

  //Setup the middleware
  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(router)

  //Run the server
  log.Fatal(http.ListenAndServe(":"+port, n))
}
