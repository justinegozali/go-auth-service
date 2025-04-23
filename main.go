package handler

import (
	"auth-service/config"
	"auth-service/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}


func init() {
	 config.DatabaseCon()
}

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Hello", port)

	r := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowAllOrigins: true, // Allow all origins
		// Alternatively, you can specify allowed origins:
		// AllowOrigins: []string{"http://example.com", "http://anotherdomain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
	}

	

	r.Use(cors.New(corsConfig))
	routeGroup := r.Group("/auth-service")
	routes.Routes(routeGroup)
	tokenRoutes := r.Group("/token")
	routes.TokenRoutes(tokenRoutes)
	roleRoutes := r.Group("/role-service")
	routes.RoleRoutes(roleRoutes)
	// r.Run()
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// CompileDaemon -command="./dummyservice"