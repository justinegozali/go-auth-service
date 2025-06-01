package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/payloads"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	config.EnvInit()
}

// Register pengguna anyar
func UserCreate(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// log.Printf(user.UserName)

	// Nggoleki pengguna sing wis ana
	var existingUser models.User
	if err := config.DB.Where("user_name = ?", user.UserName).First(&existingUser).Error; err == nil {

		c.JSON(http.StatusConflict, gin.H{"error": "User with this username already exists"})
		return
	}

	// Ngubah sandi dadi hash
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed hash password",
		})
		return
	}

	user.Password = string(hash)

	// Lebokno pengguna anyar menyang db
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Authenticate penggguna
func Authenticate(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Nggoleki user nang jero DB
	var existingUser models.User
	if err := config.DB.Where("user_name = ?", user.UserName).First(&existingUser).Error; err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	if err := config.DB.Where("is_logged_in = ?", user.IsLoggedIn).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is already logged in"})
		return
	}

	// Bandingke password sing di lebokno karo sing sampun di dekrip teng db
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Ngubah isloggedin dadi true
	existingUser.IsLoggedIn = true
	// if err := config.DB.Save(&existingUser).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update login status"})
	// 	return
	// }

	// Nggawe token jwt

	// Access token
	accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": existingUser.ID,
		"exp": time.Now().Add(30 * time.Minute).Unix(),
	})

	// Refresh token
	// refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": existingUser.ID,
	// 	"exp": time.Now().Add(1 * time.Hour).Unix(),
	// })

	// Nandhatangani akses token
	secretSign := os.Getenv("SECRET")
	tokenString, err := accesstoken.SignedString([]byte(secretSign))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token" + err.Error(),
		})
		return
	}

	// signedRefreshtoken, err := refreshToken.SignedString([]byte(secretSign))

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Failed to create token" + err.Error(),
	// 	})
	// 	return
	// }

	// kirim respon
	c.JSON(http.StatusOK, gin.H{
		"accessToken": tokenString,
		// "refreshToken": signedRefreshtoken,
		"userID": existingUser.ID,
		"roleID": existingUser.RoleId,
	})
}

// Nampilake pengguna kabeh
func ShowAllUser(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Logout pengguna
func Logout(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Nggoleki user miturut Id
	var currentUser models.User
	if err := config.DB.First(&currentUser, user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User  not found"})
		return
	}

	// Ngubah status login dadi false
	currentUser.IsLoggedIn = false
	if err := config.DB.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User is successfully logged out",
	})
}

func GetUserRoleViews(c *gin.Context) {
	var userRoleViews []models.UserRoleView

	if err := config.DB.Table("user_role_view").Find(&userRoleViews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user role views: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success fetching use view",
		"data":    userRoleViews,
	})
}

func UpdateUser(c *gin.Context) {
	var user payloads.UpdateUserRequest

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	is_change_password := c.Query("is_change_password")

	if is_change_password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Need to add change password",
			"data":    nil,
		})
		return
	}

	var existingUser models.User
	if err := config.DB.First(&existingUser, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is not valid"})
		return
	}

	if is_change_password == "false" {

		log.Println(existingUser)

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed hash password",
			})
			return
		}

		user.Password = string(hash)

		if err := config.DB.Model(&existingUser).Updates(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update user",
				"data":    nil,
			})
			return
		}
	}

	if is_change_password == "true" {

		log.Println(existingUser)

		hash, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed hash password",
			})
			return
		}

		user.Password = string(hash)

		if err := config.DB.Model(&existingUser).Updates(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update user",
				"data":    nil,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successflly update user",
		"data":    existingUser,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var existingUser models.User

	if err := config.DB.First(&existingUser, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Failed to fetch user data",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to fetch user data",
				"data":    nil,
			})
		}
		return
	}

	if err := config.DB.Where("id = ?", id).Delete(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"messae": "Failed to delete user",
			"data":   nil,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"data":    existingUser,
	})

}
