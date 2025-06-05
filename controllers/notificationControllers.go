package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func init() {
	config.EnvInit()
}

type RequestBody struct {
	Results   string `json:"results"`
	ScannedAt string `json:"scanned_at"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

// websocket handler
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}
	}
}

func NotifyClients(message string) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func FetchDataFromPython(c *gin.Context) {
	var requestBody RequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	fmt.Println("reqq", requestBody)

	response := gin.H{
		"scannedAt":    requestBody.ScannedAt,
		"licensePlate": requestBody.Results,
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		NotifyClients("New data received: (failed to marshal response)")
		return
	}

	NotifyClients(fmt.Sprintf("New data received: %s", &responseJSON))

	// Parse the scannedAt timestamp
	const layout = "2006-01-02T15:04:05.999999"
	scannedAtTime, err := time.Parse(layout, requestBody.ScannedAt)
	if err != nil {
		log.Printf("Failed to parse scannedAt: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse scannedAt"})
		return
	}
	// Create a new LogKendaraan entry
	logKendaraan := models.LogKendaraan{
		NomorPolisi:  requestBody.Results,                // Assuming Results is a slice and you want the first entry
		JamMasuk:     scannedAtTime.Format("15:04:05"),   // Format time as HH:mm:ss
		TanggalMasuk: scannedAtTime.Format("2006-01-02"), // Format date as YYYY-MM-DD
		IsHarian:     false,
	}
	// Save to database
	if err := config.DB.Create(&logKendaraan).Error; err != nil {
		log.Printf("Failed to save log kendaraan to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save data to database"})
		return
	}

	// Create a notification after saving the log
	notification := models.Notification{
		UserId:       1, // Set the appropriate user ID
		NomorPolisi:  logKendaraan.NomorPolisi,
		JamMasuk:     logKendaraan.JamMasuk,
		TanggalMasuk: logKendaraan.TanggalMasuk,
		LogId:        int(logKendaraan.ID), // Fill log_id with the ID of the saved log
	}
	// Save the notification to the database
	if err := config.DB.Create(&notification).Error; err != nil {
		log.Printf("Failed to save notification to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save notification"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func StoreNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindBodyWithJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

func GetNotificationList(c *gin.Context) {
	var notifications []models.Notification
	params := c.Request.URL.Query()

	if err := config.DB.Find(&notifications).Where(&params).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching notification",
			"data":    nil,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification fetched successfully",
		"data":    notifications,
	})
}

func UpdateLogStatus(c *gin.Context) {
	var log models.LogKendaraan

	if err := c.ShouldBindBodyWithJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var existingLog models.LogKendaraan
	if err := config.DB.First(&existingLog, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve log"})
		}
		return
	}

	existingLog.IsHarian = log.IsHarian

	if err := config.DB.Save(&existingLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update log status"})
		return
	}

	c.JSON(http.StatusOK, existingLog)
}

func UpdateNotificationStatus(c *gin.Context) {
	var notification models.Notification

	if err := c.ShouldBindBodyWithJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var existingNotification models.Notification
	if err := config.DB.First(&existingNotification, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notification"})
		}
		return
	}

	existingNotification.IsRead = notification.IsRead

	if err := config.DB.Save(&existingNotification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update log status"})
		return
	}

	c.JSON(http.StatusOK, existingNotification)
}
