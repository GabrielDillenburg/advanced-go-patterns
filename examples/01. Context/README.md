
# Go Microservices with Docker Compose

This project demonstrates a microservice architecture using Go and Docker Compose. The architecture includes three services:

1. **Order Service**: Handles user requests for placing an order and communicates with the Inventory and Payment services.
2. **Inventory Service**: Checks inventory availability for products.
3. **Payment Service**: Processes payments for orders.

All services are built using Go, and cancellation and timeout propagation are managed using `context.Context`. The services are orchestrated using Docker Compose.

## Project Structure

```
├── order-service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── inventory-service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── payment-service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
└── docker-compose.yml
```

## Services

### 1. Order Service

The **Order Service** handles user requests to place an order. It communicates with the **Inventory Service** to check product availability and the **Payment Service** to process the order's payment.

#### Endpoint:
- `GET /order`: Processes an order by checking inventory and processing payment.

### 2. Inventory Service

The **Inventory Service** simulates an inventory check for a product. It accepts a `productID` parameter and returns whether the product is available.

#### Endpoint:
- `GET /check?productID={productID}`: Checks inventory for the specified product.

### 3. Payment Service

The **Payment Service** simulates payment processing for an order. It accepts an `orderID` parameter and processes the payment for that order.

#### Endpoint:
- `POST /pay?orderID={orderID}`: Processes payment for the specified order.

## Cancellation Handling

The project demonstrates the use of `context.Context` for managing request cancellation and timeout propagation across multiple services. If the **Order Service** cancels a request (due to user cancellation or timeout), the cancellation signal is propagated to the **Inventory Service** and **Payment Service**, stopping their respective operations.

## Setup

### Prerequisites

- Docker and Docker Compose installed on your machine.
- Go 1.20+ (if you want to run the services locally outside of Docker).

### Running the Application

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/go-microservices-example.git
   cd go-microservices-example
   ```

2. **Build and run the services using Docker Compose**:

   ```bash
   docker-compose build
   docker-compose up
   ```

   This command will build the Docker images and start the services:
   - Order Service: `http://localhost:8080`
   - Inventory Service: `http://localhost:8081`
   - Payment Service: `http://localhost:8082`

### Testing the Services

1. **Place an order**:
   You can place an order by making a request to the **Order Service**:

   ```bash
   curl http://localhost:8080/order
   ```

   The **Order Service** will call the **Inventory Service** to check if the product is available and then call the **Payment Service** to process the payment.

2. **Simulate a cancellation**:
   While the order request is processing, you can cancel the request (e.g., by pressing `Ctrl+C` in the terminal). The cancellation will propagate to the **Inventory Service** and **Payment Service**, and their operations will stop.

### Logs

You can check the logs of each service in real time using Docker Compose:

```bash
docker-compose logs -f
```

This will display logs for all services, showing their interaction and any cancellation events.

### Stopping the Services

To stop the services, run:

```bash
docker-compose down
```

This will stop and remove the containers.

## Code Overview

### Order Service (main.go)

The **Order Service** communicates with both the **Inventory Service** and **Payment Service** using HTTP. It propagates the `context.Context` from the user request to downstream services.

```go
func (os *OrderService) HandleOrder(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    productID := "1234"
    log.Println("Placing order for product:", productID)

    if err := os.checkInventory(ctx, productID); err != nil {
        http.Error(w, "Inventory check failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if err := os.processPayment(ctx, "orderID-5678"); err != nil {
        http.Error(w, "Payment processing failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "Order successfully processed")
}
```

### Inventory Service (main.go)

The **Inventory Service** simulates a delay when checking inventory. It listens for cancellation via `ctx.Done()`.

```go
select {
case <-time.After(2 * time.Second):
    log.Println("Inventory check completed for product:", productID)
    fmt.Fprintf(w, "Product %s is available", productID)
case <-ctx.Done():
    log.Println("Inventory check canceled for product:", productID)
    http.Error(w, "Request canceled", http.StatusRequestTimeout)
}
```

### Payment Service (main.go)

The **Payment Service** processes payments and listens for cancellation in the same way as the **Inventory Service**.

```go
select {
case <-time.After(3 * time.Second):
    log.Println("Payment processed for order:", orderID)
    fmt.Fprintf(w, "Payment for order %s completed", orderID)
case <-ctx.Done():
    log.Println("Payment processing canceled for order:", orderID)
    http.Error(w, "Payment canceled", http.StatusRequestTimeout)
}
```

