start:
	go run cmd/main.go
migrate:
	goose -dir pkg/db/migrations postgres "host=${PSQL_HOST} port=${PSQL_PORT} user=${PSQL_USERNAME} password=${PSQL_PASSWORD} dbname=${PSQL_DBNAME} sslmode=disabled" up
rollback:
	goose -dir pkg/db/migrations postgres "host=${PSQL_HOST} port=${PSQL_PORT} user=${PSQL_USERNAME} password=${PSQL_PASSWORD} dbname=${PSQL_DBNAME} sslmode=disabled" down
test:
	go test -v ./...