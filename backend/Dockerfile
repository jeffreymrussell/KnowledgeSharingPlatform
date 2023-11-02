# Start from the latest golang base image
FROM golang:1.20.1 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o main .

######## Start a new stage from alpine #######
FROM alpine:latest

# Add the packages required for sqlite3 and TLS/SSL (if needed)
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

COPY --from=builder /app/config.json /root/config.json
COPY --from=builder /app/main.sqlite /root/main.sqlite

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
