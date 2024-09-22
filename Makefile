BINARY_NAME=kmc-palette
PROJECT_NAME=kmc-palette
VERSION=0.1.0
ARCH=amd64
ARGS=./test-image.png
.DEFAULT_GOAL := run

build:
	GOARCH=$(ARCH) GOOS=linux go build -o ./bin/$(BINARY_NAME)-$(VERSION)-linux-$(ARCH) ./cmd/$(PROJECT_NAME)/main.go
	GOARCH=$(ARCH) GOOS=darwin go build -o ./bin/$(BINARY_NAME)-$(VERSION)-darwin-$(ARCH) ./cmd/$(PROJECT_NAME)/main.go
	GOARCH=$(ARCH) GOOS=windows go build -o ./bin/$(BINARY_NAME)-$(VERSION)-windows-$(ARCH).exe ./cmd/$(PROJECT_NAME)/main.go

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
	rm ./bin/$(BINARY_NAME)-$(VERSION)-linux-$(ARCH)
	rm ./bin/$(BINARY_NAME)-$(VERSION)-darwin-$(ARCH)
	rm ./bin/$(BINARY_NAME)-$(VERSION)-windows-$(ARCH).exe
	rm ./bin/$(BINARY_NAME)-$(VERSION)-sha256-checksums.txt

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

vet:
	go vet

lint:
	golangci-lint run --enable-all

checksum: build
	@echo "Generating checksums..."
	@sha256sum ./bin/$(BINARY_NAME)-$(VERSION)-linux-$(ARCH) | awk '{print $$1, $$2}' | awk '{print $$1, substr($$2, index($$2, "/")+5)}' > ./bin/$(BINARY_NAME)-$(VERSION)-sha256-checksums.txt
	@sha256sum ./bin/$(BINARY_NAME)-$(VERSION)-darwin-$(ARCH) | awk '{print $$1, $$2}' | awk '{print $$1, substr($$2, index($$2, "/")+5)}' >> ./bin/$(BINARY_NAME)-$(VERSION)-sha256-checksums.txt
	@sha256sum ./bin/$(BINARY_NAME)-$(VERSION)-windows-$(ARCH).exe | awk '{print $$1, $$2}' | awk '{print $$1, substr($$2, index($$2, "/")+5)}' >> ./bin/$(BINARY_NAME)-$(VERSION)-sha256-checksums.txt
	@echo "Checksums written to ./bin/$(BINARY_NAME)-$(VERSION)-sha256-checksums.txt"
