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
    restaurantID uint
    name string `gorm:"size:99"`
    dateCreated time.Time
}

type ActiveParty struct {
    activePartyID uint
    restaurant Restaurant
    restaurantID uint
    partyName string `gorm:"size:100"`
    partySize string `gorm:"size:100"`
    timeCreated time.Time
    timeSeated time.Time
    phoneAhead bool
    waitTimeExpected time.Time
    waitTimeCalculated time.Time
}

type Buzzer struct {
    buzzerID uint
    restaurant Restaurant
    name string `gorm:"size:45"`
    last_heartbeat time.Time
    is_active bool
    activeParty ActiveParty
}

type HistoricalParties struct {
    historicalPartiesID int
    restaurant Restaurant
    restaurantId int
    partyName string `gorm:"size:50;not null"`
    partySize int
    dateCreated time.Time
    dateSeated time.Time
    waitTimeExpected int
    waitTimeCalc int
}

type WebAppUser struct {
    WebAppUserID int
    restaurant Restaurant
    restaurantId int
    username string `gorm:"size:100;not null"`
    password string `gorm:"size:100; not null"`
    passSalt string `gorm:"size:50; not null"`
    dateCreated time.Time
}
