.PHONY: test up stop down

test:
	go test -v ./...

up:
	@echo "Starting services..."
	docker-compose up -d --build

stop:
	@echo "Stopping services..."
	docker-compose stop

down:
	@echo "Downing services..."
	docker-compose down

