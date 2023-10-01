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
    make docker
    ```

3. **Testing the API**: It's recommended to test the API with the "ping" route to ensure everything is running smoothly.
    ```bash
   curl http://localhost:8080/ping
   ```

## Commands

The project uses a Makefile to simplify the build process. The following commands are available:
- `make docker`: Builds the project and runs it in a detached docker container.
  - **Note** : After running this command, the logs of the API will be displayed in the terminal. To stop this behavior, press `Ctrl + C`; the container will continue to run in the background. 
- `make test`: Runs the unit tests locally.
  - **Note** : The command will make a `go mod vendor` before running the tests. This is to ensure that the tests are run with the correct dependencies.
- `make doc` : Generates, locally, the documentation for the project. 
  - **Note** : Open the 'coverage.html' file in the root project directory to view the coverage report.
- `make clean`: Removes the local generated files.

## API Endpoints

### 1. Play FizzBuzz

- **Endpoint**: `/play`
- **Method**: `POST`
- **Body**:
   ```json
   {
       "number1": <int64>,
       "number2": <int64>,
       "replace1": <string>,
       "replace2": <string>,
       "limit": <int64>
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