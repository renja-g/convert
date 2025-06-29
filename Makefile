BINARY_NAME=convert
BUILD_DIR=bin

## build: builds the application binary.
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	@CGO_LDFLAGS="-Wl,-no_warn_duplicate_libraries" go build -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/convert
	@echo "${BINARY_NAME} built in ${BUILD_DIR}"

## clean: cleans up build artifacts.
clean:
	@echo "Cleaning up..."
	@rm -rf ${BUILD_DIR}
	@echo "Cleaned."

.PHONY: build clean 