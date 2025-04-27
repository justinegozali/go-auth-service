package handler

import (
	"auth-service/config"
	"auth-service/routes"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	config.DatabaseCon()
	config.EnvInit()

	router = gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
	}
	router.Use(cors.New(corsConfig))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Golang Backend")
	})

	routeGroup := router.Group("/auth-service")
	routes.Routes(routeGroup)

	tokenRoutes := router.Group("/token")
	routes.TokenRoutes(tokenRoutes)

	roleRoutes := router.Group("/role-service")
	routes.RoleRoutes(roleRoutes)

	memberRoutes := router.Group("/member-service")
	routes.MemberRoutes(memberRoutes)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
