package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rijurajarshi/charging-station-management/config"
	"github.com/rijurajarshi/charging-station-management/routes"
)

func main() {

	router := gin.Default()

	routes.ChargingStationRoute(router)

	config.DBConnect()

	log.Fatal(router.Run(":8080"))
}
