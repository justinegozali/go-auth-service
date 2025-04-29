package routes

import (
	"auth-service/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.POST("/regist-user", controllers.UserCreate)
	r.POST("/login", controllers.Authenticate)
	r.POST("/logout", controllers.Logout)
	// r.GET("/user", middleware.ValidateToken ,controllers.ShowAllUser)
	r.GET("/user", controllers.ShowAllUser)
}

func TokenRoutes(r *gin.RouterGroup) {
	r.POST("/refresh-access-token", controllers.RefreshToken)
}

func RoleRoutes(r *gin.RouterGroup) {
	r.POST("/role", controllers.CreateRole)
	r.PUT("/role/:id", controllers.EditRole)
}

func MemberRoutes(r *gin.RouterGroup) {
	r.POST("/member", controllers.CreateMember)
	r.GET("/member", controllers.ShowAllMember)
	r.PUT("/member/:id", controllers.UpdateMember)
	r.PUT("/delete-member/:id", controllers.SoftDeleteMember)
	r.GET("/paginated-member", controllers.PaginatedMember)
}
