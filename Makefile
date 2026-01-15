# Makefile for go-verbosity project

.PHONY: help build build-static test lint clean install examples deps

# Default target
help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Основные команды для библиотеки
build: ## Собрать библиотеку
	go build ./...

test: ## Запустить тесты
	go test -v ./...

lint: ## Запустить линтер
	golangci-lint run

vet: ## Запустить go vet
	go vet ./...

fmt: ## Форматировать код
	go fmt ./...

deps: ## Обновить зависимости
	go mod tidy
	go mod verify

clean: ## Очистить скомпилированные файлы
	find . -name "*.exe" -delete
	find . -name "info-bot" -delete
	find . -name "dist" -type d -exec rm -rf {} + 2>/dev/null || true

# Команды для примеров
examples: build-examples ## Собрать все примеры

build-examples: ## Собрать пример info-bot
	@echo "Сборка info-bot для текущей платформы..."
	cd examples/info-bot && go build -o info-bot .

build-static: ## Собрать статический бинарник для Linux x64
	@echo "Сборка статического бинарника для Linux x64..."
	cd examples/info-bot && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags="-s -w" -trimpath -o info-bot-linux-x64 .

build-multiplatform: ## Собрать бинарники для всех платформ (требует установки upx)
	@echo "Сборка для всех платформ..."
	mkdir -p dist
	platforms="linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64"; \
	for platform in $$platforms; do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		BINARY_NAME="info-bot"; \
		if [ "$$GOOS" = "windows" ]; then BINARY_NAME="$$BINARY_NAME.exe"; fi; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH \
		go build -ldflags="-s -w" -trimpath -o "dist/$$BINARY_NAME" \
		examples/info-bot/main.go; \
		cd dist && \
		tar -czf "info-bot-$$GOOS-$$GOARCH.tar.gz" "$$BINARY_NAME" && \
		if [ "$$GOOS" = "windows" ]; then zip -q "info-bot-$$GOOS-$$GOARCH.zip" "$$BINARY_NAME"; fi && \
		cd ..; \
	done
	@echo "Готово! Файлы в папке dist/:"
	@ls -lh dist/

# Установка
install: ## Установить библиотеку
	go install ./...

install-examples: ## Установить примеры
	go install github.com/ivmaks/go-verbosity/examples/info-bot@latest

# Полная проверка (как в CI)
ci: fmt vet test lint ## Запустить все проверки как в CI

# Установка дополнительных инструментов
install-tools: ## Установить дополнительные инструменты разработки
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Установка UPX для сжатия бинарников..."
	@if command -v apt-get >/dev/null 2>&1; then \
		sudo apt-get update && sudo apt-get install -y upx-ucl; \
	elif command -v brew >/dev/null 2>&1; then \
		brew install upx; \
	else \
		echo "Установите UPX вручную: https://upx.github.io/"; \
	fi

# Отчеты
report: ## Показать отчет о покрытии кода
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчет сохранен в coverage.html"

# Генерация документации
docs: ## Сгенерировать документацию
	go doc -all ./... > docs.txt
	godocdown ./... > README-GODOC.md

# Проверка безопасности
security: ## Проверить безопасность кода
	gosec ./...

# Все команды сборки и проверки
all: clean build test lint build-static ## Полная сборка и проверка