DB_URL=postgresql://root:root@localhost:5432/my_bank?sslmode=disable

sqlc:
	sqlc generate
test: 
	go test -v -cover ./...
createdb:
	docker exec -it postgres createdb --username=root --owner=root my_bank

dropdb:
	docker exec -it postgres dropdb my_bank
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down