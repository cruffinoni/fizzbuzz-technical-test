FROM golang:1.20

WORKDIR /app

COPY . .

# Generate the documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseInternal -g cmd/api/main.go

RUN go mod vendor
RUN go build -o main cmd/api/main.go

EXPOSE 8080

CMD ["./main"]
