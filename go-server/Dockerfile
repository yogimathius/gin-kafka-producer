# Use the official Golang image as a build stage
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Install any required dependencies (optional, depends on your app)
RUN apk add --no-cache ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /main

# Make sure the binary is executable
RUN chmod +x /main

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/main"]
