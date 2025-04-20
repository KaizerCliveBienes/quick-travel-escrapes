package main

import (
	"findeventsberlin/summarizer"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	router := gin.Default()
	router.GET("/events/berlin", func(context *gin.Context) {
		month := context.DefaultQuery("month", strings.ToLower(time.Now().Month().String()))
		summary, err := summarizer.ScrapeBerlinEvents(month)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Unable to %v", err),
			})
		}

		context.JSON(http.StatusOK, summary)
	})

	router.Run(":8080")
}
