// inventory-service/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handleInventoryCheck(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // Extract context from the HTTP request
    log.Println("Received inventory check request")

    productID := r.URL.Query().Get("productID")
    log.Println("Checking inventory for product:", productID)

    select {
    case <-time.After(2 * time.Second): // Simulate inventory check delay
        log.Println("Inventory check completed for product:", productID)
        fmt.Fprintf(w, "Product %s is available", productID)
    case <-ctx.Done(): // Check if context is canceled
        log.Println("Inventory check canceled for product:", productID)
        http.Error(w, "Request canceled", http.StatusRequestTimeout)
    }
}

func main() {
    http.HandleFunc("/check", handleInventoryCheck)
    log.Println("Inventory Service running on :8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}
