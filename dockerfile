FROM golang:latest AS base

# BUILDER -------------------------------
FROM base AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api 


# RUNNER --------------------------------
FROM base AS runner
WORKDIR /app

COPY --from=builder /app/bin/api ./
COPY --from=builder /app/.env.production ./
EXPOSE 8080

CMD ["./api"]