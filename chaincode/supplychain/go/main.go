// ProductManager offers methods to handle the product lifecycle
type ProductManager struct {
    contractapi.Contract
}

// Product holds essential information about a product in the supply chain
type Product struct {
    ID                string `json:"productID"`
    Name              string `json:"name"`
    Description       string `json:"description"`
    ManufacturingDate string `json:"manufacturingDate"`
    BatchNumber       string `json:"batchNumber"`
    CurrentStatus     string `json:"status"`
    SupplyDate        string `json:"supplyDate"`
    WarehouseLocation  string `json:"warehouseLocation"`
    WholesaleDate     string `json:"wholesaleDate"`
    WholesaleLocation  string `json:"wholesaleLocation"`
    Quantity          int    `json:"quantity"`
}

// InitializeLedger populates the ledger with sample product data
func (pm *ProductManager) InitializeLedger(ctx contractapi.TransactionContextInterface) error {
    sampleProducts := []Product{
        {ID: "P001", Name: "Product1", Description: "Description for Product1", ManufacturingDate: "2023-09-25", BatchNumber: "B001", CurrentStatus: "Created"},
        {ID: "P002", Name: "Product2", Description: "Description for Product2", ManufacturingDate: "2023-09-26", BatchNumber: "B002", CurrentStatus: "Created"},
    }

    for _, item := range sampleProducts {
        itemData, err := json.Marshal(item)
        if err != nil {
            return fmt.Errorf("failed to marshal product: %s", err.Error())
        }
        
        if err := ctx.GetStub().PutState(item.ID, itemData); err != nil {
            return fmt.Errorf("failed to initialize ledger: %s", err.Error())
        }
    }
    return nil
}

func main() {
    chaincodeInstance, err := contractapi.NewChaincode(&ProductManager{})
    if err != nil {
        fmt.Printf("Error creating product management smart contract: %s", err.Error())
        return
    }

    if err := chaincodeInstance.Start(); err != nil {
        fmt.Printf("Error starting product management smart contract: %s", err.Error())
    }
}

