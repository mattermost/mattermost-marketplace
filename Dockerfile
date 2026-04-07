# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG BUILD_UPSTREAM_URL=""

RUN BUILD_TAG=$(git describe --abbrev=0 --tags 2>/dev/null || echo "dev") && \
    BUILD_HASH=$(git rev-parse HEAD 2>/dev/null || echo "unknown") && \
    BUILD_HASH_SHORT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
      -X 'github.com/mattermost/mattermost-marketplace/internal/api.buildTag=${BUILD_TAG}' \
      -X 'github.com/mattermost/mattermost-marketplace/internal/api.buildHash=${BUILD_HASH}' \
      -X 'github.com/mattermost/mattermost-marketplace/internal/api.buildHashShort=${BUILD_HASH_SHORT}' \
      -X 'main.upstreamURL=${BUILD_UPSTREAM_URL}'" \
    -o /marketplace ./cmd/marketplace/

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

RUN adduser -D -h /app marketplace
USER marketplace
WORKDIR /app

COPY --from=builder /marketplace /app/marketplace
COPY plugins.json /app/plugins.json

EXPOSE 8085

ENTRYPOINT ["/app/marketplace"]
CMD ["server", "--database", "/app/plugins.json", "--listen", ":8085"]
