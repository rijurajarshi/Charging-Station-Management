package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rijurajarshi/charging-station-management/controller"
)

func ChargingStationRoute(router *gin.Engine) {
	router.POST("/charging-station", controller.AddChargingStation)
	router.POST("/start-charging", controller.StartCharging)
	router.GET("/available-charging-stations", controller.GetAvailableChargingStations)
	router.GET("/occupied-charging-stations", controller.GetOccupiedChargingStations)
	router.GET("/all-charging-stations", controller.GetAllChargingStations)
	router.GET("/charging-station/:id", controller.GetChargingStationByID)

}
