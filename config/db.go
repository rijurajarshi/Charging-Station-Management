package config

import (
	"github.com/rijurajarshi/charging-station-management/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	dsn := "root:GIVE_YOUR_DB_PASSWORD@tcp(127.0.0.1:3306)/charging_station?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.ChargingStation{})
	db.AutoMigrate(&models.StartCharging{})
	DB = db
}
