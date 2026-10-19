BIN   := bin/catalog
IMAGE := catalog-poc:local
PORT  ?= 8080

.PHONY: run build test vet check docker-build docker-run clean

run:
	PORT=$(PORT) go run ./cmd/server

build:
	CGO_ENABLED=0 go build -trimpath -o $(BIN) ./cmd/server

test:
	go test ./...

vet:
	go vet ./...

check: vet test

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm -p $(PORT):8080 $(IMAGE)

clean:
	rm -rf bin
