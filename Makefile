run-db:
	@echo "Running database migration..."
	@docker compose -f 'docker-compose.yaml' up -d --build 'db'

init-test-db:
	@echo "Initializing database..."
	@docker compose -f 'docker-compose.yaml' down 'db-test'
	@docker compose -f 'docker-compose.yaml' up -d --build 'db-test'

test:
	@echo "Running tests..."
	go test -v ./tests/...

clean-test:
	@go clean -testcache
	go test -v ./tests/...