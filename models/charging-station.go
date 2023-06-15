package models

import "time"

type ChargingStation struct {
	StationID        uint      `json:"stationID" gorm:"primaryKey"`
	EnergyOutput     string    `json:"energyOutput" validate:"required"`
	Type             string    `json:"type" validate:"required"`
	IsOccupied       bool      `json:"status"`
	AvailabilityTime time.Time `json:"availabilityTime"`
}
