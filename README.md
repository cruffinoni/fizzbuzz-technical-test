# fizzbuzz

## Introduction

The `fizzbuzz` project is a modern twist on the classic fizz-buzz game. Traditionally, the game replaces all numbers from 1 to 100 where multiples of 3 are substituted by "fizz", multiples of 5 by "buzz", and multiples of 15 by "fizzbuzz". The output would then look like: "1,2,fizz,4,buzz,fizz,7,...". This project offers a web server that exposes a REST API endpoint allowing users to customize the fizz-buzz rule.

## Getting Started

1. **Clone the repository**:
    ```bash
    git clone https://github.com/cruffinoni/fizzbuzz-technical-test.git
    ```
   
2. **Docker Setup**: This project uses docker-compose. Ensure that Docker is installed on your machine. Navigate to the project directory and run:
    ```bash
    docker-compose up
    ```

3. **Testing the API**: It's recommended to test the API with the "ping" route to ensure everything is running smoothly.
    ```bash
   curl http://localhost:8080/ping
   ```

## API Endpoints

### 1. Play FizzBuzz

- **Endpoint**: `/play`
- **Method**: `POST`
- **Body**:
```json
{
	"number1": <int>,
	"number2": <int>,
	"replace1": <string>,
	"replace2": <string>,
	"limit": <int>
}
```
- **Response**:
```json
{
	"result": <string>
}
```

### 2. Get Most Used Request

- **Endpoint**: `/most-used`
- **Method**: `GET`
- **Response**:
```json
{
	"int1": <int64>,
	"int2": <int64>,
	"hints": <int64>
}
```

### 3. Ping

- **Endpoint**: `/ping`
- **Method**: `GET`
- **Response**:
```json
{
	"message": "pong"
}
```

---

## Remark
Is it possible for the initial setup to take some time (about 3-5 mins). 
MySQL 8.0 needs to be downloaded and initialized. The database schema is then created and the project is compiled.