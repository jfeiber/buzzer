package main

import (
    "time"
)

type Restaurant struct {
    id int
    name string `gorm:"size:99; not null"`
    dateCreated time.Time `gorm:"not null"`
}

type ActiveParty struct {
    id int
    restaurant Restaurant
    restaurantID int `gorm:"not null"`
    partyName string `gorm:"size:50; not null"`
    partySize int `gorm:"not null"`
    timeCreated time.Time `gorm:"not null"`
    timeSeated time.Time
    phoneAhead bool `gorm:"not null"`
    waitTimeExpected time.Time
    waitTimeCalculated time.Time
}

type Buzzer struct {
    id int
    restaurant Restaurant
    buzzerName string `gorm:"size:45; not null"`
    lastHeartbeat time.Time
    isActive bool `gorm:"not null"`
    activeParty ActiveParty
    activePartyID int `gorm:"not null"`
}

type HistoricalParty struct {
    id int
    restaurant Restaurant
    restaurantId int `gorm:"not null"`
    partyName string `gorm:"size:50;not null"`
    partySize int `gorm:"not null"`
    dateCreated time.Time `gorm:"not null"`
    dateSeated time.Time `gorm:"not null"`
    waitTimeExpected int `gorm:"not null"`
    waitTimeCalculated int `gorm:"not null"`
}

type User struct {
    id int
    restaurant Restaurant
    restaurantId int `gorm:"not null"`
    username string `gorm:"size:100;not null"`
    password string `gorm:"size:100; not null"`
    passSalt string `gorm:"size:50; not null"`
    dateCreated time.Time `gorm:"not null"`
}
