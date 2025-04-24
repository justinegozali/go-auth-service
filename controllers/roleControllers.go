package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/payloads"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}
}

func CreateRole(c *gin.Context) {
	var request payloads.CreateRoleRequest
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingRole models.Role
	if err := config.DB.Where("name = ?", request.Name).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Role already exist"})
		return
	}

	var response payloads.CreateRoleResponse
	role := models.Role{
		Name: request.Name,
	}
	if err := config.DB.Create(&role).Scan(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func EditRole(c *gin.Context) {
	var role models.Role

	if err := c.ShouldBindBodyWithJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var existingRole models.Role
	if err := config.DB.First(&existingRole, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve role"})
		}
		return
	}

	existingRole.Name = role.Name

	if err := config.DB.Save(&existingRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	// Return the updated role
	c.JSON(http.StatusOK, existingRole)
}
