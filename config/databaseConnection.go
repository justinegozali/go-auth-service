package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
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

func DatabaseCon() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	fmt.Println("DSN", dsn)
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("Connect to db", dsn)
	}

}
