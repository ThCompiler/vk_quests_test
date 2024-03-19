LOG_DIR=./logs
SWAG_DIRS=./internal/delivery/http/v1/,./internal/delivery/http/v1/handlers,./internal/delivery/http/v1/model/request,./internal/delivery/http/v1/model/response,./pkg/operate

.PHONY: build
build:
	go build -o server -v ./cmd

.PHONY: build-docker
build-docker:
	docker build --no-cache --network host -f ./Dockerfile . --tag main

.PHONY: run
run:
	docker compose up -d

.PHONY: run-verbose
run-verbose:
	docker compose up

.PHONY: open-last-log
open-last-log:
	cat $(LOG_DIR)/`ls -t $(LOG_DIR) | head -1 `

.PHONY: clear-logs
clear-logs:
	rm -rf $(LOG_DIR)/*.log

.PHONY: mocks
mocks:
	go generate -n $$(go list ./internal/...)

.PHONY: swag-gen
swag-gen:
	swag init --parseDependency --parseInternal --parseDepth 1 -d $(SWAG_DIRS) -g ./swag_info.go -o docs

.PHONY: swag-fmt
swag-fmt:
	swag fmt -d $(SWAG_DIRS) -g ./swag_info.go

.PHONY: run-coverage
run-coverage:
	go test -covermode=atomic -coverprofile=cover ./...
	cat cover | fgrep -v "mock" | fgrep -v "docs" | fgrep -v "config" > cover2
	go tool cover -func=cover2
