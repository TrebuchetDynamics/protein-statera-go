APP=protein
UI=protein-ui

.PHONY: test build run ui vet clean

test:
	CGO_ENABLED=0 go test ./...

vet:
	CGO_ENABLED=0 go vet ./...

build:
	CGO_ENABLED=0 go build -o bin/$(APP) ./cmd/protein
	CGO_ENABLED=0 go build -o bin/$(UI) ./cmd/protein-ui

run:
	CGO_ENABLED=0 go run ./cmd/protein

ui:
	CGO_ENABLED=0 go run ./cmd/protein-ui

clean:
	go clean
