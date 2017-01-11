func returnHistoricalPartiesFromRestaurantAndTime(restaurantName string, timeCreatedText string, timeSeatedText string) []HistoricalParty {

    // clean up restaurant name
    var currRestaurant Restaurant
    db.First(&currRestaurant, "name = ?", restaurantName)

    if currRestaurant == (Restaurant{}) {
      returnObj["status"] = "failure"
      returnObj["error_message"] = "Party with the provided ID not found"
    }

    var historicalParties []HistoricalParty
    var format = "asdasdad"

    var timeCreated, err := time.Parse(format, timeCreatedText)
    var timeSeated, err := time.Parse(restaurantName, timeSeatedText)

    var timeCreatedFormatted := timeCreated.Format("2006-01-02 15:04:05")

    db.Where("restaurant_id = ? AND time_created >= ? AND time_seated <= ?", currRestaurant.id, timeCreated, timeSeated).find(&historicalParties)
}

resturnHistoricalPartiesFromRestaurantAndTime()