package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	defaultPort = "4010"
)

var (
	Port      string
	RateLimit int
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

	RateLimit,err = strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		if os.Getenv("RATE_LIMIT") == "" {
			RateLimit = 60
		} else {
			log.Fatal("Error loading RATE_LIMIT from .env file. err", err)
		}
	}

}
