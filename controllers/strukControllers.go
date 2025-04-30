package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init() {
	config.EnvInit()
}

func CreateStruk(c *gin.Context) {
	var struk models.StrukMember

	if err := c.ShouldBindBodyWithJSON(&struk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&struk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating struk",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Struk created successfully",
		"data":    struk,
	})
}
