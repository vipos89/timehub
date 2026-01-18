.PHONY: build run up down clean swag-gen

SERVICES := api-gateway auth-service company-service booking-service crm-service report-service

build:
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		go build -o bin/$$service ./services/$$service/cmd; \
	done

run: build
	@echo "Starting all services locally..."
	@# This is a simple run, in reality you probably want separate terminals or docker-compose
	@echo "Use 'make up' to run with Docker Compose"

up:
	@echo "Starting Docker Compose..."
	docker-compose up -d --build

down:
	@echo "Stopping Docker Compose..."
	docker-compose down

clean:
	@echo "Cleaning binaries..."
	rm -rf bin/

swag-gen:
	@echo "Generating Swagger docs..."
	@for service in $(SERVICES); do \
		echo "Generating docs for $$service..."; \
		if [ -d "services/$$service/cmd" ]; then \
			cd services/$$service/cmd; \
			SEARCH_DIRS="."; \
			if [ -d "../internal" ]; then \
				SEARCH_DIRS=".,../internal"; \
			fi; \
			go run github.com/swaggo/swag/cmd/swag init -g main.go -d $$SEARCH_DIRS -o ../docs --parseDependency --parseInternal; \
			cd ../../../; \
		fi \
	done

tidy:
	@echo "Tidying modules..."
	@for service in $(SERVICES); do \
		cd services/$$service && go mod tidy && cd ../../; \
	done
	cd pkg && go mod tidy
