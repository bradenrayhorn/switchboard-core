package main

import (
	"github.com/bradenrayhorn/switchboard-core/config"
	"github.com/bradenrayhorn/switchboard-core/database"
	"github.com/bradenrayhorn/switchboard-core/routing"
	"log"
)

func main() {
	log.Printf("starting switchboard...")

	log.Printf("loading config...")
	config.LoadConfig()
	log.Printf("config loaded!")

	log.Printf("initializing database...")
	database.Setup()
	log.Printf("database ready!")

	startServer()
}

func startServer() {
	r := routing.MakeRouter()

	err := r.Run()

	if err != nil {
		panic(err)
	}
}
