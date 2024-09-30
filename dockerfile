# Use the official Golang image as the base image
FROM --platform=linux/amd64 golang:1.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Generate Swagger documentation (if applicable)
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]