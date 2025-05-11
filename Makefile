# Name of the final binary
BINARY_NAME=app

# Default target: generate, build, and run
all: generate build run

# Step 1: Generate Go code from .templ files
generate:
	@echo "Generating templ files..."
	templ generate

# Step 2: Build the Go project
build:
	@echo "Building the project..."
	go build -o $(BINARY_NAME) ./cmd/server

# Step 3: Run the built binary
run:
	@echo "Running the app..."
	./$(BINARY_NAME)

# Clean built binary
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)