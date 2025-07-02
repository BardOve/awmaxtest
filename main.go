package main

import (
	"awmaxtest/internal"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type GenerateMeasurementRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")
	router.StaticFile("/favicon.ico", "./web/static/favicon.png")
 
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "Interactive Go Chart Generator"})
	})

	router.POST("/generate-measurement", func(c *gin.Context) {
		var req GenerateMeasurementRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		log.Printf("Received request to generate measurement data for: %s (%s)", req.Name, req.Email)

		if err := internal.GenerateData(req.Name, req.Email); err != nil { 
			log.Printf("Error generating data: %v", err)
			c.JSON(500, gin.H{"success": false, "error": "Failed to generate data"})
			return
		}
		c.JSON(200, gin.H{"success": true})
	})
	
	router.GET("/generator", func(c *gin.Context) {
		email := c.DefaultQuery("email", "Guest")
		name := c.DefaultQuery("name", "User")
		c.HTML(200, "generator.html", gin.H{
			"title":   "AWMax Data Generator",
			"userEmail": email,
			"userName":  name,
		})
	})
 
	router.POST("/generate-chart", internal.GenerateChartHandler)
	

	router.GET("/charts/:name/:filename", func(c *gin.Context) {
		name := c.Param("name")
		filename := c.Param("filename")

		if strings.Contains(name, "..") || strings.Contains(filename, "..") {
			c.String(http.StatusBadRequest, "Invalid path")
			return
		}
		filePath := filepath.Join("./charts", name, filename)
		c.File(filePath)
	})

	router.POST("/generatepfd", func (c *gin.Context) {
		var req struct {
			Name string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"success": false, "error": "Invalid request body"})
			return
		}

		fmt.Printf("Received request to generate PDF for: %s (%s)\n", req.Name, req.Email)
		

		if err := internal.GeneratePDF("./charts/"+req.Name+".pdf", req.Name, "./charts/Bård/Bård_Backpressure_2025-07-02.png"); err != nil {
			log.Printf("Error generating PDF: %v", err)
			c.JSON(500, gin.H{"success": false, "error": "Failed to generate PDF"})
			return
		}
		c.JSON(200, gin.H{"success": true})
		
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}

//filepath string, name string, inputPNG []string