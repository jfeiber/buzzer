package main

import (
    "log"
    "net/http"
    "errors"
    "html/template"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/sessions"
    "encoding/json"
    "math/rand"
    "time"
    "io/ioutil"
    _ "github.com/jinzhu/gorm/dialects/postgres"
  )

func MakeRandAlphaNumericStr(n int) string {
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
  b := make([]rune, n)
  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}

func Handle500Error(w http.ResponseWriter, err error) {
  http.Error(w, http.StatusText(500), 500)
  log.Println(err)
}

func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
  session, err := sessionStore.Get(r, "buzzer-session")
  if err != nil {
    Handle500Error(w, err)
  }
  return session
}

func RenderTemplate(w http.ResponseWriter, template_name string, template_params map[string] interface{}) {
  t, err := template.ParseFiles(template_name)
  if err != nil{
    Handle500Error(w, err)
  }
  t.Execute(w, template_params)
}

func RenderJSONFromMap(w http.ResponseWriter, obj_map map[string] interface{}) {
  json_obj, err := json.Marshal(obj_map)
  if err != nil {
    Handle500Error(w, err)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(json_obj)
}

func AddFlashToSession(w http.ResponseWriter, r *http.Request, flash string, session *sessions.Session) {
  session.AddFlash(flash)
  session.Save(r, w)
}

func IsUserLoggedIn(session *sessions.Session) bool {
  username, found := session.Values["username"]
  return found && username != ""
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[LoginURLHandler] ")
  session := GetSession(w, r)
  if IsUserLoggedIn(session) {
    http.Redirect(w, r, "/wait_list", 302)
    return
  }
  if r.Method == "POST" {
    log.Println("in post")
    username := r.FormValue("username")
    password := r.FormValue("password")

    if username == "" || password == "" {
      AddFlashToSession(w, r, "Missing form field", session)
    } else {
      var user User;
      db.First(&user, "Username = ?", username)
      if (user != (User{})) {
        passSalt := user.PassSalt
        if (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+passSalt)) == nil) {
            session.Values["username"] = username
            session.Save(r, w)
            http.Redirect(w, r, "/wait_list", 302)
            return
        }
      }
      AddFlashToSession(w, r, "Username or password is incorrect", session)
    }
  }
  template_data := map[string]interface{}{}
  if flashes := session.Flashes(); len(flashes) > 0 {
      template_data["failure_message"] = flashes[0]
  }
  session.Save(r, w)
  RenderTemplate(w, "assets/templates/login.html.tmpl", template_data)
}

func WaitListHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[WaitListHandler] ")
  session := GetSession(w, r)
  if !IsUserLoggedIn(session) {
    http.Redirect(w, r, "/login", 302)
    return
  }

  username, _ := session.Values["username"]
  restaurantID := GetRestaurantIDFromUsername(username.(string))

  var parties []ActiveParty
  db.Find(&parties, "restaurant_id = ?", restaurantID)

  party_data := map[string]interface{}{}
  party_data["waitlist_data"] = parties

  RenderTemplate(w, "assets/templates/waitlist.html.tmpl", party_data)
}

func UpdateWaitlist(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[UpdateWaitlist] ")
  session := GetSession(w, r)
//  if !IsUserLoggedIn(session) {}
  username, _ := session.Values["username"]
  restaurantID := GetRestaurantIDFromUsername(username.(string))

  var parties []ActiveParty
  db.Find(&parties, "restaurant_id = ?", restaurantID)

  party_data := map[string]interface{}{}
  party_data["waitlist_data"] = parties

  RenderJSONFromMap(w, party_data);
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RootHandler] ")
  http.Redirect(w, r, "/login", 302)
}

func RegisterBuzzerHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RegisterBuzzerHandler] ")
  session := GetSession(w, r)
  if !IsUserLoggedIn(GetSession(w, r)) {
    http.Redirect(w, r, "/login", 302)
  }
  if r.Method == "POST" {
    buzzerName := r.FormValue("buzzer_name")

    var currUser User
    db.First(&currUser, "username = ?", session.Values["username"])
    if currUser == (User{}) {
      Handle500Error(w, errors.New("Big problem: The user that is currently logged in does not have an entry in the users table."))
    } else {
      var buzzer Buzzer
      db.First(&buzzer, "buzzer_name = ?", buzzerName)
      if buzzer == (Buzzer{}) {
        AddFlashToSession(w, r, "No buzzer with that name found.", session)
      } else {
        db.Model(&buzzer).Update("restaurant_id", currUser.RestaurantID)
      }
    }
  }
  templateData := map[string]interface{}{}
  if flashes := session.Flashes(); len(flashes) > 0 {
    templateData["failure_message"] = flashes[0]
  }
  session.Save(r, w)
  RenderTemplate(w, "assets/templates/register_buzzer.html.tmpl", templateData)
}

func AddErrorMessageToResponseObj(responseObj map[string] interface{}, err_message string) {
  responseObj["status"] = "failure"
  responseObj["error_message"] = err_message
}

func ParseReqBody(r *http.Request, responseObj map[string] interface{},
                  reqBodyObj map[string] interface{}) bool {
  responseObj["status"] = "success"
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    responseObj["status"] = "failure"
    responseObj["error_message"] = "Failed to parse request body."
    return false
  }
  err = json.Unmarshal(body, &reqBodyObj)
  if err != nil {
    responseObj["status"] = "failure"
    responseObj["error_message"] = "Failed to parse JSON."
    return false
  }
  return true
}

func GetBuzzerObjFromName(reqBodyObj map[string] interface{}, responseObj map[string] interface {}, buzzer *Buzzer) bool {
  buzzerName := reqBodyObj["buzzer_name"]
  if buzzerName == nil {
    AddErrorMessageToResponseObj(responseObj, "buzzer_name field required.")
    return false
  } else {
    db.First(buzzer, "buzzer_name = ?", buzzerName)
    if *buzzer == (Buzzer{}) {
      AddErrorMessageToResponseObj(responseObj, "Buzzer with that name not found.")
      return false
    }
  }
  return true
}

func GetActivePartyFromBuzzerID(responseObj map[string] interface{}, buzzer Buzzer, activeParty *ActiveParty) bool {
  db.First(activeParty, "buzzer_id = ?", buzzer.ID)
  if *activeParty == (ActiveParty{}) {
    AddErrorMessageToResponseObj(responseObj, "No active party with that buzzer id found and the buzzer is active.")
    return false
  }
  return true
}

func GetActivePartyFromID(reqBodyObj map[string] interface{}, responseObj map[string] interface{}, activeParty *ActiveParty) bool {
  db.First(activeParty, "id = ?", reqBodyObj["party_id"])
  if *activeParty == (ActiveParty{}) {
    AddErrorMessageToResponseObj(responseObj, "Party with that ID not found.")
    return false
  }
  return true
}

func AcceptPartyHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[AcceptPartyHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  if ParseReqBody(r, responseObj, reqBodyObj) {
    if reqBodyObj["buzzer_name"] == nil || reqBodyObj["party_id"] == nil {
      AddErrorMessageToResponseObj(responseObj, "Missing required fields.")
    } else {
      var activeParty ActiveParty
      if GetActivePartyFromID(reqBodyObj, responseObj, &activeParty) {
        if activeParty.BuzzerID == 0 {
          var buzzer Buzzer
          if GetBuzzerObjFromName(reqBodyObj, responseObj, &buzzer) {
            db.Model(&activeParty).Update("buzzer_id", buzzer.ID)
            db.Model(&buzzer).Update("is_active", true)
          }
        } else {
          AddErrorMessageToResponseObj(responseObj, "Can't accept a party that already has a buzzer")
        }
      }
    }
  }
  RenderJSONFromMap(w, responseObj)
}

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[HeartbeatHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  if ParseReqBody(r, responseObj, reqBodyObj) {
    var buzzer Buzzer
    if GetBuzzerObjFromName(reqBodyObj, responseObj, &buzzer) {
      responseObj["is_active"] = buzzer.IsActive
      if buzzer.IsActive {
        db.Model(&buzzer).Update("last_heartbeat", time.Now().UTC())
        var activeParty ActiveParty
        if GetActivePartyFromBuzzerID(responseObj, buzzer, &activeParty) {
          responseObj["wait_time"] = activeParty.WaitTimeExpected
          responseObj["buzz"] = activeParty.IsTableReady
        }
      }
    }
  }
  RenderJSONFromMap(w, responseObj)
}

func GetAvailablePartyHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[GetAvailablePartyHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  if ParseReqBody(r, responseObj, reqBodyObj) {
    var buzzer Buzzer
    if GetBuzzerObjFromName(reqBodyObj, responseObj, &buzzer) {
      var activeParty ActiveParty
      db.First(&activeParty, "restaurant_id = ? and buzzer_id is null and phone_ahead is false", buzzer.RestaurantID)
      responseObj["party_avail"] = true;
      if activeParty != (ActiveParty{}) {
        responseObj["party_name"] = activeParty.PartyName
        responseObj["wait_time"] = activeParty.WaitTimeExpected
        responseObj["party_id"] = activeParty.ID
      } else {
        responseObj["party_avail"] = false;
      }
    }
  }
  RenderJSONFromMap(w, responseObj)
}

func IsBuzzerRegisteredHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[IsBuzzerRegisteredHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  if ParseReqBody(r, responseObj, reqBodyObj) {
    var buzzer Buzzer
    if GetBuzzerObjFromName(reqBodyObj, responseObj, &buzzer) {
      responseObj["is_buzzer_registered"] = buzzer.RestaurantID != 0
    }
  }
  RenderJSONFromMap(w, responseObj)
}

func GetNewBuzzerNameHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[GenerateBuzzerNameHandler] ")
  buzzerName := buzzerNameGenerator.GenerateName()
  var buzzer Buzzer
  db.First(&buzzer, "buzzer_name = ?", buzzerName)
  for buzzer != (Buzzer{}) {
    buzzerName = buzzerNameGenerator.GenerateName()
    db.First(&buzzer, "buzzer_name = ?", buzzerName)
  }
  buzzer = Buzzer{BuzzerName: buzzerName, LastHeartbeat: time.Now().UTC(), IsActive: false}
  // db.NewRecord(buzzer)
  if err := db.Create(&buzzer).Error; err != nil {
    Handle500Error(w, err)
  }
  obj_map := map[string] interface{} {"status": "success", "buzzer_name": buzzerName}
  RenderJSONFromMap(w, obj_map)
}

func HandleAuthErrorJson(responseObj map[string] interface{}) {
  responseObj["status"] = "failure"
  responseObj["error_message"] = "User not logged in."
}

func GetRestaurantIDFromUsername(username string) int {
  var currUser User
  db.First(&currUser, "username = ?", username)
  if currUser == (User{}) {
    return -1;
  }
  return currUser.RestaurantID
}

func CreateNewPartyHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[CreateNewPartyHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  session := GetSession(w, r)
  if !IsUserLoggedIn(session) {
    HandleAuthErrorJson(responseObj)
  } else {
    if ParseReqBody(r, responseObj, reqBodyObj) {
      log.Println(reqBodyObj)
      partyName := reqBodyObj["party_name"]
      partySize := reqBodyObj["party_size"]
      waitTimeExpected := reqBodyObj["wait_time_expected"]
      phoneAhead := reqBodyObj["phone_ahead"]
      if partyName == nil || partySize == nil || waitTimeExpected == nil || phoneAhead == nil {
        responseObj["status"] = "failure"
        responseObj["error_message"] = "Missing parameters."
      } else {
        username, _ := session.Values["username"]
        restaurantID := GetRestaurantIDFromUsername(username.(string))
        if restaurantID == -1 {
          Handle500Error(w, errors.New("Big problem: The user that is currently logged in does not have an entry in the users table."))
        } else {
          activeParty := ActiveParty{RestaurantID: restaurantID, PartyName: partyName.(string), PartySize: int(partySize.(float64)), PhoneAhead: phoneAhead.(bool), WaitTimeExpected: int(waitTimeExpected.(float64))}
          db.Create(&activeParty)
          responseObj["status"] = "success"
          responseObj["active_party_id"] = activeParty.ID
        }
      }
    }
  }
  RenderJSONFromMap(w, responseObj)
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[AddUserHandler] ")
  session := GetSession(w, r)
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

      var restaurant Restaurant
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
        AddFlashToSession(w, r, "Username already exists", session)
      } else {
        //add the user
        user = User{RestaurantID: restaurant.ID, Username: username, Password: string(hashedPassword), PassSalt: passSalt}
        db.NewRecord(user)
        db.Create(&user)
        AddFlashToSession(w, r, "User successfully added", session)
      }
    } else {
      AddFlashToSession(w, r, "Could not add user. Did you forget a field?", session)
    }
    http.Redirect(w, r, "/add_user", 302)
  } else {
    template_data := map[string]interface{}{}
    if flashes := session.Flashes(); len(flashes) > 0 {
      template_data["flash"] = flashes[0]
    }
    session.Save(r, w)
    RenderTemplate(w, "assets/templates/adduser.html.tmpl", template_data)
  }
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[NotFoundHandler] ")
  RenderTemplate(w, "assets/templates/404.html", nil)
}
