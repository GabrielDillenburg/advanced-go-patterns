package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type OrderService struct {
    inventoryServiceURL string
    paymentServiceURL   string
}

func (os *OrderService) HandleOrder(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    productID := "1234"
    log.Println("Placing order for product:", productID)

    // Call inventory service
    if err := os.checkInventory(ctx, productID); err != nil {
        http.Error(w, "Inventory check failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Call payment service
    if err := os.processPayment(ctx, "orderID-5678"); err != nil {
        http.Error(w, "Payment processing failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "Order successfully processed")
}

func (os *OrderService) checkInventory(ctx context.Context, productID string) error {
    req, err := http.NewRequestWithContext(ctx, "GET", os.inventoryServiceURL+"/check?productID="+productID, nil)
    if err != nil {
        return err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("inventory check failed, status: %d", resp.StatusCode)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("Inventory response:", string(body))
    return nil
}

func (os *OrderService) processPayment(ctx context.Context, orderID string) error {
    req, err := http.NewRequestWithContext(ctx, "POST", os.paymentServiceURL+"/pay?orderID="+orderID, nil)
    if err != nil {
        return err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("payment processing failed, status: %d", resp.StatusCode)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("Payment response:", string(body))
    return nil
}

func main() {
    inventoryServiceURL := os.Getenv("INVENTORY_SERVICE_URL")
    paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")

    orderService := &OrderService{
        inventoryServiceURL: inventoryServiceURL,
        paymentServiceURL:   paymentServiceURL,
    }

    http.HandleFunc("/order", orderService.HandleOrder)

    log.Println("Order Service running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
