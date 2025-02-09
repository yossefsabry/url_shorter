package main

import (
	"fmt"
	"github.com/yossefsabry/url_shorter/handler"
	"github.com/yossefsabry/url_shorter/store"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery()) // Add Recovery manually
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hey Go URL Shortener !",
		})
	})

	err := r.Run(":4440")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}

