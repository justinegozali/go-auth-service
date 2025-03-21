package routes

import (
	"auth-service/controllers"
	"auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup){
	r.POST("/regist-user", controllers.UserCreate)
	r.POST("/login", controllers.Authenticate)
	r.POST("/logout", controllers.Logout)
	r.GET("/user", middleware.ValidateToken ,controllers.ShowAllUser)
}

func TokenRoutes(r *gin.RouterGroup){
	r.POST("/refresh-access-token", controllers.RefreshToken)
}