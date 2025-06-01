package routes

import (
	"auth-service/controllers"
	"auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.POST("/regist-user", controllers.UserCreate)
	r.POST("/login", controllers.Authenticate)
	r.POST("/logout", middleware.ValidateToken, controllers.Logout)
	r.GET("/user", middleware.ValidateToken, controllers.ShowAllUser)
	r.GET("/user-view", middleware.ValidateToken, controllers.GetUserRoleViews)
	r.PUT("/user/:id", middleware.ValidateToken, controllers.UpdateUser)
	r.DELETE("/user/:id", middleware.ValidateToken, controllers.DeleteUser)
}

func TokenRoutes(r *gin.RouterGroup) {
	r.POST("/refresh-access-token", controllers.RefreshToken)
}

func RoleRoutes(r *gin.RouterGroup) {
	r.POST("/role", controllers.CreateRole)
	r.GET("/role", middleware.ValidateToken, controllers.ShowAllRole)
	r.PUT("/role/:id", middleware.ValidateToken, controllers.EditRole)
}

func MemberRoutes(r *gin.RouterGroup) {
	r.POST("/member", middleware.ValidateToken, controllers.CreateMember)
	r.GET("/member", middleware.ValidateToken, controllers.ShowAllMember)
	r.PUT("/member/:id", middleware.ValidateToken, controllers.UpdateMember)
	r.PUT("/delete-member/:id", middleware.ValidateToken, controllers.SoftDeleteMember)
	r.GET("/paginated-member", middleware.ValidateToken, controllers.PaginatedMember)
}

func StrukRoutes(r *gin.RouterGroup) {
	r.POST("/struk", middleware.ValidateToken, controllers.CreateStruk)
	r.GET("/kupon", controllers.CekKupon)
	r.PUT("/struk/:id", controllers.EditStruk)
}

func OcrRoutes(r *gin.RouterGroup) {
	r.POST("/fetch-ocr-data", controllers.FetchDataFromPython)
}

func NotificationRoutes(r *gin.RouterGroup) {
	r.GET("/ws", controllers.HandleWebSocket)
	r.POST("/fetch-data", controllers.FetchDataFromPython)
	r.POST("/notification-storage", controllers.StoreNotification)
	r.GET("/notification-storage", controllers.GetNotificationList)
}

func JenisKendaraanRoutes(r *gin.RouterGroup) {
	r.POST("/jenis-kendaraan", controllers.CreateJenisKendaraan)
	r.GET("/jenis-kendaraan", controllers.ShowAllJenisKendaraan)
}
