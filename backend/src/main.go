package main

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return ":" + port
}

func main() {
	port := GetPort()

	if port != "3000" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()

	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})

	engine.Run(port)
}
