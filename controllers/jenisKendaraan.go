package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	config.EnvInit()
}

func CreateJenisKendaraan(c *gin.Context) {
	var request payloads.CreateKendaraan

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingKendaraan models.Kendaraan
	if err := config.DB.Where("jenis_kendaraan = ?", request.JenisKendaraan).First(&existingKendaraan).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Jenis Kendaraan is already exist"})
		return
	}

	var response payloads.CreateKendaraan
	kendaraan := models.Kendaraan{
		JenisKendaraan: request.JenisKendaraan,
	}

	if err := config.DB.Create(&kendaraan).Scan(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create jenis kendaraan"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func ShowAllJenisKendaraan(c *gin.Context) {
	var Kendaraan []models.Kendaraan
	if err := config.DB.Find(&Kendaraan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"data":    nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Jenis kendaaran successfully fetched",
		"data":    Kendaraan,
	})
}
