package main

import (
	"auth-service/config"
	"auth-service/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	 config.DatabaseCon()
}

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  envtest := os.Getenv("PORT")

	fmt.Println("Hello", envtest)

	r := gin.Default()

	routeGroup := r.Group("/auth-service")
	routes.Routes(routeGroup)
	tokenRoutes := r.Group("/token")
	routes.TokenRoutes(tokenRoutes)
	r.Run()
}

// CompileDaemon -command="./dummyservice"