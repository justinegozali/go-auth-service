package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/payloads"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init() {
	config.EnvInit()
}

func CreateMember(c *gin.Context) {
	var member models.Member

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

	var response payloads.CreateMemberResponse
	if err := config.DB.Create(&member).Scan(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create member",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New member successfully added",
		"data":    response,
	})
}

func ShowAllMember(c *gin.Context) {
	var members []models.Member

	if err := config.DB.Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fecth data",
		"data":    members,
	})
}

func UpdateMember(c *gin.Context) {
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
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to fetch member",
				"data":    nil,
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
		"data":    existingMember,
	})
}

func SoftDeleteMember(c *gin.Context) {

	id := c.Param("id")

	var existingMember models.Member

	if err := config.DB.First(&existingMember, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Failed to fetch existing member",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to fetch member data",
				"data":    nil,
			})
		}
		return
	}

	existingMember.IsActive = false

	if err := config.DB.Save(&existingMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete a member",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Member successfully deleted",
		"data":    existingMember,
	})
}

func PaginatedMember(c *gin.Context) {
	var members []models.Member
	var totalMembers int64

	// Get pagination parameters from query string
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	// Set default values for page and limit
	page := 1
	limit := 10

	// Parse page and limit from query parameters
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get total number of members
	if err := config.DB.Model(&models.Member{}).Count(&totalMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching total members count",
			"data":    nil,
		})
		return
	}

	// Fetch paginated members
	if err := config.DB.Offset(offset).Limit(limit).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"data":    nil,
		})
		return
	}

	// Return paginated data
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fetched data",
		"data":    members,
		"pagination": gin.H{
			"current_page":  page,
			"total_pages":   int(math.Ceil(float64(totalMembers) / float64(limit))),
			"total_members": totalMembers,
			"limit":         limit,
		},
	})
}
