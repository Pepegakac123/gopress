
BINARY_NAME=gopress
BUILD_DIR=bin
LDFLAGS=-ldflags "-s -w"

# Ziga jako kompilator dla Windowsa
ZIG_CC=zig cc -target x86_64-windows-gnu
ZIG_CXX=zig c++ -target x86_64-windows-gnu

all: clean windows linux mac

windows:
	@echo "🔨 Budowanie dla Windows (x64) przy użyciu Zig (HEIC enabled)..."
	# CGO_ENABLED=1 jest kluczowe dla goheif!
	# Przekazujemy CC i CXX do Ziga, żeby obsłużył kod C/C++
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="$(ZIG_CC)" CXX="$(ZIG_CXX)" go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe ./cmd/gopress

linux:
	@echo "🐧 Budowanie dla Linux (x64)..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/gopress

mac:
	@echo "🍎 Budowanie dla macOS (Apple Silicon)..."
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="zig cc -target aarch64-macos" CXX="zig c++ -target aarch64-macos" go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac-m1 ./cmd/gopress

clean:
	@echo "🧹 Czyszczenie..."
	rm -rf $(BUILD_DIR)