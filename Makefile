db-up:
	docker run --name db-tech-school -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secretpassword -d postgres:12-alpine

db-down:
	docker container rm -f db-tech-school

db-start:
	docker start db-tech-school

db-stop:
	docker stop db-tech-school

db-logs:
	docker logs -f db-tech-school

migration:
	migrate create -ext sql -dir db/migrations -seq init_schema

migrateup:
	migrate -path db/migrations -database "postgresql://root:secretpassword@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secretpassword@localhost:5432/simple_bank?sslmode=disable" -verbose down

db-create-database:
	docker exec -it db-tech-school createdb --username=root --owner=root simple_bank

db-drop:
	docker exec -it db-tech-school dropdb simple_bank

db-psql:
	docker exec -it db-tech-school psql -U root simple_bank

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

test-wo-cache-1:
	go test -v -cover -count=1 ./...

test-wo-cache-2:
	go clean -testcache && go test -v -cover ./...
