package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "html/template"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
  )

var db *gorm.DB

func CreateDeviceHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[CreateDeviceHandler] ")
  vars := mux.Vars(r)
  CustomerID, _ := strconv.Atoi(vars["customer_id"])
  PartySize, _ := strconv.Atoi(vars["party_size"])
  device := Device{CustomerID: CustomerID, DeviceName: vars["device_name"], IsActive: true, PartySize: PartySize}
  db.NewRecord(device)
  db.Create(&device)
  fmt.Fprintln(w, "Device created!")
}

func FindDevicesHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[DisplayDeviceHandler] ")
  vars := mux.Vars(r)
  CustomerID, _ := strconv.Atoi(vars["customer_id"])
  var devices []Device
  db.Where("customer_id = ?", CustomerID).Find(&devices)
  for _, device := range devices {
    log.Println("Device name: " + device.DeviceName)
  }
  t, err := template.ParseFiles("assets/templates/find_device.html.tmpl")
  if err != nil{
    log.Println(err)
    log.Fatal("fail")
  } else {
    //Use an anonymous struct to pass data to the template.
    data := struct {
      CustomerID int
      Devices []Device
    }{
      CustomerID,
      devices,
    }
    t.Execute(w, data)
  }
}

func RandomURLHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RandomURLHandler] ")
  log.Println("hallo from the random URL handler")
  vars := mux.Vars(r)
  name := vars["name"]
  t, err := template.ParseFiles("assets/templates/index.html.tmpl")
  if err != nil{
    //deal with 500s later
    log.Println("this is a problem")
    log.Fatal(err)
  } else {
    t.Execute(w, map[string] string {"Name": name})
  }
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

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[NotFoundHandler] ")
  log.Println("hit the not found handler")
  t, _ := template.ParseFiles("assets/templates/404.html")
  t.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

    // log.Println(r)

    err := r.ParseForm()
    if err != nil{
           panic(err)
    }

    params := r.PostFormValue("params")
    log.Println(params)

    // type WebAppUser struct {
    //     WebAppUserID int
    //     restaurant Restaurant
    //     restaurantId int
    //     username string `gorm:"size:100;not null"`
    //     password string `gorm:"size:100; not null"`
    //     passSalt string `gorm:"size:50; not null"`
    //     dateCreated time.Time
    // }
    // webAppUser := WebAppUser{
	// 	NotifyType:     notifyType,
	// 	UserId:         userId,
	// 	ActorId:        actorId,
	// 	NotifyableType: notifyableType,
	// 	NotifyableId:   notifyableId,
	// }
    //
	// exitCount := 0
	// db.Model(Notification{}).Where(
	// 	"user_id = ? and actor_id = ? and notifyable_type = ? and notifyable_id = ?",
	// 	userId, actorId, notifyableType, notifyableId).Count(&exitCount)
	// if exitCount > 0 {
	// 	return nil
	// }
    //
	// err := db.Save(&note).Error
    //
	// go PushNotifyInfoToUser(userId, note, true)
    //
	// return err
    //
    //


	// err := db.Save(n).Error
	// if err != nil {
	// 	v.Error("服务器异常创建失败")
	// }

    //redirect
    //

}

func LoggedInHandler(w http.ResponseWriter, r *http.Request) {
    log.SetPrefix("[dickinButt]")
    log.Println("Logged in bitch!!!!")

    // TODO: If logged in then do something else

    t, err := template.ParseFiles("assets/templates/loggedInPage.html.tmpl")
    if err != nil{
      //deal with 500s later
      log.Println("this is a problem")
      log.Fatal(err)
    } else {
      t.Execute(w, nil)
    }
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
  if err != nil{
    log.Fatal(err)
    panic("error connecting to db")
  }

  //Setup the routes
  router := mux.NewRouter()
  router.HandleFunc("/", RootHandler)
  router.HandleFunc("/test/{name}", RandomURLHandler)
  router.HandleFunc("/create_device/{customer_id}/{device_name}/{party_size}", CreateDeviceHandler)
  router.HandleFunc("/find_devices/{customer_id}", FindDevicesHandler)
  router.HandleFunc("/loggedIn", LoggedInHandler)
  router.HandleFunc("/register", RegisterHandler)
  router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

  //Tell the router to server the assets folder as static files
  router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

  //Setup the middleware
  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(router)

  //Run the server
  log.Fatal(http.ListenAndServe(":"+port, n))
}
