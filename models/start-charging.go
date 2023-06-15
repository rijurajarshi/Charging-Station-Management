package models

import "time"

type StartCharging struct {
	StationID              uint      `json:"stationID" gorm:"primaryKey"`
	VehicleBatteryCapacity string    `json:"vehicleBatteryCapacity" validate:"required"`
	CurrentVehicleCharge   string    `json:"currentVehicleCharge" validate:"required"`
	ChargingStartTime      time.Time `json:"chargingStartTime"`
}
