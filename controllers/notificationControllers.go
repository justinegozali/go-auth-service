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

	// logFile, err := os.OpenFile("data_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Printf("Failed to open log file: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to open log file"})
	// 	return
	// }
	// defer logFile.Close()

	// logger := log.New(logFile, "", log.LstdFlags)
	// logger.Printf("Scanned at: %s\n", requestBody.ScannedAt)
	// var licensePlates []string
	// for _, result := range requestBody.Results {
	// 	logger.Printf("Detected text: %s, Probability: %.2f\n", result.Text, result.Probability)
	// 	licensePlates = append(licensePlates, result.Text)
	// }

	var licensePlates []models.LogKendaraan
	for _, result := range requestBody.Results {
		licensePlate := models.LogKendaraan{
			NomorPolisi:  result.Text,
			JamMasuk:     time.Now().Format("15:04:05"),   // Current time
			TanggalMasuk: time.Now().Format("2006-01-02"), // Current date
			IsHarian:     false,                           // Set default value
		}
		licensePlates = append(licensePlates, licensePlate)
		// Save to database
		if err := config.DB.Create(&licensePlate).Error; err != nil {
			log.Printf("Failed to save license plate to database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save data to database"})
			return
		}
	}

	var licensePlateNumbers []string
	for _, lp := range licensePlates {
		licensePlateNumbers = append(licensePlateNumbers, lp.NomorPolisi)
	}

	response := gin.H{
		"scannedAt":    requestBody.ScannedAt,
		"licensePlate": strings.Join(licensePlateNumbers, " "),
	}

	NotifyClients(fmt.Sprintf("New data received: %s", response))

	c.JSON(http.StatusOK, response)
}
