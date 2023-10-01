tests:
	go mod vendor
	go test ./... -coverprofile=coverage.out -cover
	go tool cover -html="coverage.out" -o "coverage.html"

docker:
	docker-compose down
	docker-compose up -d --build
	sleep 3
	docker logs --follow "$(docker ps --format 'table {{.ID}} {{.Image}}' | grep 'fizzbuzz-api' | awk '{print $1}')"

doc:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseInternal -g cmd/api/main.go

clean:
	rm -rf coverage.out coverage.html vendor/

.PHONY: docker tests doc clean