package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func CekKupon(c *gin.Context) {
	var struk []models.StrukMember

	kodeKupon := c.Query("kode_kupon")
	nomorPolisi := c.Query("nomor_polisi")

	if kodeKupon == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "kode kupon is required",
			"data":    nil,
		})
		return
	}

	if nomorPolisi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "nomor polisi is required",
			"data":    nil,
		})
		return
	}

	if err := config.DB.Where("kode_kupon = ? AND nomor_polisi = ?", kodeKupon, nomorPolisi).Find(&struk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching struk",
			"data":    nil,
		})
		return
	}

	if len(struk) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Kupon not found",
			"data":    nil,
		})
		return
	}

	for _, s := range struk {
		// Parse the KadaluarsaKupon string to time.Time
		kadaluarsa, err := time.Parse("2006-01-02", s.KadaluarsaKupon) // Adjust the layout as needed
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error parsing expiration date",
				"data":    nil,
			})
			return
		}
		if kadaluarsa.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Kupon has expired",
				"data":    nil,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kupon fetched successfully",
		"data":    struk,
	})
}

func EditStruk(c *gin.Context) {
	var struk models.StrukMember

	if err := c.ShouldBindBodyWithJSON(&struk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var existingStruk models.StrukMember
	if err := config.DB.First(&existingStruk, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Struk not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve struk"})
		}
		return
	}

	if err := config.DB.Model(&existingStruk).Updates(struk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update struk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Struk updated successfully", "data": existingStruk})

}
