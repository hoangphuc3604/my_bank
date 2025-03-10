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
	/home/hoangphuc3604/go/bin/migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	/home/hoangphuc3604/go/bin/migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	/home/hoangphuc3604/go/bin/migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	/home/hoangphuc3604/go/bin/migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

server:
	go run main.go

.PHONY: sqlc test createdb dropdb migrateup migratedown server