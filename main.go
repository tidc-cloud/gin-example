package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware to log requests and responses
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		log.Printf("Request: %s %s %s\n", c.Request.Method, c.Request.URL, string(bodyBytes))

		// Record response
		w := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Log response
		log.Printf("Response: %d %s %s\n", w.Status(), http.StatusText(w.Status()), w.body.String())
		log.Printf("Request processed in %v\n", duration)
	}
}

// responseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Logger())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/echo", func(c *gin.Context) {
		var json map[string]interface{}
		if err := c.BindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, json)
	})

	if err := r.Run(); err != nil {
		log.Fatalf("Could not run the server: %v", err)
	}
}
