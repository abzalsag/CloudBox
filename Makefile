DB_URL ?= postgres://cloudbox:cloudbox@localhost:5434/cloudbox?sslmode=disable

migrate-up:
	migrate -path database/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path database/migrations -database "$(DB_URL)" down 1

migrate-new:
	migrate create -ext sql -dir database/migrations -seq $(name)

sqlc:
	sqlc generate