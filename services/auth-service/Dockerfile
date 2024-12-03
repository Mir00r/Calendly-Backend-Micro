# Start with the Go image
FROM golang:1.20

# Set working directory
WORKDIR /app

# Copy files
COPY . .

# Build the Go application
RUN go build -o auth-service ./cmd/main.go

# Expose the service port
EXPOSE 8080

# Command to run the executable
CMD ["./auth-service"]
