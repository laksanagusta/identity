# ============================================
# Docker Configuration
# ============================================
DOCKER_USERNAME ?= laksanadika
IMAGE_NAME ?= identity
VERSION ?= latest
FULL_IMAGE_NAME = $(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION)

# ============================================
# Build Commands
# ============================================

## build: Build Go binary locally
build:
	@echo "📦 Building Go binary..."
	go build -o identity .

## docker-build: Build Docker image for linux/amd64
docker-build:
	@echo "🐳 Building Docker image: $(FULL_IMAGE_NAME) (platform: linux/amd64)"
	docker build --no-cache --platform linux/amd64 -t $(FULL_IMAGE_NAME) .
	docker tag $(FULL_IMAGE_NAME) $(DOCKER_USERNAME)/$(IMAGE_NAME):latest
	@echo "✅ Docker image built successfully!"

# ============================================
# Docker Hub Commands
# ============================================

## docker-login: Login to Docker Hub
docker-login:
	@echo "🔐 Logging in to Docker Hub..."
	docker login
	@echo "✅ Login successful!"

## docker-push: Push image to Docker Hub
docker-push:
	@echo "🚀 Pushing image to Docker Hub: $(FULL_IMAGE_NAME)"
	docker push $(FULL_IMAGE_NAME)
	docker push $(DOCKER_USERNAME)/$(IMAGE_NAME):latest
	@echo "✅ Image pushed successfully!"
	@echo ""
	@echo "📋 Pull command for your cloud server:"
	@echo "   docker pull $(FULL_IMAGE_NAME)"

## docker-release: Build and push image to Docker Hub (one command)
docker-release: docker-build docker-push
	@echo ""
	@echo "🎉 Release complete!"
	@echo "📋 To run on your cloud server:"
	@echo "   docker pull $(FULL_IMAGE_NAME)"
	@echo "   docker run -d --name identity -p 5001:5001 --env-file .env $(FULL_IMAGE_NAME)"

# ============================================
# Cloud Server Helper Commands
# ============================================

## docker-pull: Pull latest image from Docker Hub
docker-pull:
	@echo "📥 Pulling image from Docker Hub..."
	docker pull $(FULL_IMAGE_NAME)
	@echo "✅ Image pulled successfully!"

## docker-run: Run the container locally
docker-run:
	@echo "🚀 Running container..."
	docker run -d \
		--name identity \
		-p 5001:5001 \
		--env-file .env \
		--restart unless-stopped \
		$(FULL_IMAGE_NAME)
	@echo "✅ Container started!"

## docker-stop: Stop and remove container
docker-stop:
	@echo "🛑 Stopping container..."
	docker stop identity || true
	docker rm identity || true
	@echo "✅ Container stopped and removed!"

## docker-logs: View container logs
docker-logs:
	docker logs -f identity

## docker-restart: Restart the container
docker-restart: docker-stop docker-run

# ============================================
# Deployment Workflow (for cloud server)
# ============================================

## deploy: Pull latest image and restart container (run on cloud server)
deploy: docker-stop docker-pull docker-run
	@echo ""
	@echo "🎉 Deployment complete!"
	@echo "📋 Check logs with: make docker-logs"

# ============================================
# Cleaning Commands
# ============================================

## docker-clean: Remove local Docker images
docker-clean:
	@echo "🧹 Cleaning Docker images..."
	docker rmi $(FULL_IMAGE_NAME) || true
	docker rmi $(DOCKER_USERNAME)/$(IMAGE_NAME):latest || true
	@echo "✅ Cleaned!"

## clean: Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f identity
	@echo "✅ Cleaned!"

# ============================================
# Development Commands
# ============================================

## dev: Run development server
dev:
	go run main.go

## test: Run tests
test:
	go test -v ./...

# ============================================
# Help
# ============================================

## help: Show this help message
help:
	@echo "Identity Service - Docker Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
	@echo ""
	@echo "Configuration:"
	@echo "  DOCKER_USERNAME: $(DOCKER_USERNAME)"
	@echo "  IMAGE_NAME:      $(IMAGE_NAME)"
	@echo "  VERSION:         $(VERSION)"
	@echo "  FULL_IMAGE_NAME: $(FULL_IMAGE_NAME)"
	@echo ""
	@echo "Override example:"
	@echo "  make docker-release VERSION=v1.0.0"

.PHONY: build docker-build docker-login docker-push docker-release \
        docker-pull docker-run docker-stop docker-logs docker-restart \
        deploy docker-clean clean dev test help
