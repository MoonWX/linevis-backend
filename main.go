package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"linevis-backend/database"
	"linevis-backend/routes"
	"log"
)

func main() {
	fmt.Println("Init")

	gin.SetMode(gin.ReleaseMode)

	db := database.InitDB("linevis.db")
	r := gin.Default()
	routes.SetupRoutes(r, db)

	log.Fatal(r.Run(":9999"))
}
