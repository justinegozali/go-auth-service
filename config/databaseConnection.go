package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB  *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file")
	}
}

func DatabaseCon(){
	var err error
	dsn := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	fmt.Println("DSN",dsn)
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("Connect to db", dsn)
	}
	
}

