REDIS_ADDR ?= 'postgres://ps_user:ps_password@localhost:7002/backend?sslmode=disable'

# Test
.PHONY: test
test:
	docker-compose up -d
	go test ./... -v -count=1
	docker-compose down

# Run Commands
.PHONY: run-redis
run-redis:
	docker-compose up

.PHONY: run-api
run-api:
	go run app/api/*.go

.PHONY: run-matchmaker
run-matchmaker:
	go run app/matchmaking_worker/*.go

.PHONY: run-cleaner
run-cleaner:
	go run app/tickets_cleaner_worker/*.go