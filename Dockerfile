# Use the official Golang image as the base image
FROM golang:1.22-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git (required for go modules) and other dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/scattergories

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main /app/main

# Copy the .env file into the container
COPY .env .env

# Copy the prompts file into the container
COPY config/prompts.txt config/prompts.txt

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/main"]
