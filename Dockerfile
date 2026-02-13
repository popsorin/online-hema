FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go mod tidy
RUN go build -o /app/hema-api ./cmd/api

FROM alpine:3.20

RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /app/hema-api /app/hema-api

USER appuser
EXPOSE 8080
CMD ["/app/hema-api"]
