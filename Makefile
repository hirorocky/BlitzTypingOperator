# BlitzTypingOperator Makefile

.PHONY: build build-release clean test run lint help

# デフォルトターゲット
all: build

# 開発用ビルド
build:
	@echo "Building BlitzTypingOperator..."
	go build -o BlitzTypingOperator ./cmd/BlitzTypingOperator

# リリース用ビルド（最適化あり）
build-release:
	@echo "Building BlitzTypingOperator (release)..."
	go build -ldflags="-s -w" -o BlitzTypingOperator ./cmd/BlitzTypingOperator

# テスト実行
test:
	@echo "Running tests..."
	go test ./...

# テスト実行（詳細出力）
test-verbose:
	@echo "Running tests (verbose)..."
	go test -v ./...

# Lint実行
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# アプリケーション実行（埋め込みデータ使用）
run: build
	@echo "Running BlitzTypingOperator..."
	./BlitzTypingOperator

# クリーンアップ
clean:
	@echo "Cleaning up..."
	@rm -f BlitzTypingOperator

# ヘルプ
help:
	@echo "BlitzTypingOperator Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-release  - Build with optimizations"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  lint           - Run golangci-lint"
	@echo "  run            - Build and run with embedded data"
	@echo "  clean          - Remove build artifacts"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "To use external data:"
	@echo "  ./BlitzTypingOperator -data /path/to/custom_data"
