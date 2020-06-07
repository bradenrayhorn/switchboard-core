package main

import (
	"github.com/bradenrayhorn/switchboard-backend/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Printf("starting switchboard...")

	log.Printf("loading config...")
	config.LoadConfig()
	log.Printf("config loaded!")

	startServer()
}

func startServer() {
	router := gin.Default()

	registerRoutes(router)

	router.Run()
}
