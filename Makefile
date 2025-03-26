# Variables
APP_NAME=img-sizer
AWS_REGION?=eu-central-1
AWS_ACCOUNT_ID:=$(shell aws sts get-caller-identity --query Account --output text)
DOCKER_REGISTRY=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
VERSION=$(shell git describe --tags --always --dirty)
GOOS?=linux
GOARCH?=amd64

.PHONY: build dev docker-build docker-run docker-push create-ecr deploy clean

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(APP_NAME) ./cmd/img-sizer
	@echo "Build complete"

# Start development server with air
dev:
	@echo "Starting development server..."
	APP_ENV=local air

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) --platform $(GOOS)/$(GOARCH) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

# Run Docker container locally
docker-run:
	@echo "Starting Docker container..."
	docker run \
		--name $(APP_NAME) \
		--env-file .env \
		-p 8080:8080 \
		$(APP_NAME):$(VERSION)

# Push Docker image to AWS ECR
docker-push:
	@echo "Pushing to AWS ECR..."
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(DOCKER_REGISTRY)
	docker tag $(APP_NAME):$(VERSION) $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION)
	docker tag $(APP_NAME):latest $(DOCKER_REGISTRY)/$(APP_NAME):latest
	docker push $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(APP_NAME):latest

# Create ECR repository if it doesn't exist
create-ecr:
	@echo "Creating ECR repository..."
	aws ecr describe-repositories --repository-names $(APP_NAME) --no-cli-pager >/dev/null 2>&1 || \
	aws ecr create-repository --repository-name $(APP_NAME) --no-cli-pager

# Deploy to AWS (build, push, and create ECR if needed)
deploy: create-ecr docker-build docker-push

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/ tmp/
	docker rm -f $(APP_NAME) 2>/dev/null || true


# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  dev          - Start development server with air"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container locally"
	@echo "  docker-push  - Push Docker image to AWS ECR"
	@echo "  clean        - Clean up build artifacts"
	@echo "  create-ecr   - Create ECR repository"
	@echo "  deploy       - Deploy to AWS (build, push, create ECR)" 
