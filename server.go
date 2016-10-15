package main

import (
    // "fmt"
    "log"
    // "strconv"
    "net/http"
    "os"
    // "time"
    "html/template"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    "github.com/jinzhu/gorm"
    "golang.org/x/crypto/bcrypt"
    // "github.com/gorilla/sessions"
    "math/rand"
    _ "github.com/jinzhu/gorm/dialects/postgres"
  )

var db *gorm.DB

//Leaving this so people can still see the code but it won't work anymore with the Devices table removed.

// func CreateDeviceHandler(w http.ResponseWriter, r *http.Request) {
//   log.SetPrefix("[CreateDeviceHandler] ")
//   vars := mux.Vars(r)
//   CustomerID, _ := strconv.Atoi(vars["customer_id"])
//   PartySize, _ := strconv.Atoi(vars["party_size"])
//   device := Device{CustomerID: CustomerID, DeviceName: vars["device_name"], IsActive: true, PartySize: PartySize}
//   db.NewRecord(device)
//   db.Create(&device)
//   fmt.Fprintln(w, "Device created!")
// }

// func FindDevicesHandler(w http.ResponseWriter, r *http.Request) {
//   log.SetPrefix("[DisplayDeviceHandler] ")
//   vars := mux.Vars(r)
//   CustomerID, _ := strconv.Atoi(vars["customer_id"])
//   var devices []Device
//   db.Where("customer_id = ?", CustomerID).Find(&devices)
//   for _, device := range devices {
//     log.Println("Device name: " + device.DeviceName)
//   }
//   t, err := template.ParseFiles("assets/templates/find_device.html.tmpl")
//   if err != nil{
//     log.Println(err)
//     log.Fatal("fail")
//   } else {
//     //Use an anonymous struct to pass data to the template.
//     data := struct {
//       CustomerID int
//       Devices []Device
//     }{
//       CustomerID,
//       devices,
//     }
//     t.Execute(w, data)
//   }
// }

func makeRandAlphaNumericStr(n int) string {
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
  b := make([]rune, n)
  for i := range b {
      b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RootHandler] ")
  log.Println("hallo from the root handler")
  t, err := template.ParseFiles("assets/templates/loginpage.html.tmpl")
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
  log.Println(r.Method)
  log.Println(makeRandAlphaNumericStr(50))
  if r.Method == "POST" {
    username := r.FormValue("username")
    password := r.FormValue("password")
    restaurantName := r.FormValue("restaurant_name")
    if username != "" && password != "" && restaurantName != "" {
      passSalt := makeRandAlphaNumericStr(50)
      hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+passSalt), bcrypt.DefaultCost)
      if err != nil {
        panic(err)
      }
      log.Println(string(hashedPassword))
      log.Println(restaurantName)
      var restaurant Restaurant;
      db.First(&restaurant, "name = ?", restaurantName)
      if restaurant == (Restaurant{}) {
        restaurant = Restaurant{Name: restaurantName}
        db.NewRecord(restaurant)
        db.Create(&restaurant)
      }
      user := User{RestaurantID: restaurant.ID, Username: username, Password: string(hashedPassword), PassSalt: passSalt}
      db.NewRecord(user)
      db.Create(&user)
    }
  }
  t, err := template.ParseFiles("assets/templates/adduser.html.tmpl")
  if err != nil{
    //deal with 500s later
    log.Println("this is a problem")
    log.Fatal(err)
  } else {
    t.Execute(w, nil)
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
  router.HandleFunc("/add_user", AddUserHandler)
  router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

  //Tell the router to server the assets folder as static files
  router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

  //Setup the middleware
  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(router)

  //Run the server
  log.Fatal(http.ListenAndServe(":"+port, n))
}
