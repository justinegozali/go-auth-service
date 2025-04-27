package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvInit() {
	if os.Getenv("VERCEL") == "" {
		// Local machine: load .env file
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Error loading .env file, proceeding with system environment variables")
		} else {
			log.Println(".env file loaded successfully")
		}
	} else {
		// Production (Vercel): don't load .env, just use system env
		log.Println("Running on Vercel, using system environment variables")
	}
}
