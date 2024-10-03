package main

import (
	"fabric-rest-api/controller"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Set up API endpoints
	router.POST("/addProduct", controller.CreateProductHandler)
	router.POST("/dispatchProduct", controller.SupplyProductHandler)
	router.POST("/bulkPurchase", controller.WholesaleProductHandler)
	router.GET("/fetchProduct", controller.QueryProductHandler)
	router.POST("/processSale", controller.SellProductHandler)

	// Launch the web server
	fmt.Println("Server is running at http://localhost:3000")
	if err := router.Run("localhost:3000"); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}

