FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.work go.work.sum ./
RUN go mod download

COPY cmd/api ./cmd/api
COPY internal ./internal
COPY pkg ./pkg
WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 go build -o /api .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /api /api
EXPOSE $PORT
CMD ["/api"]
