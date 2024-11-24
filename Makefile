demo:
	docker compose up
	open http://localhost:3000
	open http://localhost:8080

test:
	go test ./...

api-docs:
	swag init -g cmd/server/main.go -o api