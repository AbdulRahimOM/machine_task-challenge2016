package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultPort = "4010"
)

var (
	Port string
)

func init() {
	loadEnv()
}

func loadEnv() {
	fmt.Println("Loading .env file...")
	//parse .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. err", err)
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = defaultPort
	}

}
