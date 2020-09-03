SHELL := /bin/bash
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /config)
.DEFAULT_GOAL := all

# Make the repo
all: clean mock test

# Cleanup builds and generated code
clean:
	go clean $(shell go list ./...)
	rm -rf .gen
	rm -rf tmp

# Generate mocks for testing
mock:
	for SRC in $(shell find .$(MDIR) -not -name "*_test.go" -not -name main.go -path "./src/*.go") ; do \
		mkdir -p .gen/mock/$$SRC ; \
		rm -rf .gen/mock/$$SRC ; \
		mockgen -source=$$SRC -destination=.gen/mock/$$SRC ; \
	done

# Run tests
test:
	go test -cover $(PKGS)

# Start helix with hot-reloading
start:
	docker-compose up -d

# Stop the local instance
stop:
	docker-compose down

# Display the logs
logs:
	docker-compose logs -f

# Starts Helix and migrates the database to the latest version
bootstrap:
	make start ; \
	make db/up

# Connect to the local PostgeSQL instance
db/connect:
	docker-compose run helixdb bash

# Generates a new set of database migration files (up and down)
db/migration:
	migrate create -ext sql -dir bin/migrate/sql -seq $(name)

# Migrates the database up to the latest version
db/up:
	set -o allexport; source cfg/dev.env; set +o allexport ; \
	go run bin/migrate/main.go up

# Migrates the database down to version 0
db/down:
	set -o allexport; source cfg/dev.env; set +o allexport ; \
	go run bin/migrate/main.go down
