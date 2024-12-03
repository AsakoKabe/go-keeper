package main

import (
	"fmt"
	"log"

	"go-keeper/config/server"
	"go-keeper/internal/http/rest"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
)

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		return
	}

	app, err := rest.NewApp(cfg)

	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer app.Stop()

	if err := app.Run(cfg); err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func printBuildInfo() {
	fmt.Println("Build version: " + buildVersion)
	fmt.Println("Build date: " + buildDate)
}
