package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"linevis-backend/database"
	"net/http"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:1420"}, // 允许的域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Create
	r.POST("/products", func(c *gin.Context) {
		var product database.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&product)
		c.JSON(http.StatusOK, product)
	})

	// Read
	r.GET("/products/:id", func(c *gin.Context) {
		var product database.Product
		if err := db.First(&product, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	// Read All
	r.GET("/products", func(c *gin.Context) {
		var products []database.Product
		if err := db.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error when fetching products"})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	// Update
	r.PUT("/products/:id", func(c *gin.Context) {
		var product database.Product
		if err := db.First(&product, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		err := c.ShouldBindJSON(&product)
		if err != nil {
			return
		}
		db.Save(&product)
		c.JSON(http.StatusOK, product)
	})

	// Delete
	r.DELETE("/products/:id", func(c *gin.Context) {
		if err := db.Delete(&database.Product{}, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
