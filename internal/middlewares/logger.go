package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		end := time.Now()
		latency := end.Sub(start)

		log.Printf("[%s] %s %s | Status: %d | Latency: %v",
			c.Request.Method,
			path,
			raw,
			c.Writer.Status(),
			latency,
		)
	}
}
