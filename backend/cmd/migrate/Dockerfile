# Build from backend/: docker build -f cmd/migrate/Dockerfile . -t zb-migrate:local
# Standalone build (no go.work): migrate is its own module in cmd/migrate.
FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY cmd/migrate/ ./

RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /migrate .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /migrate /migrate
ENTRYPOINT ["/migrate"]
CMD ["up"]
