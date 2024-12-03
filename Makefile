# Variables
DOCKER_COMPOSE=docker-compose
SERVICES=auth-service event-service notification-service

.PHONY: all build start stop restart clean

# Build all services
build:
	$(DOCKER_COMPOSE) build

# Start services
start:
	$(DOCKER_COMPOSE) up -d

# Stop services
stop:
	$(DOCKER_COMPOSE) down

# Restart services
restart: stop start

# Clean Docker environment
clean:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	docker system prune -f
