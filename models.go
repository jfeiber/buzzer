package main

import (
    "time"
)

type Device struct {
    ID int
    CustomerID int `gorm:"not null"`
    DeviceName string `gorm:"not null unique"`
    LastHeartbeat time.Time
    IsActive bool `gorm:"not null"`
    PartyName string `gorm:"size:50"`
    PartySize int
    WaitTime int
}

type Restaurant struct {
    restaurantID int
    name string `gorm:"size:99; not null"`
    dateCreated time.Time `gorm:"not null"`
}

type ActiveParty struct {
    activePartyID int
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
    buzzerID int
    restaurant Restaurant
    buzzerName string `gorm:"size:45; not null"`
    lastHeartbeat time.Time
    isActive bool `gorm:"not null"`
    activeParty ActiveParty
    activePartyID int `gorm:"not null"`
}

type HistoricalParties struct {
    historicalPartiesID int
    restaurant Restaurant
    restaurantId int `gorm:"not null"`
    partyName string `gorm:"size:50;not null"`
    partySize int `gorm:"not null"`
    dateCreated time.Time `gorm:"not null"`
    dateSeated time.Time `gorm:"not null"`
    waitTimeExpected int `gorm:"not null"`
    waitTimeCalculated int `gorm:"not null"`
}

type WebAppUser struct {
    webAppUserID int
    restaurant Restaurant
    restaurantId int `gorm:"not null"`
    username string `gorm:"size:100;not null"`
    password string `gorm:"size:100; not null"`
    passSalt string `gorm:"size:50; not null"`
    dateCreated time.Time `gorm:"not null"`
}
