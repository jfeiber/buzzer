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
  router.HandleFunc("/logout", LogoutHandler)
  router.HandleFunc("/add_user", AddUserHandler)
  router.HandleFunc("/analytics", AnalyticsHandler)
  router.HandleFunc("/waitlist", WaitListHandler)
  router.HandleFunc("/admin", UserAdminHandler)
  router.HandleFunc("/buzzer_management", BuzzerManagementHandler)
  router.HandleFunc("/buzzer_api/get_new_buzzer_name", GetNewBuzzerNameHandler)
  router.HandleFunc("/buzzer_api/is_buzzer_registered", IsBuzzerRegisteredHandler)
  router.HandleFunc("/buzzer_api/get_available_party", GetAvailablePartyHandler)
  router.HandleFunc("/buzzer_api/accept_party", AcceptPartyHandler)
  router.HandleFunc("/buzzer_api/heartbeat", HeartbeatHandler)
  router.HandleFunc("/frontend_api/create_new_party", CreateNewPartyHandler)
  router.HandleFunc("/frontend_api/delete_party", DeleteActivePartyHandler)
  router.HandleFunc("/frontend_api/remove_user", RemoveUserHandler)
  router.HandleFunc("/frontend_api/get_active_parties", GetActivePartiesHandler)
  router.HandleFunc("/frontend_api/get_users", GetUsersHandler)
  router.HandleFunc("/frontend_api/is_party_assigned_buzzer", IsPartyAssignedBuzzerHandler)
  router.HandleFunc("/frontend_api/activate_buzzer", ActivateBuzzerHandler)
  router.HandleFunc("/frontend_api/update_phone_ahead_status", UpdatePhoneAheadStatusHandler)
  router.HandleFunc("/frontend_api/update_party_size", UpdatePartySizeHandler)
  router.HandleFunc("/frontend_api/unlink_buzzer", UnlinkBuzzerHandler)
  router.HandleFunc("/frontend_api/get_linked_buzzers", GetLinkedBuzzersHandler)
  router.HandleFunc("/analytics_api/get_average_party_chart", GetAveragePartySizeChartHandler)
  router.HandleFunc("/analytics_api/get_total_customers_chart", GetTotalCustomersChartHandler)
  router.HandleFunc("/analytics_api/get_parties_hour_chart", GetParitesPerHourChartHandler)
  router.HandleFunc("/analytics_api/get_party_loss_chart", GetPartyLossChartHandler)
  router.HandleFunc("/analytics_api/get_avg_wait_chart", GetAvgWaittimeChartHandler)


  router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

  //Tell the router to server the assets folder as static files
  router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

  //Setup the middleware
  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(router)

  //Run the server
  log.Fatal(http.ListenAndServe(":"+port, n))
}
