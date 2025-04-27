package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvInit() {
	if os.Getenv("VERCEL") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Error loading .env file, proceeding with system environment variables")
		} else {
			log.Println(".env file loaded successfully")
		}
	} else {
		log.Println("Running on Vercel, using system environment variables")
	}
}
