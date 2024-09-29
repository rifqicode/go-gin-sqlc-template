# Loyalty

## GO Version

This project currently use version 1.23.1

## Project Structure

This project structure is heavily inspired by [Golang Standards](https://github.com/golang-standards/project-layout).

## Migration
### Create New Migration
```
./scripts/migration/create-new.sh

```
### Running The Migration
```
./scripts/migration/run-migration.sh
```

### Rollback Previous Migration
```
./scripts/migration/rollback-migration.sh
```

## Generate Swagger Docs

Recompile the swagger docs
```
./scripts/swagger/compile.sh
```

Then visit http://localhost:8080/swagger/index.html to see the docs

## Run Server

```
go run main.go run-server
```
