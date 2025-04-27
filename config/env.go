package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvInit() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading env file")
	// }

	if os.Getenv("VERCEL") == "" {
		// Running locally → load .env file
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Error loading .env file, proceeding with system environment variables")
		} else {
			log.Println(".env file loaded successfully")
		}
	} else {
		// Running on Vercel → DO NOT load .env
		log.Println("Running on Vercel, using system environment variables")
	}
}
