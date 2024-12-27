BINARY_NAME=FineMarkdown
SOURCE_DIR=./cmd/md_editor
FYNE_PACKAGE=github.com/fyne-io/fyne/v2/cmd/fyne

build:
	@echo "Building $(BINARY_NAME) with Fyne..."
	@rm -f ${SOURCE_DIR}/*.tar.xz 
	@cd ${SOURCE_DIR} && fyne package -os linux --release

run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(SOURCE_DIR)

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f ${SOURCE_DIR}/*.tar.xz 
	@echo "Cleaned!"

test:
	go test -v ./...

.PHONY: build run clean test

