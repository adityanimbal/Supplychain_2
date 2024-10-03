package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract represents the contract for managing products
type SmartContract struct {
    contractapi.Contract
}

// Product holds the details of a product
type Product struct {
    ProductID           string `json:"productID"`
    Name                string `json:"name"`
    Description         string `json:"description"`
    ManufacturingDate   string `json:"manufacturingDate"`
    BatchNumber         string `json:"batchNumber"`
    SupplyDate          string `json:"supplyDate,omitempty"`
    WarehouseLocation   string `json:"warehouseLocation,omitempty"`
    WholesaleDate       string `json:"wholesaleDate,omitempty"`
    WholesaleLocation    string `json:"wholesaleLocation,omitempty"`
    Quantity            int    `json:"quantity,omitempty"`
    Status              string `json:"status"`
}

// AddProduct inserts a new product into the ledger
func (sc *SmartContract) AddProduct(ctx contractapi.TransactionContextInterface, productID, name, description, manufacturingDate, batchNumber string) error {
    newProduct := Product{
        ProductID:         productID,
        Name:              name,
        Description:       description,
        ManufacturingDate: manufacturingDate,
        BatchNumber:       batchNumber,
        Status:            "Created",
    }

    productData, err := json.Marshal(newProduct)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(productID, productData)
}

// RecordSupply updates a product's supply information
func (sc *SmartContract) RecordSupply(ctx contractapi.TransactionContextInterface, productID, supplyDate, warehouseLocation string) error {
    product, err := sc.RetrieveProduct(ctx, productID)
    if err != nil {
        return err
    }

    product.SupplyDate = supplyDate
    product.WarehouseLocation = warehouseLocation
    product.Status = "Supplied"

    updatedProductData, _ := json.Marshal(product)
    return ctx.GetStub().PutState(productID, updatedProductData)
}

// RecordWholesale updates a product's wholesale information
func (sc *SmartContract) RecordWholesale(ctx contractapi.TransactionContextInterface, productID, wholesaleDate, wholesaleLocation string, quantity int) error {
    product, err := sc.RetrieveProduct(ctx, productID)
    if err != nil {
        return err
    }

    product.WholesaleDate = wholesaleDate
    product.WholesaleLocation = wholesaleLocation
    product.Quantity = quantity
    product.Status = "Wholesaled"

    updatedProductData, _ := json.Marshal(product)
    return ctx.GetStub().PutState(productID, updatedProductData)
}

// RetrieveProduct fetches a product's details from the ledger
func (sc *SmartContract) RetrieveProduct(ctx contractapi.TransactionContextInterface, productID string) (*Product, error) {
    productData, err := ctx.GetStub().GetState(productID)
    if err != nil {
        return nil, fmt.Errorf("error reading from world state: %s", err.Error())
    }
    if productData == nil {
        return nil, fmt.Errorf("product %s not found", productID)
    }

    var product Product
    if err := json.Unmarshal(productData, &product); err != nil {
        return nil, err
    }

    return &product, nil
}

// ChangeProductStatus modifies the status of a product (e.g., sold)
func (sc *SmartContract) ChangeProductStatus(ctx contractapi.TransactionContextInterface, productID, newStatus string) error {
    product, err := sc.RetrieveProduct(ctx, productID)
    if err != nil {
        return err
    }

    product.Status = newStatus
    updatedProductData, _ := json.Marshal(product)
    return ctx.GetStub().PutState(productID, updatedProductData)
}

