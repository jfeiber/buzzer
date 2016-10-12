package main

import "time"

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

type HistoricalParties struct {
    historicalPartiesID int
    restaurant Restaurant
    restaurantId int
    partyName string `gorm:"size:50;not null"`
    partySize int
    dateCreated
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
