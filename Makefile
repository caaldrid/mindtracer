.PHONY: run test db-up db-down migrate

run:
	cd backend && go run main.go

test:
	cd backend && go test ./...

db-up:
	cd backend && docker compose up -d

db-down:
	cd backend && docker compose down

migrate:
	cd backend && go run main.go -migrate
