// payment-service/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handlePayment(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // Extract context from the HTTP request
    log.Println("Received payment request")

    orderID := r.URL.Query().Get("orderID")
    log.Println("Processing payment for order:", orderID)

    select {
    case <-time.After(3 * time.Second): // Simulate payment processing delay
        log.Println("Payment processed for order:", orderID)
        fmt.Fprintf(w, "Payment for order %s completed", orderID)
    case <-ctx.Done(): // Check if context is canceled
        log.Println("Payment processing canceled for order:", orderID)
        http.Error(w, "Payment canceled", http.StatusRequestTimeout)
    }
}

func main() {
    http.HandleFunc("/pay", handlePayment)
    log.Println("Payment Service running on :8082")
    log.Fatal(http.ListenAndServe(":8082", nil))
}
