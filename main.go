package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/willianVini-dev/hexagonal/adapter/input/routes"
	"github.com/willianVini-dev/hexagonal/configuration/logger"
)

func main() {

	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", err)
	}

	logger.Info("About to start application")

	router := gin.Default()
	routes.InitRoutes(router)

	if err := router.Run(":8080"); err != nil {
		logger.Error("Error trying to start application", err)
		return
	}
}
