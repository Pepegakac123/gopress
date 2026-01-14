# Nazwa pliku wynikowego
BINARY_NAME=gopress
BUILD_DIR=bin
LDFLAGS=-ldflags "-s -w"

# --- KONFIGURACJA SDK ---
# Tutaj podaj PEŁNĄ ścieżkę, gdzie rozpakowałeś SDK.
# Jeśli zrobiłeś to w folderze projektu, użyj $(CURDIR)
MACOS_SDK_PATH=/home/Kacper/sdk/MacOSX12.3.sdk

# --- KONFIGURACJA ZIG ---
# Windows
ZIG_WIN_CC=zig cc -target x86_64-windows-gnu
ZIG_WIN_CXX=zig c++ -target x86_64-windows-gnu

# macOS
ZIG_MAC_CC=zig cc -target aarch64-macos -isysroot $(MACOS_SDK_PATH)
ZIG_MAC_CXX=zig c++ -target aarch64-macos -isysroot $(MACOS_SDK_PATH)

all: clean windows linux mac

windows:
	@echo "🔨 Budowanie dla Windows (x64)..."
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="$(ZIG_WIN_CC)" CXX="$(ZIG_WIN_CXX)" go build -ldflags "-s -w -extldflags '-static'" -o $(BUILD_DIR)/$(BINARY_NAME).exe ./cmd/gopress

linux:
	@echo "🐧 Budowanie dla Linux (x64)..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/gopress


mac:
	@echo "🍎 Budowanie dla macOS (Apple Silicon)..."
	@if [ ! -d "$(MACOS_SDK_PATH)" ]; then echo "❌ BŁĄD: Nie znaleziono SDK w $(MACOS_SDK_PATH)"; exit 1; fi
	
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 \
	CC="$(ZIG_MAC_CC)" \
	CXX="$(ZIG_MAC_CXX)" \
	CGO_CFLAGS="-isysroot $(MACOS_SDK_PATH) -I$(MACOS_SDK_PATH)/usr/include -Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wno-macro-redefined" \
	CGO_LDFLAGS="-isysroot $(MACOS_SDK_PATH) -L$(MACOS_SDK_PATH)/usr/lib -F$(MACOS_SDK_PATH)/System/Library/Frameworks" \
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac ./cmd/gopress

clean:
	@echo "🧹 Czyszczenie..."
	rm -rf $(BUILD_DIR)