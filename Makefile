# Include variables from the .envrc file
include .envrc

# ===================================================================================== #
# DEVELOPMENT
# ===================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run:
	@go run ./cmd/api -db-dsn=${TASK_LOGGER_DSN}

## db/mysql: connect to the database using
.PHONY: db/mysql
db/mysql:
	mysql ${TASK_LOGGER_DSN}
