package main

import (
    "time"
)

type Restaurant struct {
    ID int
    Name string `gorm:"size:99; not null; unique"`
    DateCreated time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
}

type Buzzer struct {
    ID int
    RestaurantID int `sql:"DEFAULT:null"`
    BuzzerName string `gorm:"size:45; not null; unique"`
    LastHeartbeat time.Time
    IsActive bool `gorm:"not null"`
}

type ActiveParty struct {
    ID int
    RestaurantID int `gorm:"not null"`
    PartyName string `gorm:"size:50; not null"`
    PartySize int `gorm:"not null"`
    TimeCreated time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
    PhoneAhead bool `gorm:"not null"`
    WaitTimeExpected int
    WaitTimeCalculated int
    BuzzerID int
}

type HistoricalParty struct {
    ID int
    RestaurantID int `gorm:"not null"`
    PartyName string `gorm:"size:50;not null"`
    PartySize int `gorm:"not null"`
    TimeCreated time.Time `gorm:"not null"`
    TimeSeated time.Time `gorm:"not null"`
    WaitTimeExpected int `gorm:"not null"`
    WaitTimeCalculated int `gorm:"not null"`
}

type User struct {
    ID int
    RestaurantID int `gorm:"not null"`
    Username string `gorm:"size:100;not null unique"`
    Password string `gorm:"size:100; not null"`
    PassSalt string `gorm:"size:50; not null"`
    DateCreated time.Time `gorm:"not null" sql:"DEFAULT:current_timestamp"`
}
