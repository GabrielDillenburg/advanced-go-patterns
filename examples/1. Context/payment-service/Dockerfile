# payment-service/Dockerfile

# Start from the official Golang image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go app
RUN go build -o /payment-service

# Expose port 8082 to the outside world
EXPOSE 8082

# Command to run the executable
CMD ["/payment-service"]
