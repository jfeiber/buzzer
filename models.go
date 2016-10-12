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

type Buzzer struct {
    gorm.Model

    buzzerID uint
    Restaurant restaurant
    name string `gorm:"size:45"`
    last_heartbeat time.Time
    is_active bool
    ActiveParty activePartyID

}

type ActiveParty struct {
    activePartyID uint
    Restaurant restaurant
    partyName string `gorm:"size:100"`
    partySize string `gorm:"size:100"`
    timeCreated time.Time
    timeSeated time.Time
    phoneAhead bool
    waitTimeExpected time.Time
    waitTimeCalculated time.Time
}

type Restaurant struct {
    restaurantID uint
    name string `gorm:"size:99"`
    dateCreated time.Time
}
