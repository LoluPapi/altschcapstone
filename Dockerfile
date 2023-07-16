# Use the official Golang image as the base image
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download


# Copy the source code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Set the entry point command
CMD ["./main"]
