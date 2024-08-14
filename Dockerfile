# Stage 1
FROM golang:1.22-alpine AS builder
LABEL maintainer="savioruz <jakueenak@gmail.com>"

# Install dependencies
RUN apk --no-cache add ca-certificates

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o roastgithub-api .

# Stage 2
FROM scratch

# Copy CA certificates from the builder stage to enable SSL verification
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/roastgithub-api", "/build/.env", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/roastgithub-api"]
