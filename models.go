package main

import (
    "time"
)

type Restaurant struct {
    ID int
    Name string `gorm:"size:99; not null"`
    DateCreated time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
}

type ActiveParty struct {
    ID int
    RestaurantID int `gorm:"not null"`
    PartyName string `gorm:"size:50; not null"`
    PartySize int `gorm:"not null"`
    TimeCreated time.Time `gorm:"not null"`
    TimeSeated time.Time
    PhoneAhead bool `gorm:"not null"`
    WaitTimeExpected int
    WaitTimeCalculated int
}

type Buzzer struct {
    ID int
    RestaurantID int `gorm:"not null"`
    BuzzerName string `gorm:"size:45; not null"`
    LastHeartbeat time.Time
    IsActive bool `gorm:"not null"`
    ActiveParty ActiveParty
    ActivePartyID int `gorm:"not null"`
}

type HistoricalParty struct {
    ID int
    RestaurantID int `gorm:"not null"`
    PartyName string `gorm:"size:50;not null"`
    PartySize int `gorm:"not null"`
    DateCreated time.Time `gorm:"not null"`
    DateSeated time.Time `gorm:"not null"`
    WaitTimeExpected int `gorm:"not null"`
    WaitTimeCalculated int `gorm:"not null"`
}

type User struct {
    ID int
    RestaurantID int `gorm:"not null"`
    Username string `gorm:"size:100;not null"`
    Password string `gorm:"size:100; not null"`
    PassSalt string `gorm:"size:50; not null"`
    DateCreated time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
}
