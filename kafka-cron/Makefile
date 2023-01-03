SHELL=/bin/bash
.PHONY: run test develop

up:  ## Start all or c=<name> containers in foreground
	docker-compose -f docker-compose.yml up
start: ## Start all or c=<name> containers in background
	docker-compose -f docker-compose.yml start $(c)
down: ## Clean all data
	docker-compose -f docker-compose.yml down $(c)
stop: ## Show status of containers
	docker-compose -f docker-compose.yml stop $(c)
restart: ## Restart all or c=<name> containers
	docker-compose -f docker-compose.yml stop $(c)
	docker-compose -f docker-compose.yml up -d $(c)
destroy: ## Clean all containers, images, volumes, networks, and undefined containers created with docker-compose.yml
	docker-compose -f docker-compose.yml down --rmi all -v --remove-orphans $(c)
up-more:  ## Start all or c=<name> containers in foreground
	docker-compose -f docker-compose-add-consumers.yml up $(c)
start-more: ## Start all or c=<name> containers in background
	docker-compose -f docker-compose-add-consumers.yml start $(c)
down-more: ## Clean all data
	docker-compose -f docker-compose-add-consumers.yml down $(c)
stop-more: ## Show status of containers
	docker-compose -f docker-compose-add-consumers.yml stop $(c)
restart-more: ## Restart all or c=<name> containers
	docker-compose -f docker-compose-add-consumers.yml stop $(c)
	docker-compose -f docker-compose-add-consumers.yml up -d $(c)
destroy-more: ## Clean all containers, images, volumes, networks, and undefined containers created with docker-compose-add-consumers.yml
	docker-compose -f docker-compose-add-consumers.yml down --rmi all -v --remove-orphans $(c)
up-ui:  ## Start all or c=<name> containers in foreground
	docker-compose -f docker-compose-ui.yml up $(c)
start-ui: ## Start all or c=<name> containers in background
	docker-compose -f docker-compose-ui.yml start $(c)
down-ui: ## Clean all data
	docker-compose -f docker-compose-ui.yml down $(c)
stop-ui: ## Show status of containers
	docker-compose -f docker-compose-ui.yml stop $(c)
restart-ui: ## Restart all or c=<name> containers
	docker-compose -f docker-compose-ui.yml stop $(c)
	docker-compose -f docker-compose-ui.yml up -d $(c)
destroy-ui: ## Clean all containers, images, volumes, networks, and undefined containers created with docker-compose-ui.yml
	docker-compose -f docker-compose-ui.yml down --rmi all -v --remove-orphans $(c)
up-metrics:  ## Start all or c=<name> containers in foreground
	docker-compose -f docker-compose-metrics.yml up
start-metrics: ## Start all or c=<name> containers in background
	docker-compose -f docker-compose-metrics.yml start $(c)
down-metrics: ## Clean all data
	docker-compose -f docker-compose-metrics.yml down $(c)
stop-metrics: ## Show status of containers
	docker-compose -f docker-compose-metrics.yml stop $(c)
restart-metrics: ## Restart all or c=<name> containers
	docker-compose -f docker-compose-metrics.yml stop $(c)
destroy-metrics: ## Clean all containers, images, volumes, networks, and undefined containers created with docker-compose.yml
	docker-compose -f docker-compose-metrics.yml down --rmi all -v --remove-orphans $(c)