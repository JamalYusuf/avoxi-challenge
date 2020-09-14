# Build stage
# Start from the version of Go that it was developed on
FROM golang:1.15-alpine AS builder

# Add Maintainer info
LABEL maintainer="Jamal Yusuf <contact@JamalYusuf.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose https and grpc ports
EXPOSE 10000 11000


# Command to run the executable
CMD ["./main"]