FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.work go.work.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
RUN go mod download

WORKDIR /app/cmd/api
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /api .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /api /api
EXPOSE 8080
CMD ["/api"]
