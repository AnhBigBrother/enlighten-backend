# BUILDER -------------------------------
FROM golang:latest AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api 


# RUNNER --------------------------------
FROM alpine:latest AS runner
WORKDIR /app

COPY --from=builder /app/bin/api ./
COPY --from=builder /app/.env.production ./

EXPOSE 8080

ENTRYPOINT ["./api", "-env=production"]