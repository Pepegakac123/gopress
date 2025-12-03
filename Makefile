# Nazwa pliku wynikowego
BINARY_NAME=gopress

# Katalog wyjściowy
BUILD_DIR=bin

# Flagi kompilatora
LDFLAGS=-ldflags "-s -w"

all: clean windows linux mac

windows:
	@echo "🔨 Budowanie dla Windows (x64)..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe ./cmd/gopress

linux:
	@echo "🐧 Budowanie dla Linux (x64)..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/gopress

mac:
	@echo "🍎 Budowanie dla macOS (Apple Silicon + Intel)..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac-m1 ./cmd/gopress
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac-intel ./cmd/gopress

clean:
	@echo "🧹 Czyszczenie..."
	rm -rf $(BUILD_DIR)

run:
	go run ./cmd/gopress
