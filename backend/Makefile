.PHONY: migrate-up migrate-down migrate-create db-reset run

# Database migration commands
migrate-up:
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable" up

migrate-down:
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/database/migrations -seq $$name

# Reset database (down all migrations then up)
db-reset:
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable" down -all
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable" up

# Run the application
run:
	go run cmd/server/main.go
