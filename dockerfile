FROM golang:1.24-bullseye

WORKDIR /app

# Install git for module downloads
RUN apt-get update && apt-get install -y git

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download **only remote dependencies**
RUN go mod download

# Copy the full project (local/internal packages included)
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/main.go

EXPOSE 8082
CMD ["./main"]
