BINARY_NAME=kmc
PROJECT_NAME=kmc-palette
ARCH=amd64
ARGS=./test-image.png
.DEFAULT_GOAL := run

build:
	GOARCH=$(ARCH) GOOS=linux go build -o ./bin/$(BINARY_NAME)-linux-$(ARCH) ./cmd/$(PROJECT_NAME)/main.go
	GOARCH=$(ARCH) GOOS=darwin go build -o ./bin/$(BINARY_NAME)-darwin-$(ARCH) ./cmd/$(PROJECT_NAME)/main.go
	GOARCH=$(ARCH) GOOS=windows go build -o ./bin/$(BINARY_NAME)-windows-$(ARCH).exe ./cmd/$(PROJECT_NAME)/main.go

run-linux: build
	./bin/$(BINARY_NAME)-linux-$(ARCH)

run-darwin: build
	./bin/$(BINARY_NAME)-darwin-$(ARCH)

run-windows: build
	./bin/$(BINARY_NAME)-windows-$(ARCH).exe

run:
	go run ./cmd/$(PROJECT_NAME)/main.go $(ARGS)

clean:
	go clean
	rm ./bin/$(BINARY_NAME)-linux-$(ARCH)
	rm ./bin/$(BINARY_NAME)-darwin-$(ARCH)
	rm ./bin/$(BINARY_NAME)-windows-$(ARCH).exe

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

vet:
	go vet

lint:
	golangci-lint run --enable-all
