# Variables
DC = docker-compose

.PHONY: up down restart logs ps clean

# Start and build containers
up:
	$(DC) up -d --build

# Stop and remove containers
down:
	$(DC) down

# Restart
restart: down up

# Logs reading
logs:
	$(DC) logs -f

# Container status
ps:
	$(DC) ps

# Clean up unused Docker resources (containers, networks, images, and build cache)
clean:
	docker system prune -f