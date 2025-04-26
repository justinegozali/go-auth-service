package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}

func GinHandler(c *gin.Context) {
	// Create a ResponseWriter and Request from the Gin context
	w := c.Writer
	r := c.Request

	// Call the original Handler
	Handler(w, r)
}

func Main() {
	r := gin.Default()

	r.GET("/", GinHandler)

	http.ListenAndServe(":8080", r)
}
