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
		if [ -d "./services/$$service/cmd" ]; then \
			echo "Generating docs for $$service..."; \
			cd services/$$service/cmd && go run github.com/swaggo/swag/cmd/swag init -g main.go -o ../docs --parseDependency --parseInternal; \
			cd ../../../; \
		fi \
	done

tidy:
	@echo "Tidying modules..."
	@for service in $(SERVICES); do \
		cd services/$$service && go mod tidy && cd ../../; \
	done
	cd pkg && go mod tidy
