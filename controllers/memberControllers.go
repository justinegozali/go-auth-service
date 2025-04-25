package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init() {
	config.EnvInit()
}

func CreateMember(c *gin.Context){
	var member models.Member

	log.Printf("data: %+v\n", member)

	if err := c.ShouldBindBodyWithJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingPlat models.Member
	if err := config.DB.Where("nomor_polisi = ?", member.NomorPolisi).First(&existingPlat).Error; err == nil {
			
		c.JSON(http.StatusConflict, gin.H{"error": "Nopol is already exists"})
		return
	}

	if member.TanggalMasuk == "" {
		member.TanggalMasuk = time.Now().Format("2006-01-02") 
	}


	if err := config.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create member",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New member successfully added",
		"data": member, 
	})
}

func ShowAllMember(c *gin.Context){
	var members []models.Member

	if err := config.DB.Find(&members).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fecth data",
		"data": members,
	})
}

func UpdateMember(c *gin.Context){
	var member models.Member

	if err := c.ShouldBindBodyWithJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var existingMember models.Member
	if err := config.DB.First(&existingMember, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Member not found",
				"data": nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to fetch member",
				"data": nil,
			})
		}
		return
	}

	if err := config.DB.Model(&existingMember).Updates(member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update member",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successflly update member",
		"data": existingMember,
	})
}
