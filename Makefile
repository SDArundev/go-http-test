dev:
	docker-compose -p golang-test up

sql: 
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate