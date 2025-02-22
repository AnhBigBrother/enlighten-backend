gen-proto:
	cd proto; sh gen.sh; cd ..;

build-img:
	docker build -t anhbruhh/enlighten-api:0.2 .

goose-up: 
	goose postgres postgres://postgres:abc123@localhost:5432/enlighten up -dir sql/schema

goose-down: 
	goose postgres postgres://postgres:abc123@localhost:5432/enlighten down -dir sql/schema

start-dev: 
	go build -o bin/api cmd/api/main.go && ./bin/api

start-pro: 
	go build -o bin/api cmd/api/main.go && ./bin/api -env=production

.PHONY: start-dev start-pro gen-proto build-img goose-up goose-down