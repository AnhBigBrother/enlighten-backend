FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api 

EXPOSE 8080

CMD ["bin/api"]