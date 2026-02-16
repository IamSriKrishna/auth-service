FROM golang:1.23-alpine AS builder

# Build arguments
ARG VERSION=dev
ARG BUILDTIME=unknown
ARG GITHUB_SHA=unknown
ARG GITHUB_REF=unknown

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with build info
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILDTIME} -X main.gitCommit=${GITHUB_SHA} -X main.gitRef=${GITHUB_REF}" \
    -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 3001

# Run the binary
CMD ["./main"]
