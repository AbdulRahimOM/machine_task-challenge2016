package main

import (
	"challenge16/internal/config"
	"challenge16/internal/regions"
	"challenge16/internal/server"
	"fmt"
)

const (
	csvFile = "cities.csv"
	envPath = ".env"
)

func main() {
	//initialize the region data
	regions.LoadDataIntoMap(csvFile)

	//initialize the environment configuration
	config.LoadEnv(envPath)

	app := server.NewServer()

	err := app.Listen(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		panic("Couldn't start the server. Error: " + err.Error())
	}
}