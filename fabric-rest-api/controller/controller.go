package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// Paths for wallet and connection configuration
const (
	walletDir = "path/to/wallet"           // Update with your wallet directory
	ccpFile   = "path/to/connection.yaml"   // Update with your connection profile file path
)

// Structs for request payloads
type ProductCreationRequest struct {
	ID   string `json:"productID"`
	Name string `json:"productName"`
}

type ProductSupplyRequest struct {
	ID     string `json:"productID"`
	Status string `json:"status"`
}

type ProductWholesaleRequest struct {
	ID     string `json:"productID"`
	Status string `json:"status"`
}

type ProductQueryRequest struct {
	ID string `json:"productID"`
}

type ProductSaleRequest struct {
	ID       string `json:"productID"`
	Buyer    string `json:"buyerInfo"`
}

// Establish a connection to the Fabric network
func establishGatewayConnection(channel, userID string) (*gateway.Gateway, *gateway.Network, error) {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	wallet, err := gateway.NewFileSystemWallet(walletDir)
	if err != nil {
		fmt.Printf("Unable to create wallet: %s\n", err)
		os.Exit(1)
	}

	gatewayInstance, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpFile))),
		gateway.WithIdentity(wallet, userID),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to gateway: %w", err)
	}

	networkInstance, err := gatewayInstance.GetNetwork(channel)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get network: %w", err)
	}

	return gatewayInstance, networkInstance, nil
}

// Handle POST /createProduct
func CreateProductHandler(c *gin.Context) {
	var req ProductCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request payload"})
		return
	}

	gatewayInstance, network, err := establishGatewayConnection("Channel2", "ProducerUser")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gatewayInstance.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.SubmitTransaction("createProduct", req.ID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// Handle POST /supplyProduct
func SupplyProductHandler(c *gin.Context) {
	var req ProductSupplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request payload"})
		return
	}

	gatewayInstance, network, err := establishGatewayConnection("Channel2", "SupplierUser")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gatewayInstance.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.SubmitTransaction("supplyProduct", req.ID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// Handle POST /wholesaleProduct
func WholesaleProductHandler(c *gin.Context) {
	var req ProductWholesaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request payload"})
		return
	}

	gatewayInstance, network, err := establishGatewayConnection("Channel3", "WholesalerUser")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gatewayInstance.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.SubmitTransaction("wholesaleProduct", req.ID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// Handle GET /queryProduct
func QueryProductHandler(c *gin.Context) {
	productID := c.Query("productID")

	gatewayInstance, network, err := establishGatewayConnection("Channel1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gatewayInstance.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.EvaluateTransaction("queryProduct", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": string(result)})
}

// Handle POST /sellProduct
func SellProductHandler(c *gin.Context) {
	var req ProductSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request payload"})
		return
	}

	gatewayInstance, network, err := establishGatewayConnection("Channel1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gatewayInstance.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.SubmitTransaction("sellProduct", req.ID, req.Buyer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

