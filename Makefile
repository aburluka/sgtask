DOCKER_COMPOSE_ENV ?= ./docker-compose.yml

all:
	go build -o srv main.go

.PHONY: env
env: web_env_stop
	mkdir -p sg-clickhouse-data
	docker-compose -f ${DOCKER_COMPOSE_ENV} up

.PHONY: web_env_stop
env_stop:
	docker-compose -f ${DOCKER_COMPOSE_ENV} down
