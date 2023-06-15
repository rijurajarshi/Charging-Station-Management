package controller

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"github.com/rijurajarshi/charging-station-management/config"
	"github.com/rijurajarshi/charging-station-management/models"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}

	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	logger.Println("Application started.....")

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		logger.Println("Application closed!!")
		os.Exit(0)
	}()

}

var validation *validator.Validate = validator.New()

var catch *cache.Cache = cache.New(2*time.Minute, 3*time.Minute)

func AddChargingStation(c *gin.Context) {
	var station models.ChargingStation
	err := c.ShouldBindJSON(&station)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	err = validation.Struct(station)
	if err != nil {
		logger.Println("Validation error:", err)
		c.JSON(400, err.Error())
		c.Abort()
		return
	}

	station.AvailabilityTime = time.Now()
	err = config.DB.Create(&station).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create charging station"})
		return
	}
	logger.Println("Charging station added successfully")
	c.JSON(201, gin.H{"message": "Charging station added successfully"})
}

func StartCharging(c *gin.Context) {
	var chargingRequest models.StartCharging
	if err := c.ShouldBindJSON(&chargingRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := validation.Struct(chargingRequest)
	if err != nil {
		logger.Println("Validation error:", err)
		c.JSON(400, err.Error())
		c.Abort()
		return
	}
	chargingRequest.ChargingStartTime = time.Now()
	config.DB.Create(&chargingRequest)
	c.JSON(200, gin.H{"success": chargingRequest})

	var chargingStation models.ChargingStation
	result := config.DB.First(&chargingStation, chargingRequest.StationID)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Charging station not found"})
		return
	}
	chargingStation.IsOccupied = true
	config.DB.Save(&chargingStation)
	availabilityTime := CalculateAvailabilityTime(chargingRequest, chargingStation)
	result = config.DB.Model(&chargingStation).Update("availability_time", availabilityTime)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to start charging"})
		return
	}
	logger.Println("Charging started successfully")
	c.JSON(200, gin.H{"message": "Charging started successfully"})
}

func GetAvailableChargingStations(c *gin.Context) {
	var chargingStations []models.ChargingStation

	result := config.DB.Where("is_occupied=?", false).Find(&chargingStations)
	if result.Error != nil {

		c.JSON(500, gin.H{"error": "Failed to retrieve available charging stations"})
		return
	}
	logger.Println("Retrieved available charging stations")
	c.JSON(200, chargingStations)

}

func GetOccupiedChargingStations(c *gin.Context) {
	var chargingStations []models.ChargingStation
	result := config.DB.Where("is_occupied=?", true).Find(&chargingStations)
	if result.Error != nil {

		c.JSON(500, gin.H{"error": "Failed to retrieve available charging stations"})
		return
	}
	logger.Println("Retrieved occupied charging stations")
	c.JSON(200, chargingStations)
}

func GetAllChargingStations(c *gin.Context) {
	cache_key := "charging-station"

	fmt.Println(catch.ItemCount())

	if result, found := catch.Get(cache_key); found {
		logger.Println("Retrieved all charging stations from cache")
		c.JSON(200, gin.H{"chargingStations from cache": result})
		return
	} else {
		var chargingStations []models.ChargingStation
		config.DB.Find(&chargingStations)
		catch.Set(cache_key, chargingStations, cache.DefaultExpiration)
		logger.Println("Accessing all Charging Stations using Database")
		c.JSON(200, gin.H{
			"All Charging Stations from database": chargingStations,
		})
	}
}

func GetChargingStationByID(c *gin.Context) {

	cache_key1 := "charging-station-byID"
	id := c.Param("id")

	if result, found := catch.Get(cache_key1); found {
		logger.Println("Retrieved charging station by ID from cache")
		c.JSON(200, gin.H{"chargingStation By ID from cache": result})
		return
	} else {

		var chargingStations []models.ChargingStation

		config.DB.Find(&chargingStations, id)
		catch.Set(cache_key1, chargingStations, cache.DefaultExpiration)
		logger.Println("Accessing Charging Station By ID using Database")
		c.JSON(200, gin.H{
			"Charging Station By ID from database": chargingStations,
		})
	}

}

func CalculateAvailabilityTime(request models.StartCharging, chargingStation models.ChargingStation) *time.Time {
	vehicleBatteryCapacity, _ := strconv.Atoi(strings.ReplaceAll(request.VehicleBatteryCapacity, "kWh", ""))
	currentVehicleCharge, _ := strconv.Atoi(strings.ReplaceAll(request.CurrentVehicleCharge, "kWh", ""))
	energyOutput, _ := strconv.Atoi(strings.ReplaceAll(chargingStation.EnergyOutput, "kWh", ""))

	remainingEnergy := time.Duration(vehicleBatteryCapacity-currentVehicleCharge) * time.Hour
	availabilityTime := request.ChargingStartTime.Add(remainingEnergy / time.Duration(energyOutput))

	return &availabilityTime
}
