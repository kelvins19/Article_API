# README

## Tech Stack
1. Go V1.19
2. GorillaMux Framework
3. PostgreSQL
4. golang-migrate for Database Migrations

## How to migrate
1. Install golang-migrate[https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g]
2. Run migration
> migrate -path db/migration -database "postgresql://kelvins19:123456@localhost:5432/postgres?sslmode=disable" -verbose up


## How to setup the application
1. Go to the project directory
2. Run `go build`
3. Run `go mod tidy`
4. Run `go run main.go` to start the application

## How to run the unit tests
1. Make sure the server has already been started
2. If not, run `go run main.go`
3. Go to test directory, run `go test`