package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func init() {
	config.EnvInit()
}

type DetectedText struct {
	Text        string  `json:"text"`
	Probability float64 `json:"probability"`
}

type RequestBody struct {
	Results   []DetectedText `json:"results"`
	ScannedAt string         `json:"scanned_at"`
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

	// Initialize a slice to hold license plates
	var licensePlates []models.LogKendaraan
	var licensePlateNumbers []string // To hold the license plate numbers for logging

	for _, result := range requestBody.Results {
		licensePlate := models.LogKendaraan{
			NomorPolisi:  result.Text,
			JamMasuk:     time.Now().Format("15:04:05"),   // Current time
			TanggalMasuk: time.Now().Format("2006-01-02"), // Current date
			IsHarian:     false,                           // Set default value
		}
		licensePlates = append(licensePlates, licensePlate)
		licensePlateNumbers = append(licensePlateNumbers, result.Text) // Collect license plate numbers

		// Save to database
		if err := config.DB.Create(&licensePlate).Error; err != nil {
			log.Printf("Failed to save license plate to database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save data to database"})
			return
		}
	}

	// Join license plate numbers into a single string
	joinedLicensePlates := strings.Join(licensePlateNumbers, " ")

	// Log the joined license plate numbers
	log.Printf("Detected license plates: %s", joinedLicensePlates)

	response := gin.H{
		"scannedAt":    requestBody.ScannedAt,
		"licensePlate": joinedLicensePlates,
	}

	NotifyClients(fmt.Sprintf("New data received: %s", response))

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
