all: build-services

build-services:
	make -C query-client-service
	make -C storage-service

run-env: build-services
	make -C query-client-service
	make -C storage-service
	@docker-compose up -d

test:
	make -C query-client-service test

stop-env:
	@docker-compose stop
	@docker-compose rm

deploy-local: run-env

.PHONY: deploy-local
