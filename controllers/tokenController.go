package controllers

import (
	"auth-service/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func init() {
	config.EnvInit()
}

func RefreshToken(c *gin.Context) {
	// Njuput refresh token saka header
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Misahke header karo token kanggo njuput tokene
	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Dekrip token
	hmacSampleSecret:= os.Getenv("SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Priksa algoritma signing, kudu cocok karo sing dikarepake
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		log.Fatal(err)
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Mriksa apa token sampun kadaluwarsa
		if float64(time.Now().Unix()) > claims["exp"].(float64){
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		sub := claims["sub"]

		// Nggawe access token ingkang anyar
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": sub,
			"exp": time.Now().Add(1 * time.Minute).Unix(),
		})

		// Nandhatangani akses token anyar
		secretSign := os.Getenv("SECRET")
		tokenString, err := newAccessToken.SignedString([]byte(secretSign))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create token" + err.Error(),
			})
			return 
		}

		c.JSON(http.StatusOK, gin.H{
			"newAccessToken": tokenString,
		})

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		fmt.Println(err)
	}
}