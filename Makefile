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
	/home/usr/go/bin/migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	/home/usr/go/bin/migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	/home/usr/go/bin/migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	/home/usr/go/bin/migrate -path db/migration -database "$(DB_URL)" -verbose down 1

server:
	go run main.go

.PHONY: sqlc test createdb dropdb migrateup migratedown server