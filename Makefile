DB_URL ?= postgres://cloudbox:cloudbox@localhost:5432/cloudbox?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1

migrate-new:
	migrate create -ext sql -dir db/migrations -seq $(name)

sqlc:
	sqlc generate