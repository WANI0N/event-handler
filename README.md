# Event Handler
Simple gin app for handling events.

## Setup app
1. create `.env` file in root directory, add following variables:
```bash
export ADMIN_TOKEN=<insert_any_string>
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379
```
2. in project's root directory, run: `go get ./...`
3. install Redis (https://developer.redis.com/create/windows/)
    - run `service redis-server start` (defaults to 127.0.0.1:6379)
    - verify if redis is running by `redis-cli` command
4. source the .env file `source .env` 
5. run application using `go run main.go`
   -  API documentation should available in http://localhost:3000/docs/swagger/index.html

## Run unit tests
- tests can be run by `go test ./...` in root directory

## Smoke test via REST Client (VS Code) - optional
- install VS code extension "REST Client" https://marketplace.visualstudio.com/items?itemName=humao.rest-client
- open `devHttpClient.rest` file and declare admin_token on line:1 (same token as in .env file)
- execute calls in order: HealthCheck, CreateEvent, GetEvent, DeleteEvent by mouse click "Send Request"