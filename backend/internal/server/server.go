package server

import (
	"github.com/gin-gonic/gin"
)

// NewServer creates a new Gin server instance
func NewServer() *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create new router
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	return r
} 