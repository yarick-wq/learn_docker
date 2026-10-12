include .env
export

.PHONY: run-http-app migrate-up migrate-down run-local postgres-up

run-http-app:
	docker run -d -p 8080:8080 http_docker

migrate-up:
	@migrate -path migrations -database "${CONN_STRING}" up

migrate-down:
	@migrate -path migrations -database "${CONN_STRING}" down 1

run-local:
	go run .

postgres-up:
	@docker run -d --rm --name sql_docker -p 5432:5432 \
	-e POSTGRES_PASSWORD=${CONN_PASSWORD} -e POSTGRES_DB=${POSTGRES_DB} \
	-v employees_pgdata:/var/lib/postgresql postgres:18-bookworm
