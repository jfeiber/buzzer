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
    "fmt"
    "math"
    _ "github.com/jinzhu/gorm/dialects/postgres"
  )

// RootHandler Handles roots.
func RootHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[RootHandler] ")
  http.Redirect(w, r, "/login", 302)
}

// MakeRandAlphaNumericStr is a back-end method that generates random string for password salt hash.
func MakeRandAlphaNumericStr(n int) string {
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
  b := make([]rune, n)
  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}

// Handle500Error handles 500 error messages.
func Handle500Error(w http.ResponseWriter, err error) {
  w.WriteHeader(500)
  http.Error(w, http.StatusText(500), 500)
  log.Println(err)
}

// GetSession is a back-end method gets and returns information about current session.
func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
  session, err := sessionStore.Get(r, "buzzer-session")
  if err != nil {
    Handle500Error(w, err)
  }
  return session
}

// RenderTemplate is a back-end method to render html templates to user.
func RenderTemplate(w http.ResponseWriter, template_name string, template_params map[string] interface{}) {
  t, err := template.ParseFiles(template_name, "assets/templates/navbar.html.tmpl",
                                "assets/templates/header.html.tmpl")
  template_params["template_name"] = template_name
  if err != nil{
    Handle500Error(w, err)
  }
  t.Execute(w, template_params)
}

// RenderJSONFromMap is a back-end method to create JSON object from passed object map.
func RenderJSONFromMap(w http.ResponseWriter, obj_map map[string] interface{}) {
  json_obj, err := json.Marshal(obj_map)
  if err != nil {
    Handle500Error(w, err)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(json_obj)
}

// ParseReqBody is a back-end method to parse recieved JSON into reqBodyObj object.
func ParseReqBody(r *http.Request, responseObj map[string] interface{},
                  reqBodyObj map[string] interface{}) bool {
  responseObj["status"] = "success"
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    AddErrorMessageToResponseObj(responseObj, "Failed to parse request body.")
    return false
  }
  err = json.Unmarshal(body, &reqBodyObj)
  if err != nil {
    AddErrorMessageToResponseObj(responseObj, "Failed to parse JSON.")
    return false
  }
  return true
}

// ParseReqBodyBuzzer is a back-end method that performs the same functionality as the one above
// but uses the more succinct response language for Buzzer API methods.
func ParseReqBodyBuzzer(r *http.Request, responseObj map[string] interface{},
                  reqBodyObj map[string] interface{}) bool {
  responseObj["e"] = 0
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    AddErrorMessageToResponseObjBuzzer(responseObj, "Failed to parse request body.")
    return false
  }
  err = json.Unmarshal(body, &reqBodyObj)
  if err != nil {
    AddErrorMessageToResponseObjBuzzer(responseObj, "Failed to parse JSON.")
    return false
  }
  return true
}

// AddFlashToSession is a back-end method to add flash to session.
//  "This method is both great and helpful" - jfeiber
func AddFlashToSession(w http.ResponseWriter, r *http.Request, flash string, session *sessions.Session) {
  session.AddFlash(flash)
  session.Save(r, w)
}

// AddErrorMessageToResponseObj is a back-end method to add error message information to responseObj.
// "I think any method I wrote should be extolled in virtues of how great and wonderful they are" - jfeiber
func AddErrorMessageToResponseObj(responseObj map[string] interface{}, errMessage string) {
  responseObj["status"] = "failure"
  responseObj["error_message"] = errMessage
}

// AddErrorMessageToResponseObjBuzzer performs the same functionality as the above method but
// use the more succinct API response used for API endpoints that interact with the Buzzer.
func AddErrorMessageToResponseObjBuzzer(responseObj map[string] interface{}, errMessage string) {
  responseObj["e"] = 1
  responseObj["e_msg"] = errMessage
}

// LoginHandler checks credentials against database and establish session if valid.
// Redirects to/renders Wailtlist page if valid user, display error if not.
// POST contains 'username' and 'password' which are attemped username and password.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[LoginURLHandler] ")
  session := GetSession(w, r)
  if IsUserLoggedIn(session) {
    http.Redirect(w, r, "/waitlist", 302)
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
            http.Redirect(w, r, "/waitlist", 302)
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

// LogoutHandler logs out the current user.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    log.SetPrefix("[LogoutHandler] ")
    session := GetSession(w, r)
    if IsUserLoggedIn(session) && r.Method =="GET" {
        session.Values["username"] = ""
    } else {
        AddFlashToSession(w, r, "Already Logged Out", session)
    }
    session.Save(r, w)
    http.Redirect(w, r, "/login", 302)
    return
}

// IsUserLoggedIn is a back-end method to verify user is logged in and has valid session.
func IsUserLoggedIn(session *sessions.Session) bool {
  username, found := session.Values["username"]
  return found && username != ""
}

// WaitListHandler renders the Wailtlist page after url call recieved.
func WaitListHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[WaitListHandler] ")
  session := GetSession(w, r)
  //confirms valid session
  if !IsUserLoggedIn(session) {
    http.Redirect(w, r, "/login", 302)
    return
  }
  //get current session values
  username, _ := session.Values["username"]
  restaurantID := GetRestaurantIDFromUsername(username.(string))

  var parties []ActiveParty
  //query database for all parties associated with this restaurantID, order by time created asc
  db.Order("time_created asc").Find(&parties, "restaurant_id = ?", restaurantID)

  partyData := map[string]interface{}{}
  partyData["waitlist_data"] = parties
  //This function is called by the template to format the time an ActiveParty was created it in
  //HH:MM form.
  partyData["formatElapsedWaitingTime"] = func (partyCreatedTime time.Time) string {
    duration := time.Now().Sub(partyCreatedTime)
    hours := math.Floor(duration.Hours())
    minutes := math.Floor((duration.Hours()-hours)*60)
    return fmt.Sprintf("%02d:%02d", int(hours), int(minutes))
  }

  //This function is called by the template to format the estimated waiting time into
  //HH:MM format.
  //Kevin had fun making this with his friends.
  partyData["formatEstimatedWaitingTime"] = func (duration int) string {
    var hours = (duration/60)
    var mins = (duration-(hours*60))
    return fmt.Sprintf("%02d:%02d", int(hours), int(mins))
  }

  //render the html template, passing along the data
  RenderTemplate(w, "assets/templates/waitlist.html.tmpl", partyData)
}

// GetActivePartiesHandler is a frontend API Call to update Table of Active Parties.
func GetActivePartiesHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[GetActivePartiesHandler] ")
  session := GetSession(w, r)
  //confirms user session is valid
  if !IsUserLoggedIn(session) {
    http.Redirect(w, r, "/login", 302)
    return
  }
  //retrieve username then restaurantID of current user
  username, _ := session.Values["username"]
  restaurantID := GetRestaurantIDFromUsername(username.(string))

  var parties []ActiveParty
  //query database for all parties with currect restaurantID, order by time partied created asc
  db.Order("time_created asc").Find(&parties, "restaurant_id = ?", restaurantID)

  //create struct and store resuting data from query
  partyData := map[string]interface{}{}
  partyData["waitlist_data"] = parties

  //send to format as JSON and return to frontend
  RenderJSONFromMap(w, partyData);
}

// BuzzerManagementHandler renders the device/buzzer management page.
func BuzzerManagementHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[BuzzerManagerHandler] ")
  session := GetSession(w, r)
  //verify session
  if !IsUserLoggedIn(session) {
    http.Redirect(w, r, "/login", 302)
    return
  }
  //process to assign restaurant id to buzzer
  if r.Method == "POST" {
    //get the name of desired buzzer from POST
    buzzerName := r.FormValue("buzzer_name")
    //pull user data from database based on username from session
    var currUser User
    db.First(&currUser, "username = ?", session.Values["username"])
    //verify that user is in database
    if currUser == (User{}) {
      Handle500Error(w, errors.New("Big problem: The user that is currently logged in does not have an entry in the users table."))
    } else {
      var buzzer Buzzer
      //pull buzzer info from database based on buzzerName from POST
      db.First(&buzzer, "buzzer_name = ? and restaurant_id is null", buzzerName)
      //error if buzzerName does not exist in database
      if buzzer == (Buzzer{}) {
        AddFlashToSession(w, r, "No buzzer with that name found.", session)
      } else {
        //update the found buzzer entry with the current users restaurant id
        db.Model(&buzzer).Update("restaurant_id", currUser.RestaurantID)
      }
    }
  }
  //get username and restaurantID from session
  username, _ := session.Values["username"]
  restaurantID := GetRestaurantIDFromUsername(username.(string))

  var devices []Buzzer
  //query database for all buzzers with the current restaurantID, order by buzzerName asc
  db.Order("buzzer_name asc").Find(&devices, "restaurant_id = ?", restaurantID)
  buzzerData := map[string]interface{}{}
  buzzerData["buzzer_data"] = devices
  if flashes := session.Flashes(); len(flashes) > 0 {
    buzzerData["failure_message"] = flashes[0]
  }
  session.Save(r, w)

  //render buzzer management page and pass along buzzer data
  RenderTemplate(w, "assets/templates/buzzer_management.html.tmpl", buzzerData)
}

// GetLinkedBuzzersHandler is a frontend API class to return updated JSON of buzzers/devices for a specific restaurant.
func GetLinkedBuzzersHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[GetLinkedBuzzerHandler] ")
  session := GetSession(w, r)

  if !IsUserLoggedIn(session) {
    http.Redirect(w, r, "/login", 302)
    return
  }

  username, _ := session.Values["username"]
  restaurantID := GetRestaurantIDFromUsername(username.(string))

  var devices []Buzzer
  db.Order("buzzer_name asc").Find(&devices, "restaurant_id = ?", restaurantID)
  buzzerData := map[string]interface{}{}
  buzzerData["buzzer_data"] = devices

  RenderJSONFromMap(w, buzzerData)
}

// UnlinkBuzzerHandler is a frontend API call to unlink a buzzer from assigned restaurant.
// POST 'buzzer_id' has buzzerID to be unlinked, restaurantID set to null.
func UnlinkBuzzerHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[UnlinkBuzzerHandler] ")
  session := GetSession(w, r)
  if r.Method == "POST" {
    responseObj := map[string] interface{} {}
    reqBodyObj := map[string] interface{}{}
    if !IsUserLoggedIn(session) {
      HandleAuthErrorJson(w, responseObj)
    } else {
      if ParseReqBody(r, responseObj, reqBodyObj) {
        buzzerID := reqBodyObj["buzzer_id"]
        if buzzerID == nil {
          AddErrorMessageToResponseObj(responseObj, "No buzzerID provided.")
        } else {
            var foundBuzzer Buzzer
            db.First(&foundBuzzer, "id = ?", buzzerID)
          if foundBuzzer == (Buzzer{}) {
            AddErrorMessageToResponseObj(responseObj, "Buzzer with that ID not found.")
          } else {
              db.Model(&foundBuzzer).Update("restaurant_id", nil)
          }
        }
      }
    }
    RenderJSONFromMap(w, responseObj)
  }
}

// ActivateBuzzerHandler is a frontend API call to activate a buzzer and alert connected party.
// If information is valid, will set 'is_table_ready' value in database to true.
// POST contains 'active_party_id' which is the related party to be alerted.
func ActivateBuzzerHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[ActivateBuzzer] ")
  session := GetSession(w, r)
  if r.Method == "POST" {
    responseObj := map[string] interface{} {}
    reqBodyObj := map[string] interface{}{}
    if !IsUserLoggedIn(session) {
      HandleAuthErrorJson(w, responseObj)
    } else {
      if ParseReqBody(r, responseObj, reqBodyObj) {
        activePartyID := reqBodyObj["active_party_id"]
        if activePartyID == nil {
          AddErrorMessageToResponseObj(responseObj, "No activePartyID provided.")
        } else {
            var foundActiveParty ActiveParty
            db.First(&foundActiveParty, "id = ?", activePartyID)
          if foundActiveParty == (ActiveParty{}) {
            AddErrorMessageToResponseObj(responseObj, "Party with that ID not found.")
          } else {
              db.Model(&foundActiveParty).Update("is_table_ready", true)
          }
        }
      }
    }

    RenderJSONFromMap(w, responseObj)
  }
}

// UpdatePhoneAheadStatusHandler is a frontend API call to update party status from PhoneAhead to Waitlist.
//  POST contains 'active_party_id' which is the party whose status is to be updated.
func UpdatePhoneAheadStatusHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[UpdatePhoneAheadStatusHandler] ")
  session := GetSession(w, r)
  if r.Method == "POST" {
    responseObj := map[string] interface{} {}
    reqBodyObj := map[string] interface{}{}
    if !IsUserLoggedIn(session) {
      HandleAuthErrorJson(w, responseObj)
    } else {
      if ParseReqBody(r, responseObj, reqBodyObj) {
        activePartyID := reqBodyObj["active_party_id"]
        if activePartyID == nil {
          AddErrorMessageToResponseObj(responseObj, "No activePartyID provided.")
        } else {
            var foundActiveParty ActiveParty
            db.First(&foundActiveParty, "id = ?", activePartyID)
          if foundActiveParty == (ActiveParty{}) {
            AddErrorMessageToResponseObj(responseObj, "Party with that ID not found.")
          } else {
                db.Model(&foundActiveParty).Update("phone_ahead", false)
            }
          }
        }
      }

    RenderJSONFromMap(w, responseObj)
  }
}

// GetBuzzerObjFromName is a back-end method to return all information (as object) on a buzzer based on buzzerName.
// Passed reqBodyObj contains 'buzzer_name' which is the buzzerName to query by.
func GetBuzzerObjFromName(reqBodyObj map[string] interface{}, responseObj map[string] interface {}, buzzer *Buzzer) bool {
  buzzerName := reqBodyObj["bn"]
  if buzzerName == nil {
    AddErrorMessageToResponseObjBuzzer(responseObj, "buzzer_name field required.")
    return false
  } else {
    db.First(buzzer, "buzzer_name = ?", buzzerName)
    if *buzzer == (Buzzer{}) {
      AddErrorMessageToResponseObjBuzzer(responseObj, "Buzzer with that name not found.")
      return false
    }
  }
  return true
}


// GetBuzzerObjFromID is a back-end method to return all information (as object) on a buzzer based on buzzerID.
// Passed buzzerID is the buzzerID to query by.
func GetBuzzerObjFromID(buzzerID int, responseObj map[string] interface{}, buzzer *Buzzer) bool {
  db.First(buzzer, "id = ?", buzzerID)
  if *buzzer == (Buzzer{}) {
    AddErrorMessageToResponseObj(responseObj, "Buzzer with that ID not found.")
    return false
  }
  return true
}

// GetActivePartyFromBuzzerID is a back-end method to determine if buzzer is connected to an active
// party and return active party info. buzzerID retrieved from passed buzzer object.
// Returns false if no asscoicated active party, else returns true and sets passed activeParty
// pointer to found party.
func GetActivePartyFromBuzzerID(responseObj map[string] interface{}, buzzer Buzzer, activeParty *ActiveParty) bool {
  db.First(activeParty, "buzzer_id = ?", buzzer.ID)
  if *activeParty == (ActiveParty{}) {
    AddErrorMessageToResponseObj(responseObj, "No active party with that buzzer id found and the buzzer is active.")
    return false
  }
  return true
}

// GetActivePartyFromID is a back-end method to determine if active party exists.
// Passed reqBodyObj contains 'party_id' which is the activePartyID to check.
// Returns false if party does not exist, else returns true and sets passed activeParty
// pointer to found party.
func GetActivePartyFromID(reqBodyObj map[string] interface{}, responseObj map[string] interface{}, activeParty *ActiveParty) bool {
  db.First(activeParty, "id = ?", reqBodyObj["party_id"])
  if *activeParty == (ActiveParty{}) {
    AddErrorMessageToResponseObjBuzzer(responseObj, "Party with that ID not found.")
    return false
  }
  return true
}

// GetAvailablePartyHandler is a buzzer API method to get an active party to potentially be assigned to buzzer.
// Returns first result from database of party with no assigned buzzer and not phone ahead.
// The 'p_a' field in the response indicates whether or not a party is available. If a party is
// available, then the response will also contain the name of the party ('n'), the estimated wait
// time for the party ('t'), and the ID of the party ('id'). The response will also contain the
// usual error info.
func GetAvailablePartyHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[GetAvailablePartyHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  if ParseReqBodyBuzzer(r, responseObj, reqBodyObj) {
    var buzzer Buzzer
    if GetBuzzerObjFromName(reqBodyObj, responseObj, &buzzer) {
      var activeParty ActiveParty
      db.First(&activeParty, "restaurant_id = ? and buzzer_id is null and phone_ahead is false", buzzer.RestaurantID)
      responseObj["p_a"] = 1;
      if activeParty != (ActiveParty{}) {
        responseObj["n"] = activeParty.PartyName
        responseObj["t"] = activeParty.WaitTimeExpected
        responseObj["id"] = activeParty.ID
      } else {
        responseObj["p_a"] = 0;
      }
    }
  }
  RenderJSONFromMap(w, responseObj)
}

// AcceptPartyHandler is a buzzer API method handles response from buzzer accepting assignment to a active party.
// If the accepted information is valid, the database is updated to reflect buzzer assignment.
// reqBodyObj must contain 'bn' (the assigned name of requesting buzzer) and 'id' (id of active
// party that the buzzer is accepting).
func AcceptPartyHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[AcceptPartyHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  if ParseReqBodyBuzzer(r, responseObj, reqBodyObj) {
    if reqBodyObj["bn"] == nil || reqBodyObj["id"] == nil {
      AddErrorMessageToResponseObjBuzzer(responseObj, "Missing required fields.")
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
          AddErrorMessageToResponseObjBuzzer(responseObj, "Can't accept a party that already has a buzzer")
        }
      }
    }
  }
  RenderJSONFromMap(w, responseObj)
}

// HeartbeatHandler is used by the buzzers to check in periodically. Right now that period is ~30
// seconds. If a party has been marked inactive or a table is ready then the Buzzer will receive
// that info in the response to this endpoint.
// The 'i_a' field in the response indicates whether or not a party is active. In the future
// this should return that a party is inactive if it's not in the ActiveParties DB as parties
// that are no longer active will be moved to HistoricalParties.
// The "t" field represents the expected wait time.
// When the "b" field is 1 the buzzer will buzz.
func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[HeartbeatHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{} {}
  if ParseReqBodyBuzzer(r, responseObj, reqBodyObj) {
    log.Println(reqBodyObj)
    var buzzer Buzzer
    if GetBuzzerObjFromName(reqBodyObj, responseObj, &buzzer) {
      responseObj["i_a"] = 0
      if buzzer.IsActive {
        responseObj["i_a"] = 1
      }
      if buzzer.IsActive {
        db.Model(&buzzer).Update("last_heartbeat", time.Now().UTC())
        var activeParty ActiveParty
        if GetActivePartyFromBuzzerID(responseObj, buzzer, &activeParty) {
          responseObj["t"] = activeParty.WaitTimeExpected
          responseObj["b"] = 0
          if activeParty.IsTableReady {
            responseObj["b"] = 1
          }
        }
      }
    }
  }
  RenderJSONFromMap(w, responseObj)
}

// IsBuzzerRegisteredHandler is a buzzer API method checks to see if specified buzzer is
// assigned/registered with a resturant. Uses buzzer name in reqBodyObj to get the related buzzer
// object adn check for RestaurantID. The response field 'i_reg' is 1 if the buzzer is registered,
// 0 otherwise.
func IsBuzzerRegisteredHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[IsBuzzerRegisteredHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  if ParseReqBodyBuzzer(r, responseObj, reqBodyObj) {
    var buzzer Buzzer
    if GetBuzzerObjFromName(reqBodyObj, responseObj, &buzzer) {
      responseObj["i_reg"] = 0
      if  buzzer.RestaurantID != 0 {
        responseObj["i_reg"] = 1
      }
    }
  }
  RenderJSONFromMap(w, responseObj)
}

// DeleteActivePartyHandler is a frontend API method deletes the specificed active party through
// given ID. Retrieves specified activePartyID from reqBodyObj 'active_party_id' and deletes from
// activeParty table. Removed party and all related information is then stored in historicalParty
// table by called fucntion.
//TODO: Move the active parties into historical parties.
func DeleteActivePartyHandler(w http.ResponseWriter, r *http.Request) {
    log.SetPrefix("[DeleteActivePartyHandler] ")
    responseObj := map[string] interface{} {}
    reqBodyObj := map[string] interface{}{}
    session := GetSession(w, r)
    if !IsUserLoggedIn(session) {
      HandleAuthErrorJson(w, responseObj)
    } else if ParseReqBody(r, responseObj, reqBodyObj) {
        activePartyID := reqBodyObj["active_party_id"]
        if activePartyID == nil {
            responseObj["status"] = "failure"
            responseObj["error_message"] = "Missing active_party_id parameter"
        } else {
            var activeParty ActiveParty
            db.First(&activeParty, "ID=?", activePartyID)
            failedBuzzerUpdate := false
            if (activeParty.BuzzerID != 0) {
              var buzzer Buzzer
              if GetBuzzerObjFromID(activeParty.BuzzerID, responseObj, &buzzer) {
                db.Model(&buzzer).Update("is_active", false)
              } else {
                failedBuzzerUpdate = true
              }
            }
            if !failedBuzzerUpdate {
              dbInfo := db.Delete(&activeParty)
              if dbInfo.Error == nil {
                  responseObj["status"] = "success"
              } else {
                  responseObj["status"] = "failure"
                  responseObj["error_message"] = "db.Delete failed"
              }
            }
        }
    }
    RenderJSONFromMap(w, responseObj)
}

// GetNewBuzzerNameHandler is a buzzer API method to assign name to unnnamed device.
// Uses buzzerNameGenerator to generate new name in proper format.
// Creates entry in database for newly connected buzzer.
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

// HandleAuthErrorJson is a back-end method to handle authorization error message output.
func HandleAuthErrorJson(w http.ResponseWriter, responseObj map[string] interface{}) {
  w.WriteHeader(401)
  responseObj["status"] = "failure"
  responseObj["error_message"] = "User not logged in"
}

// GetRestaurantIDFromUsername is a back-end method to retrieve associated RestaurantID from
// specified username. Username is passed in via string and found RestaurantID is retuned to caller.
func GetRestaurantIDFromUsername(username string) int {
  var currUser User
  db.First(&currUser, "username = ?", username)
  if currUser == (User{}) {
    return -1;
  }
  return currUser.RestaurantID
}

// CreateNewPartyHandler is a frontend API method to add a new party to the database/waitlist.
// reqBodyObj contains 'party_name', 'party_size', 'wait_time_expected', 'phone_ahead'.
// Completion status returned in responseObj along wiht assigned 'active_party_id' if added
// successfully.
func CreateNewPartyHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[CreateNewPartyHandler] ")
  responseObj := map[string] interface{} {}
  reqBodyObj := map[string] interface{}{}
  session := GetSession(w, r)
  if !IsUserLoggedIn(session) {
    HandleAuthErrorJson(w, responseObj)
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

// AddUserHandler is a hanlder to render the add user page and handle new user additions.
// POST contains new user info in 'username', 'password', and 'restaurant_name'.
// Inputed password is run through salted hash before new user data stored in database.
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

// AnalyticsHandler renders the analytics page
func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
    log.SetPrefix("[AnalyticsHandler] ")
    session := GetSession(w, r)
    if !IsUserLoggedIn(session) {
        http.Redirect(w, r, "/login", 302)
        return
    }
    RenderTemplate(w, "assets/templates/analytics.html.tmpl", map[string]interface{}{})
}

// GetHistoricalPartiesHandler Returns a json of historicalparties within a time range
func GetHistoricalPartiesHandler(w http.ResponseWriter, r *http.Request) {
    log.SetPrefix("[GetHistoricalPartiesHandler]")
    returnObj := map[string] interface{} {"status": "success"}
    session := GetSession(w, r)

    if !IsUserLoggedIn(session) {
      HandleAuthErrorJson(w, returnObj)
    } else if r.Method == "POST" {
        startEndInfo := map[string] interface{}{}
        if ParseReqBody(r, returnObj, startEndInfo) {
            username, _ := session.Values["username"]
            restaurantID := GetRestaurantIDFromUsername(username.(string))

            var format = "01/02/2006"
            var startDate, endDate time.Time
            var err interface{}
            if val, ok := startEndInfo["start_date"].(string); ok {
                startDate, err = time.Parse(format, val)
                if err != nil {
                    returnObj["status"] = "failure"
                    returnObj["error_message"] = "time.Parse failed"
                }
            } else {
                returnObj["status"] = "failure"
                returnObj["error_message"] = "start date undefined"
            }

            if val, ok := startEndInfo["end_date"].(string); ok {
                endDate, err = time.Parse(format, val)
                if err != nil {
                    returnObj["status"] = "failure"
                    returnObj["error_message"] = "time.Parse failed"
                }
            } else {
                returnObj["status"] = "failure"
                returnObj["error_message"] = "end date undefined"
            }

            startDateFormatted := startDate.Format("2006-01-02 15:04:05")
            endDateFormatted := endDate.Format("2006-01-02 15:04:05")
            var historicalParties []HistoricalParty
            db.Where("restaurant_id = ? AND time_created >= ? AND time_seated <= ?", restaurantID, startDateFormatted, endDateFormatted).Find(&historicalParties)
            if len(historicalParties) > 0 {
                returnObj["historical_parties"] = historicalParties
            }
        }
    }
    RenderJSONFromMap(w, returnObj)
}

// IsPartyAssignedBuzzerHandler is a frontend API method to check if specified active party is
// assigned buzzer. Passed object r contains 'active_party_id' to be quieried for, returnObj
// contains response 'is_party_assigned_buzzer'. Used by fronted to check if buzzer has been
// successfully assigend to party after party was created.
func IsPartyAssignedBuzzerHandler(w http.ResponseWriter, r *http.Request) {
  returnObj := map[string] interface{} {"status": "success"}
  if !IsUserLoggedIn(GetSession(w, r)) {
    HandleAuthErrorJson(w, returnObj)
  } else if r.Method == "POST" {
    activePartyInfo := map[string] interface{}{}
    if ParseReqBody(r, returnObj, activePartyInfo) {
      var activeParty ActiveParty
      activePartyID := activePartyInfo["active_party_id"]; if activePartyID == nil {
        returnObj["status"] = "failure"
        returnObj["error_message"] = "Missing active_party_id parameter"
      } else {
        db.First(&activeParty, "id = ?", activePartyID)
        if activeParty == (ActiveParty{}) {
          returnObj["status"] = "failure"
          returnObj["error_message"] = "Party with the provided ID not found"
        }
        log.Println(activeParty);
        if (activeParty.BuzzerID == 0) {
          returnObj["is_party_assigned_buzzer"] = false
        } else {
          returnObj["is_party_assigned_buzzer"] = true
        }
        returnObj["active_party_id"] = activePartyID
      }
    }
  }
  RenderJSONFromMap(w, returnObj)
}

// NotFoundHandler is a handler to render 404 Not Found page.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
  log.SetPrefix("[NotFoundHandler] ")
  w.WriteHeader(404)
  RenderTemplate(w, "assets/templates/404.html", map[string]interface{}{})
}
