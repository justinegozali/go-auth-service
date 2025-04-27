package main

import (
	"auth-service/config"
	"auth-service/routes"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin" // Import for adapting Gin to serverless
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	config.DatabaseCon()
	config.EnvInit()

	r := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
	}
	r.Use(cors.New(corsConfig))

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello from Golang")
	})

	routeGroup := r.Group("/auth-service")
	routes.Routes(routeGroup)

	tokenRoutes := r.Group("/token")
	routes.TokenRoutes(tokenRoutes)

	roleRoutes := r.Group("/role-service")
	routes.RoleRoutes(roleRoutes)

	memberRoutes := r.Group("/member-service")
	routes.MemberRoutes(memberRoutes)

	// Here we initialize the Gin Lambda adapter
	ginLambda = ginadapter.New(r)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.Proxy(req)
}
