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
    Restaurant []restaurants
    name = string `gorm:"size:255"`
    last_heartbeat time.Time
    is_active bool
    ActiveParty []activePartyIds

}
