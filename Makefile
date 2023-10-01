tests:
	go test ./... -coverprofile=coverage.out -cover
	go tool cover -html="coverage.out" -o "coverage.html"

docker:
	docker-compose down
	docker-compose up -d --build

doc:
	swag init --parseInternal -g cmd/api/main.go

.PHONY: docker tests doc