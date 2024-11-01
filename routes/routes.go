package routes

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"linevis-backend/database"
	"linevis-backend/service"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Set release mode
	gin.SetMode(gin.ReleaseMode)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许的域名
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

	r.POST("/manual/:id", func(c *gin.Context) {
		var product database.Product
		var input database.Product
		var fileService *service.FileService

		if err := db.First(&product, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		input = product

		fileService, err := service.NewFileService("manuals")
		if err != nil {
			log.Fatalf("Failed to initialize file service: %v", err)
		}

		if !strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Content-Type must be multipart/form-data",
			})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Failed to get file: %v", err),
			})
			return
		}

		if file.Size == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "File is empty",
			})
			return
		}

		result, err := fileService.SaveUploadedFileWithName(file, service.UUIDPrefix, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		input.Manual = result.FileName

		if err := db.Model(&product).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded successfully",
			"path":     result.FullPath,
			"filename": result.FileName,
			"size":     result.Size,
		})
	})

	r.GET("/manual/:id", func(c *gin.Context) {
		var fileService *service.FileService
		// 1. 先从数据库获取产品信息
		var product database.Product
		if err := db.First(&product, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		// 2. 检查是否有手册文件
		if product.Manual == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "no manual file found"})
			return
		}

		// 3. 初始化文件服务
		fileService, err := service.NewFileService("manuals")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize file service"})
			return
		}

		// 4. 获取文件路径
		filePath := fileService.GetFilePath(product.Manual)

		// 5. 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "manual file not found"})
			return
		}

		// 6. 设置响应头，指定为图片
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "inline")

		// 根据文件扩展名设置正确的 Content-Type
		ext := strings.ToLower(filepath.Ext(product.Manual))
		switch ext {
		case ".jpg", ".jpeg":
			c.Header("Content-Type", "image/jpeg")
		case ".png":
			c.Header("Content-Type", "image/png")
		case ".gif":
			c.Header("Content-Type", "image/gif")
		case ".webp":
			c.Header("Content-Type", "image/webp")
		default:
			c.Header("Content-Type", "image/jpeg") // 默认假设为 JPEG
		}

		// 7. 发送文件
		c.File(filePath)
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
